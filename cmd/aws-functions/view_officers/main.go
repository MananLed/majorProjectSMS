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

var societyService service.SocietyService
var userService service.UserService

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK config, %v", err)
	}
	database := dynamodb.NewFromConfig(cfg)

	societyRepo := repository.NewSocietyRepository(database, "UpKeepzTable")
	societyService = *service.NewSocietyService(societyRepo)
	userRepo := repository.NewUserRepository(database, "UpKeepzTable")
	userService = *service.NewUserService(userRepo)
}

func main() {
	lambda.Start(lambdacors.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, err := userService.GetUserByID(ctx)

	if err != nil {
		return response.LambdaResponse(http.StatusNotFound, map[string]any{"errorCode": 1004}, "User not Found"), nil
	}

	if err != nil || (user.Role != model.RoleAdmin) {
		return response.LambdaResponse(http.StatusUnauthorized, map[string]any{"errorCode": 1008}, "Unauthorized access"), nil
	}

	officers, err := societyService.GetAllOfficers(ctx)

	if err != nil {
		return response.LambdaResponse(http.StatusInternalServerError, map[string]any{"errorCode": 1010}, "Server Error"), nil
	}

	return response.LambdaResponse(http.StatusOK, officers, "Officers retrieved successfully"), nil
}