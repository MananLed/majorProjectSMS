package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/MananLed/majorProjectSMS/internal/dto"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type FeedbackRepository struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

type FeedbackRepositoryInterface interface {
	SaveFeedback(model.Feedback) error
	GetAllFeedbacks() ([]model.Feedback, error)
}

func NewFeedbackRepository(ddbClient *dynamodb.Client, tableName string) *FeedbackRepository {
	return &FeedbackRepository{DynamoDbClient: ddbClient, TableName: tableName}
}

func (r *FeedbackRepository) SaveFeedback(feedback model.Feedback) error {

	if feedback.ID == uuid.Nil {
		feedback.ID = utils.GenerateUUID()
	}

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pkValue AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkValue":  &types.AttributeValueMemberS{Value: "ROLE#" + string(model.RoleResident)},
			":skPrefix": &types.AttributeValueMemberS{Value: feedback.ResidentID},
		},
	}

	ctx := context.TODO()
	response, err := r.DynamoDbClient.Query(ctx, input)
	if err != nil {
		return err
	}

	var userDetails dto.User

	err = attributevalue.UnmarshalMap(response.Items[0], &userDetails)
	if err != nil {
		return err
	}

	var name string

	if userDetails.MiddleName != "" {
		name = fmt.Sprintf("%s %s %s", userDetails.FirstName, userDetails.MiddleName, userDetails.LastName)
	} else {
		name = fmt.Sprintf("%s %s", userDetails.FirstName, userDetails.LastName)
	}

	feedback.ResidentName = name

	fetchRequestStatement := "SELECT * FROM " + r.TableName + " WHERE PK = ? AND SK = ?"

	result, err := r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fetchRequestStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: feedback.RequestID.String()},
			&types.AttributeValueMemberS{Value: feedback.RequestID.String()},
		},
	})

	if err != nil {
		return err
	}

	if len(result.Items) == 0 {
		return errors.New("no such request exist")
	}

	var request dto.Request

	err = attributevalue.UnmarshalMap(result.Items[0], &request)
	if err != nil {
		return err
	}

	updateRequestStatement := "UPDATE " + r.TableName + " SET feedback_given = ? WHERE PK = ? AND SK = ?"

	_, err = r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(updateRequestStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberBOOL{Value: true},
			&types.AttributeValueMemberS{Value: feedback.RequestID.String()},
			&types.AttributeValueMemberS{Value: feedback.RequestID.String()},
		},
	})

	if err != nil {
		return err
	}

	_, err = r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(updateRequestStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberBOOL{Value: true},
			&types.AttributeValueMemberS{Value: "REQUESTS"},
			&types.AttributeValueMemberS{Value: string(request.Status) + "#" + string(request.ServiceType) + "#" + string(request.ResidentID) + "#" + request.Date + "#" + request.ID},
		},
	})

	if err != nil {
		return err
	}

	statement := "INSERT INTO " + r.TableName + " VALUE {'PK': ?, 'SK': ?, 'assigned_to': ?, 'content': ?, 'date': ?, 'flat_no': ?, 'id': ?, 'rating': ?, 'request_id': ?, 'resident_id': ?, 'service_type': ?, 'time_slot': ?, 'username': ?}"

	_, err = r.DynamoDbClient.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "FEEDBACKS"},
			&types.AttributeValueMemberS{Value: request.ID + "#" + feedback.ID.String()},
			&types.AttributeValueMemberS{Value: request.AssignedTo},
			&types.AttributeValueMemberS{Value: feedback.Content},
			&types.AttributeValueMemberS{Value: request.Date},
			&types.AttributeValueMemberS{Value: request.Flat},
			&types.AttributeValueMemberS{Value: feedback.ID.String()},
			&types.AttributeValueMemberN{Value: fmt.Sprintf("%d", feedback.Rating)},
			&types.AttributeValueMemberS{Value: feedback.RequestID.String()},
			&types.AttributeValueMemberS{Value: feedback.ResidentID},
			&types.AttributeValueMemberS{Value: request.ServiceType},
			&types.AttributeValueMemberS{Value: request.TimeSlot},
			&types.AttributeValueMemberS{Value: name},
		},
	})

	if err != nil {
		return err
	}

	return nil
}


func (r *FeedbackRepository) GetAllFeedbacks() ([]model.Feedback, error) {

	var feedback model.Feedback
	var feedbacks []model.Feedback

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pkValue"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkValue": &types.AttributeValueMemberS{Value: "FEEDBACKS"},
		},
	}

	ctx := context.TODO()
	response, err := r.DynamoDbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}
	
	var feedbackDetails dto.FeedbackDDB

	if err != nil {
		return nil, err
	} else {
		for _, n := range response.Items {
			err = attributevalue.UnmarshalMap(n, &feedbackDetails)
			if err != nil {
				return nil, err
			}

			feedback.ID, _ = uuid.Parse(feedbackDetails.ID)
			feedback.ResidentID = feedbackDetails.ResidentID
			feedback.Flat = feedbackDetails.Flat
			feedback.Rating = feedbackDetails.Rating
			feedback.Content = feedbackDetails.Content
			feedback.ResidentName = feedbackDetails.UserName
			feedback.RequestID, _ = uuid.Parse(feedbackDetails.RequestID)
			feedback.AssignedTo = feedbackDetails.AssignedTo
			feedback.ServiceType = feedbackDetails.ServiceType
			feedback.Date = feedbackDetails.Date
			feedback.TimeSlot = feedbackDetails.TimeSlot

			feedbacks = append(feedbacks, feedback)
		}
	}

	return feedbacks, nil
}
