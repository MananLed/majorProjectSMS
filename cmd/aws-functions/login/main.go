package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/MananLed/majorProjectSMS/internal/dto"
	lambdacors "github.com/MananLed/majorProjectSMS/internal/middleware/lamdba_corsmiddleware"
	"github.com/MananLed/majorProjectSMS/internal/repository"
	"github.com/MananLed/majorProjectSMS/internal/response"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/MananLed/majorProjectSMS/internal/utils"
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
	lambda.Start(lambdacors.WithCORS(handler))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req dto.LoginRequestDTO

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return response.ErrorResponse(http.StatusBadRequest, "Invalid Request Body", 1001), nil
	}

	user, err := userService.Login(strings.ToLower(strings.TrimSpace(req.Email)), req.Password)

	if err != nil {
		return response.ErrorResponse(http.StatusUnauthorized, "Invalid Credentials", 1002), nil
	}

	var jwtTokenString string
	jwtTokenString, err = utils.GenerateJWT(user.ID, string(user.Role), user.Email, user.Flat)

	if err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Server Error", 1010), nil
	}

	return response.SuccessResponse(map[string]any{"token": jwtTokenString, "email": user.Email, "role": user.Role}, "Login Successful", http.StatusCreated), nil
}
