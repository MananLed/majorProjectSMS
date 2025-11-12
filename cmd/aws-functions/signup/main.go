package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MananLed/majorProjectSMS/internal/db"
	"github.com/MananLed/majorProjectSMS/internal/dto"
	lambdacors "github.com/MananLed/majorProjectSMS/internal/middleware/lamdba_corsmiddleware"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/repository"
	"github.com/MananLed/majorProjectSMS/internal/response"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/crypto/bcrypt"
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
	lambda.Start(lambdacors.WithCORS(handler))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req dto.SignUpRequestDTO

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return response.LambdaResponse(http.StatusBadRequest, map[string]any{"errorCode": 1001}, "Invalid Request Body"), nil
	}

	if !utils.ValidateEmail(req.Email) || !utils.ValidateMobileNumber(req.Mobile) || !utils.ValidateFlatNumber(req.Flat) || !utils.ValidatePassword(req.Password) {
		return response.LambdaResponse(http.StatusBadRequest, map[string]any{"errorCode": 1001}, "Invalid Request Body"), nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.LambdaResponse(http.StatusInternalServerError, map[string]any{"errorCode": 1010}, "Server Error"), nil
	}

	user := model.User{
		ID:           utils.GenerateUUID().String(),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		MiddleName:   req.MiddleName,
		Email:        req.Email,
		MobileNumber: req.Mobile,
		Flat:         req.Flat,
		Password:     string(hashedPassword),
		Role:         model.RoleResident,
	}

	err = userService.SignUp(user)

	if err != nil {
		return response.LambdaResponse(http.StatusInternalServerError, map[string]any{"errorCode": 1010}, "Server Error"), nil
	}

	return response.LambdaResponse(http.StatusOK, map[string]any{"successCode": 1000}, "Sign Up successful"), nil
}