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
	"golang.org/x/crypto/bcrypt"
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

	user, err := utils.GetUserFromContext(ctx)

	if err != nil {
		return response.ErrorResponse(http.StatusNotFound, "User not Found", 1004), nil
	}

	if user.Role != model.RoleAdmin {
		return response.ErrorResponse(http.StatusUnauthorized, "Unauthorized access", 1008), nil
	}

	var req dto.OfficerDetails

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return response.ErrorResponse(http.StatusBadRequest, "Invalid Request Body", 1001), nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to hash password", 1010), nil
	}

	newOfficer := model.User{
		Email:        req.Email,
		ID:           utils.GenerateUUID().String(),
		Password:     string(hashedPassword),
		Role:         model.RoleOfficer,
		FirstName:    "********",
		LastName:     "*******",
		MobileNumber: "**********",
		Flat:         "xxx",
	}

	if err := userService.SignUp(newOfficer); err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to create officer", 1010), nil
	}

	return response.SuccessResponse(newOfficer.ID, "Officer created successfully", http.StatusOK), nil
}
