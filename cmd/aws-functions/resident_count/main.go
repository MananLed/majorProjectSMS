package main

import (
	"context"
	"log"
	"net/http"

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
)

var societyService service.SocietyService

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK config, %v", err)
	}
	database := dynamodb.NewFromConfig(cfg)

	societyRepo := repository.NewSocietyRepository(database, "UpKeepzTable")
	societyService = *service.NewSocietyService(societyRepo)
}

func main() {
	lambda.Start(lambdacors.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, err := utils.GetUserFromContext(ctx)

	if err != nil {
		return response.ErrorResponse(http.StatusNotFound, "User not Found", 1004), nil
	}

	if (user.Role == model.RoleResident) {
		return response.ErrorResponse(http.StatusUnauthorized, "Unauthorized access", 1008), nil
	}

	residents, err := societyService.GetAllResidents(ctx)

	if err != nil {
		return response.ErrorResponse(http.StatusInternalServerError, "Server Error", 1010), nil
	}

	var countOfResidents int = len(residents)

	return response.SuccessResponse(countOfResidents, "Residents count retrieved successfully", http.StatusOK), nil
}