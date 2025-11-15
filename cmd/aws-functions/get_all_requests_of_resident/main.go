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

	var pendingRequests []model.ServiceRequest
	var approvedRequests []model.ServiceRequest
	var completedRequests []model.ServiceRequest

	pendingRequests, err = serviceRequestService.GetServiceRequestsByStatus(user.ID, model.StatusPending)
	if err != nil{
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to fetch requests", 1010), nil
	}

	approvedRequests, err = serviceRequestService.GetServiceRequestsByStatus(user.ID, model.StatusApproved)
	if err != nil{
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to fetch requests", 1010), nil
	}

	completedRequests, err = serviceRequestService.GetServiceRequestsByStatus(user.ID, model.StatusCompleted)
	if err != nil{
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to fetch requests", 1010), nil
	}

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
