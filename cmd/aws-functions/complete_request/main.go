package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	authenticationmiddleware "github.com/MananLed/majorProjectSMS/internal/middleware/lambda_authmiddleware"
	lambdacors "github.com/MananLed/majorProjectSMS/internal/middleware/lamdba_corsmiddleware"
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

	if user.Role != "admin" && user.Role != "officer" {
		return response.ErrorResponse(http.StatusForbidden, "Not authorized to complete requests", 1009), nil
	}

	path := strings.TrimPrefix(event.Path, "/service/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[0] != "complete" {
		return response.ErrorResponse(http.StatusBadRequest, "Invalid path format", 1001), nil
	}

	reqIDStr := parts[1]
	reqID, err := uuid.Parse(reqIDStr)
	if err != nil {
		return response.ErrorResponse(http.StatusBadRequest, "Invalid request ID", 1002), nil
	}

	if err := serviceRequestService.CompleteServiceRequest(reqID); err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to complete: "+err.Error(), 1008), nil
	}

	return response.SuccessResponse(nil, "Service Request completed successfully!!", http.StatusOK), nil
}