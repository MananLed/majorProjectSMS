package main

import (
	"context"
	"log"
	"net/http"

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

	if user.Role != "admin" && user.Role != "officer" {
		return response.ErrorResponse(http.StatusForbidden, "Not authorized", 1008), nil
	}

	var pendingRequests []model.ServiceRequest
	var approvedRequests []model.ServiceRequest
	var completedRequests []model.ServiceRequest

	pendingRequests, err = serviceRequestService.GetPendingRequestsByServiceType(model.Plumber)
	if err != nil{
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to fetch requests", 1010), nil
	}
	pendingRequestElectrician, err := serviceRequestService.GetPendingRequestsByServiceType(model.Electrician)
	if err != nil{
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to fetch requests", 1010), nil
	}
	pendingRequests = append(pendingRequests, pendingRequestElectrician...)

	approvedRequests, err = serviceRequestService.GetApprovedRequestsByServiceType(model.Plumber)
	if err != nil{
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to fetch requests", 1010), nil
	}
	approvedRequestsElectrician, err := serviceRequestService.GetApprovedRequestsByServiceType(model.Electrician)
	if err != nil{
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to fetch requests", 1010), nil
	}
	approvedRequests = append(approvedRequests, approvedRequestsElectrician...)

	completedRequests, err = serviceRequestService.GetCompletedRequestsByServiceType(model.Plumber)
	if err != nil{
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to fetch requests", 1010), nil
	}
	completedRequestsElectrician, err := serviceRequestService.GetCompletedRequestsByServiceType(model.Electrician)
	if err != nil{
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to fetch requests", 1010), nil
	}
	completedRequests = append(completedRequests, completedRequestsElectrician...)

	allRequests := struct {
		Pending   []model.ServiceRequest
		Approved  []model.ServiceRequest
		Completed []model.ServiceRequest
	}{
		Pending:   pendingRequests,
		Approved:  approvedRequests,
		Completed: completedRequests,
	}

	return response.SuccessResponse(allRequests, "Requests fetched successfully", http.StatusOK), nil
}
