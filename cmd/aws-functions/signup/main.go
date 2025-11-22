package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/MananLed/majorProjectSMS/internal/dto"
	lambdacors "github.com/MananLed/majorProjectSMS/internal/middleware/lamdba_corsmiddleware"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/repository"
	"github.com/MananLed/majorProjectSMS/internal/response"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"golang.org/x/crypto/bcrypt"
)

var userService service.UserService
var snsClient *sns.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK config, %v", err)
	}
	database := dynamodb.NewFromConfig(cfg)
	snsClient = sns.NewFromConfig(cfg)

	userRepo := repository.NewUserRepository(database, "UpKeepzTable")
	userService = *service.NewUserService(userRepo)
}

func main() {
	lambda.Start(lambdacors.WithCORS(handler))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req dto.SignUpRequestDTO

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return response.ErrorResponse(http.StatusBadRequest, "Invalid Request Body", 1001), nil
	}

	mockContext := context.Background()
	mockContext = context.WithValue(mockContext, utils.UserIDKey, "mock")
	mockContext = context.WithValue(mockContext, utils.UserEmailKey, strings.ToLower(strings.TrimSpace(req.Email)))
	mockContext = context.WithValue(mockContext, utils.UserFlatKey, "mock")
	mockContext = context.WithValue(mockContext, utils.UserRoleKey, "resident")

	existingUser, err := userService.GetUserByID(mockContext)
	if err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Server Error", 1010), nil
	}
	if existingUser != nil{
		return response.ErrorResponse(http.StatusBadRequest, "User with same email already exists", 1001), nil
	}

	if !utils.ValidateEmail(req.Email) || !utils.ValidateMobileNumber(req.Mobile) || !utils.ValidateFlatNumber(req.Flat) || !utils.ValidatePassword(req.Password) {
		return response.ErrorResponse(http.StatusBadRequest, "Invalid Request Body", 1001), nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Server Error", 1010), nil
	}

	user := model.User{
		ID:           utils.GenerateUUID().String(),
		FirstName:    strings.TrimSpace(req.FirstName),
		LastName:     strings.TrimSpace(req.LastName),
		MiddleName:   strings.TrimSpace(req.MiddleName),
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		MobileNumber: strings.TrimSpace(req.Mobile),
		Flat:         strings.TrimSpace(req.Flat),
		Password:     string(hashedPassword),
		Role:         model.RoleResident,
	}

	err = userService.SignUp(user)

	if err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Server Error", 1010), nil
	}

	_, _ = snsClient.Subscribe(context.TODO(), &sns.SubscribeInput{
		TopicArn: aws.String(os.Getenv("INVOICE_EMAIL_TOPIC_ARN")),
		Protocol: aws.String("email"),
		Endpoint: aws.String(user.Email),
	})

	return response.SuccessResponse(nil, "Sign Up successful", http.StatusCreated), nil
}
