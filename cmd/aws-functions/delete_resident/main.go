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

var credentialService service.CredentialService
var userService service.UserService

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK config, %v", err)
	}
	database := dynamodb.NewFromConfig(cfg)

	userRepo := repository.NewUserRepository(database, "UpKeepzTable")
	userService = *service.NewUserService(userRepo)
	credentialRepo := repository.NewCredentialRepository(database, "UpKeepzTable")
	credentialService = *service.NewCredentialService(credentialRepo)
}

func main() {
	lambda.Start(lambdacors.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, err := userService.GetUserByID(ctx)

	if err != nil {
		return response.ErrorResponse(http.StatusNotFound, "User not Found", 1004), nil
	}

	if (user.Role != model.RoleAdmin) {
		return response.ErrorResponse(http.StatusUnauthorized, "Unauthorized access", 1008), nil
	}

	id := event.QueryStringParameters["id"]

	if id == "" {
		return response.ErrorResponse(http.StatusBadRequest, "Invalid request", 1001), nil
	}

	if err := credentialService.DeleteResidentCredentials(ctx, id); err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Error deleting officer", 1010), nil
	}
	return response.SuccessResponse(nil, "Resident deleted successfully!!", http.StatusOK), nil
}