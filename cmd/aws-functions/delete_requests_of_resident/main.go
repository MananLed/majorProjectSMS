package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/MananLed/majorProjectSMS/internal/dto"
	"github.com/MananLed/majorProjectSMS/internal/repository"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var serviceRequestService service.ServiceRequestService

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK config, %v", err)
	}
	database := dynamodb.NewFromConfig(cfg)

	serviceRepo := repository.NewServiceRequestRepository(database, "UpKeepzTable")
	serviceRequestService = *service.NewServiceRequestService(serviceRepo)
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) (error) {
	for _, msg := range sqsEvent.Records{
		var sqsMsg dto.DeleteRequestMessage

		if err := json.Unmarshal([]byte(msg.Body), &sqsMsg); err != nil {
			return err 
		}

		userId := sqsMsg.UserID

		err := serviceRequestService.DeleteServiceRequestByID(userId)

		if err != nil{
			log.Print(err)
			return err 
		}
	}
	return nil
}