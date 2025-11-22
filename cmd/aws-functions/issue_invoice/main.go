package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var invoiceService service.InvoiceService
var snsClient *sns.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK config, %v", err)
	}
	database := dynamodb.NewFromConfig(cfg)
	snsClient = sns.NewFromConfig(cfg)

	invoiceRepo := repository.NewInvoiceRepository(database, "UpKeepzTable")
	invoiceService = *service.NewInvoiceService(invoiceRepo)
}

func main() {
	lambda.Start(lambdacors.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	currentUser, err := utils.GetUserFromContext(ctx)

	if err != nil || (currentUser.Role == model.RoleResident) {
		return response.ErrorResponse(http.StatusForbidden, "Unauthorized Access", 1008), nil
	}

	var req dto.Invoice

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil || req.Amount <= 0 {
		return response.ErrorResponse(http.StatusBadRequest, "Invalid Request Body", 1001), nil
	}

	now := time.Now()
	if err := invoiceService.GenerateInvoice(req.Amount, now.Month(), now.Year()); err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Failed to generate invoice", 1006), nil
	}

	_, _ = snsClient.Publish(context.TODO(), &sns.PublishInput{
		TopicArn: aws.String(os.Getenv("INVOICE_EMAIL_TOPIC_ARN")),
		Message:  aws.String(fmt.Sprintf("Monthly invoice of bill Rs.%.2f is issued, please check your dashboard", req.Amount)),
	})

	return response.SuccessResponse(nil, "Invoice issued successfully!!", http.StatusOK), nil
}
