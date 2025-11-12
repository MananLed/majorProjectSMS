package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MananLed/majorProjectSMS/internal/db"
	"github.com/MananLed/majorProjectSMS/internal/dto"
	lambdacors "github.com/MananLed/majorProjectSMS/internal/middleware/lamdba_corsmiddleware"
	"github.com/MananLed/majorProjectSMS/internal/repository"
	"github.com/MananLed/majorProjectSMS/internal/response"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/MananLed/majorProjectSMS/internal/utils"
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

func main(){
	lambda.Start(lambdacors.WithCORS(handler))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req dto.LoginRequestDTO

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return response.LambdaResponse(http.StatusBadRequest, map[string]any{"errorCode": 1001}, "Invalid Request Body"), nil
	}

	user, err := userService.Login(req.Email, req.Password)

	if err != nil {
		return response.LambdaResponse(http.StatusUnauthorized, map[string]any{"errorCode": 1002}, "Invalid Credentials"), nil
	}

	var jwtTokenString string
	jwtTokenString, err = utils.GenerateJWT(user.ID, string(user.Role), user.Email, user.Flat)

	if err != nil {
		return response.LambdaResponse(http.StatusInternalServerError, map[string]any{"errorCode": 1010}, "Server Error"), nil
	}

	return response.LambdaResponse(http.StatusOK, map[string]any{"successCode": 1000, "token": jwtTokenString, "email": user.Email, "role": user.Role}, "Login Successful"), nil
}

