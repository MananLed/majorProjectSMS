package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

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
	monthStr := event.QueryStringParameters["month"]
	yearStr := event.QueryStringParameters["year"]

	if yearStr == "" && monthStr == "" || yearStr == "" && monthStr != "" {
		return response.LambdaResponse(http.StatusBadRequest, map[string]any{"errorCode": 1001}, "Invalid Request"), nil
	} else if yearStr != "" && monthStr == "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return response.LambdaResponse(http.StatusBadRequest, map[string]any{"errorCode": 1001}, "Invalid Request"), nil
		}
		notices, err := noticeService.GetNoticesByYear(year)
		if err != nil {
			return response.LambdaResponse(http.StatusInternalServerError, map[string]any{"errorCode": 1010}, "Server error"), nil
		}
		return response.LambdaResponse(http.StatusOK, notices, "Notices retrieved successfully!!"), nil
	} else {
		year, err1 := strconv.Atoi(yearStr)
		monthInt, err2 := strconv.Atoi(monthStr)
		if err1 != nil || err2 != nil || monthInt < 1 || monthInt > 12 {
			return response.LambdaResponse(http.StatusBadRequest, map[string]any{"errorCode": 1001}, "Invalid Request"), nil
		}
		month := time.Month(monthInt)
		notices, err := noticeService.GetNoticesByMonthYear(month, year)
		if err != nil {
			return response.LambdaResponse(http.StatusInternalServerError, map[string]any{"errorCode": 1010}, "Server Error"), nil
		}
		return response.LambdaResponse(http.StatusOK, notices, "Notices retrieved successfully"), nil
	}
}
