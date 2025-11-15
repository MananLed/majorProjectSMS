package main

import (
	"context"
	"log"
	"net/http"

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

var noticeService service.NoticeService

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK config, %v", err)
	}
	database := dynamodb.NewFromConfig(cfg)
	noticeRepo := repository.NewNoticeRepository(database, "UpKeepzTable")
	noticeService = *service.NewNoticeService(noticeRepo)
}

func main() {
	lambda.Start(lambdacors.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	notices, err := noticeService.GetNotices()
	if err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Server Error", 1010), nil
	}
	return response.SuccessResponse(notices, "Notices fetched successfully", http.StatusOK), nil
}
