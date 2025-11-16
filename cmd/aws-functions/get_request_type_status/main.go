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

	status := event.QueryStringParameters["status"]
	serviceType := event.QueryStringParameters["serviceType"]
	if status == "" || serviceType == "" {
		return response.ErrorResponse(http.StatusBadRequest, "Missing status or serviceType parameter", 1001), nil
	}

	var requests []model.ServiceRequest

	switch {
	case serviceType == "plumber" && status == "pending":
		requests, err = serviceRequestService.GetPendingRequestsByServiceType(user, model.Plumber)
	case serviceType == "plumber" && status == "approved":
		requests, err = serviceRequestService.GetApprovedRequestsByServiceType(user, model.Plumber)
	case serviceType == "electrician" && status == "pending":
		requests, err = serviceRequestService.GetPendingRequestsByServiceType(user, model.Electrician)
	case serviceType == "electrician" && status == "approved":
		requests, err = serviceRequestService.GetApprovedRequestsByServiceType(user, model.Electrician)
	}

	if err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to fetch requests", 1010), nil
	}

	return response.SuccessResponse(requests, "Requests fetched successfully!!", http.StatusOK), nil
}
