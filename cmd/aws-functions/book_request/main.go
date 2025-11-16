package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/dto"
	authenticationmiddleware "github.com/MananLed/majorProjectSMS/internal/middleware/lambda_authmiddleware"
	lambdacors "github.com/MananLed/majorProjectSMS/internal/middleware/lamdba_corsmiddleware"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/repository"
	"github.com/MananLed/majorProjectSMS/internal/response"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

var serviceRequestService service.ServiceRequestService

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK config, %v", err)
	}
	database := dynamodb.NewFromConfig(cfg)

	serviceRepo := repository.NewServiceRequestRepository(database, "UpKeepzTable")
	serviceRequestService = *service.NewServiceRequestService(serviceRepo)
}

func main() {
	lambda.Start(lambdacors.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, err := utils.GetUserFromContext(ctx)

	if err != nil {
		return response.ErrorResponse(http.StatusUnauthorized, "User not found", 1007), nil
	}

	if user.Role != model.RoleResident {
		return response.ErrorResponse(http.StatusForbidden, "Unauthorized Access", 1008), nil
	}

	var req dto.ServiceRequest

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return response.ErrorResponse(http.StatusBadRequest, "Invalid Request Body", 1001), nil
	}

	availableSlots, err := serviceRequestService.GetAvailableTimeSlots(model.ServiceType(req.ServiceType))

	if err != nil{
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to book request", 1010), nil
	}

	if req.SlotID < 1 || req.SlotID > len(availableSlots) {
		return response.ErrorResponse(http.StatusBadRequest, "Invalid request", 1001), nil
	}
	chosenSlot := availableSlots[req.SlotID-1]

	now := time.Now()
	formattedDate := now.Format("02-01-2006")

	request := model.ServiceRequest{
		RequestID:   uuid.New(),
		ResidentID:  user.ID,
		Flat:        user.Flat,
		Status:      model.StatusPending,
		TimeSlot:    chosenSlot.StartTime.Format("3:04 PM") + " - " + chosenSlot.EndTime.Format("3:04 PM"),
		StartTime:   chosenSlot.StartTime,
		EndTime:     chosenSlot.EndTime,
		ServiceType: model.ServiceType(req.ServiceType),
		Date:        formattedDate,
	}

	err = serviceRequestService.BookServiceRequest(request)

	if err != nil {
		log.Print(err)
		log.Print(err)
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to book service request", 1010), nil
	}

	return response.SuccessResponse(request.RequestID, "Service Request created successfully!!", http.StatusOK), nil
}