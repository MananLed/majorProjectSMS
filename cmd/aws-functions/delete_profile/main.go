package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/MananLed/majorProjectSMS/internal/dto"
	authenticationmiddleware "github.com/MananLed/majorProjectSMS/internal/middleware/lambda_authmiddleware"
	lambdacors "github.com/MananLed/majorProjectSMS/internal/middleware/lamdba_corsmiddleware"
	"github.com/MananLed/majorProjectSMS/internal/repository"
	"github.com/MananLed/majorProjectSMS/internal/response"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var userService service.UserService
var sqsClient *sqs.Client
var queueURL string

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK config, %v", err)
	}
	database := dynamodb.NewFromConfig(cfg)

	sqsClient = sqs.NewFromConfig(cfg)

	queueURL = os.Getenv("QUEUE_URL")

	if queueURL == ""{
		panic("QUEUE_URL is not set")
	}

	userRepo := repository.NewUserRepository(database, "UpKeepzTable")
	userService = *service.NewUserService(userRepo)
}

func main() {
	lambda.Start(lambdacors.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		return response.ErrorResponse(http.StatusUnauthorized, "User not found", 1007), nil
	}

	msgBody := dto.DeleteRequestMessage{
		UserID: user.ID,
	}

	msgBodyBytes, err := json.Marshal(msgBody)
	if err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to marshal body", 1010), nil
	}

	_, err = sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl: aws.String(queueURL),
		MessageBody: aws.String(string(msgBodyBytes)),
	})
	if err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to send message to queue", 1010), nil
	}

	err = userService.DeleteProfile(ctx)
	if err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Error deleting user", 1010), nil
	}
	return response.SuccessResponse(nil, "Profile deleted successfully", http.StatusOK), nil
}
