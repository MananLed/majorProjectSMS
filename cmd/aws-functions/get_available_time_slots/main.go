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

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
	serviceType := event.QueryStringParameters["serviceType"]

	if serviceType != string(model.Plumber) && serviceType != string(model.Electrician) {
		return response.ErrorResponse(http.StatusBadRequest, "Invalid request", 1001), nil
	}

	slots, err := serviceRequestService.GetAvailableTimeSlots(model.ServiceType(serviceType))
	if err != nil{
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to fetch available time slots", 1010), nil
	}

	return response.SuccessResponse(slots, "Time slots fetched successfully", http.StatusOK), nil
}