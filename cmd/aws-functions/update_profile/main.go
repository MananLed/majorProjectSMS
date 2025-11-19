package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/MananLed/majorProjectSMS/internal/dto"
	authenticationmiddleware "github.com/MananLed/majorProjectSMS/internal/middleware/lambda_authmiddleware"
	lambdacors "github.com/MananLed/majorProjectSMS/internal/middleware/lamdba_corsmiddleware"
	"github.com/MananLed/majorProjectSMS/internal/repository"
	"github.com/MananLed/majorProjectSMS/internal/response"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var userService service.UserService

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK config, %v", err)
	}
	database := dynamodb.NewFromConfig(cfg)

	userRepo := repository.NewUserRepository(database, "UpKeepzTable")
	userService = *service.NewUserService(userRepo)
}

func main() {
	lambda.Start(lambdacors.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, err := userService.GetUserByID(ctx)
	if err != nil {
		return response.ErrorResponse(http.StatusUnauthorized, "User not found", 1007), nil
	}

	var req dto.UpdateProfile
	var previousEmail string

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return response.ErrorResponse(http.StatusBadRequest, "Invalid Request Body", 1001), nil
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.MiddleName != "" {
		user.MiddleName = req.MiddleName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		previousEmail = user.Email
		user.Email = req.Email
	}
	if req.MobileNumber != "" {
		user.MobileNumber = req.MobileNumber
	}

	err = userService.UpdateProfile(*user, previousEmail)

	if err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to update user", 1011), nil
	}

	return response.SuccessResponse(nil, "User updated successfully", http.StatusOK), nil
}