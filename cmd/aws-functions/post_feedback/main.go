package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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

var feedbackService service.FeedbackService

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK config, %v", err)
	}
	database := dynamodb.NewFromConfig(cfg)
	feedbackRepo := repository.NewFeedbackRepository(database, "UpKeepzTable")
	feedbackService = *service.NewFeedbackService(feedbackRepo)
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
		return response.ErrorResponse(http.StatusUnauthorized, "unauthorized", 1007), nil
	}

	var req dto.Feedback

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return response.ErrorResponse(http.StatusBadRequest, "Invalid Request Body", 1001), nil
	}

	if req.Rating < 1 || req.Rating > 5 {
		return response.ErrorResponse(http.StatusBadRequest, "Invalid rating", 1001), nil
	}

	if len(req.Content) > 500 {
		return response.ErrorResponse(http.StatusBadRequest, "Content length exceeded the permitted lenght", 1001), nil
	}

	reqID, err := uuid.Parse(req.RequestID)
	if err != nil {
		return response.ErrorResponse(http.StatusBadRequest, "Invalid request ID", 1002), nil
	}

	if err := feedbackService.IssueFeedbackOnRequest(req.Content, user.ID, user.Flat, req.Rating, reqID); err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to issue feedback", 1010), nil
	}

	return response.SuccessResponse(nil, "Feedback issued successfully!!", http.StatusOK), nil
}
