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

var invoiceService service.InvoiceService

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK config, %v", err)
	}
	database := dynamodb.NewFromConfig(cfg)
	invoiceRepo := repository.NewInvoiceRepository(database, "UpKeepzTable")
	invoiceService = *service.NewInvoiceService(invoiceRepo)
}

func main() {
	lambda.Start(lambdacors.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	monthStr := event.QueryStringParameters["month"]
	yearStr := event.QueryStringParameters["year"]

	if yearStr == "" && monthStr == "" || yearStr == "" && monthStr != "" {
		return response.ErrorResponse(http.StatusBadRequest, "Invalid Request", 1001), nil
	} else if yearStr != "" && monthStr == "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return response.ErrorResponse(http.StatusBadRequest, "Invalid Request", 1001), nil
		}
		invoices, err := invoiceService.GetInvoicesByYear(year)
		if err != nil {
			return response.ErrorResponse(http.StatusInternalServerError, "Server error", 1010), nil
		}
		return response.SuccessResponse(invoices, "Invoices retrieved successfully!!", http.StatusOK), nil
	} else {
		year, err1 := strconv.Atoi(yearStr)
		monthInt, err2 := strconv.Atoi(monthStr)
		if err1 != nil || err2 != nil || monthInt < 1 || monthInt > 12 {
			return response.ErrorResponse(http.StatusBadRequest, "Invalid Request", 1001), nil
		}
		month := time.Month(monthInt)
		invoice, err := invoiceService.GetInvoiceByMonthAndYear(month, year)
		if err != nil {
			return response.ErrorResponse(http.StatusInternalServerError, "Server Error", 1010), nil
		}
		return response.SuccessResponse(invoice, "Invoice retrieved successfully", http.StatusOK), nil
	}
}
