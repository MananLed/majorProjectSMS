package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MananLed/majorProjectSMS/internal/db"
	authenticationmiddleware "github.com/MananLed/majorProjectSMS/internal/middleware/lambda_authmiddleware"
	lambdacors "github.com/MananLed/majorProjectSMS/internal/middleware/lamdba_corsmiddleware"
	"github.com/MananLed/majorProjectSMS/internal/repository"
	"github.com/MananLed/majorProjectSMS/internal/response"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var userService service.UserService

func init() {
	database, err := db.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize DB: %v", err))
	}

	userRepo := repository.NewUserRepository(database)
	userService = *service.NewUserService(userRepo)
}

func main() {
	lambda.Start(lambdacors.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
	user, err := userService.GetUserByID(ctx)

	if err != nil{
		return response.LambdaResponse(http.StatusNotFound, map[string]any{"errorCode": 1004}, "User not Found"), nil
	}

	return response.LambdaResponse(http.StatusOK, map[string]any{"successCode": 1000, "profile": user}, "User retrieved successfully"), nil
}
