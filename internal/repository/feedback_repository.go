package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"database/sql"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type FeedbackRepository struct {
	mu             sync.Mutex
	DB             *sql.DB
	DynamoDbClient *dynamodb.Client
	TableName      string
}

type FeedbackRepositoryInterface interface {
	SaveFeedback(model.Feedback) error
	GetFeedbacksByID(string) ([]model.Feedback, error)
	GetAllFeedbacks() ([]model.Feedback, error)
	IsFeedbackPresent(requestID uuid.UUID) (bool, error)
}

func NewFeedbackRepository(ddbClient *dynamodb.Client, tableName string) *FeedbackRepository {
	return &FeedbackRepository{DynamoDbClient: ddbClient, TableName: tableName}
}

func (r *FeedbackRepository) SaveFeedback(feedback model.Feedback) error {

	// var user model.User

	// queryForUserDetails := `
	// 	SELECT u.first_name, u.middle_name, u.last_name FROM
	// 	users u WHERE u.id = $1
	// `

	// row := r.DB.QueryRow(queryForUserDetails, feedback.ResidentID)

	// err := row.Scan(&user.FirstName, &user.MiddleName, &user.LastName)

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
	// }

	// var name string

	// if user.MiddleName != "" {
	// 	name = fmt.Sprintf("%s %s %s", user.FirstName, user.MiddleName, user.LastName)
	// } else {
	// 	name = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	// }

	// feedback.ResidentName = name

	// queryForUpdatingFeedbackStatus := `
	// 	UPDATE service_requests set feedback_given = $1 WHERE request_id = $2
	// `

	// _, err = r.DB.Exec(queryForUpdatingFeedbackStatus, true, feedback.RequestID)

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 	return fmt.Errorf("failed to update feedback given status: %v", err)
	// }

	// queryForServiceDetails := `
	// 	SELECT assigned_to, service_type, date, time_slot from service_requests
	// 	WHERE request_id = $1
	// `

	// row = r.DB.QueryRow(queryForServiceDetails, feedback.RequestID)

	// err = row.Scan(&feedback.AssignedTo, &feedback.ServiceType, &feedback.Date, &feedback.TimeSlot)

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
	// }

	// query := `
	// 	INSERT INTO feedbacks (id, resident_id, rating, content, flat_no, username, request_id, assigned_to, service_type, date, time_slot)
	// 	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	// `

	// _, err = r.DB.Exec(query, feedback.ID, feedback.ResidentID, feedback.Rating, feedback.Content, feedback.Flat, feedback.ResidentName, feedback.RequestID, feedback.AssignedTo, feedback.ServiceType, feedback.Date, feedback.TimeSlot)

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
	// }
	// return err

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

	type User struct {
		PK           string `dynamodbav:"PK"`
		SK           string `dynamodbav:"SK"`
		ID           string `dynamodbav:"id"`
		Email        string `dynamodbav:"email"`
		FirstName    string `dynamodbav:"first_name"`
		LastName     string `dynamodbav:"last_name"`
		MiddleName   string `dynamodbav:"middle_name"`
		MobileNumber string `dynamodbav:"mobile_number"`
		Password     string `dynamodbav:"password"`
		Role         string `dynamodbav:"role"`
		Flat         string `dynamodbav:"flat"`
	}

	var userDetails User

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

	type Request struct {
		PK            string `dynamobdav:"PK"`
		SK            string `dynamobdav:"SK"`
		AssignedTo    string `dynamodbav:"assigned_to"`
		Date          string `dynamodbav:"date"`
		FeedbackGiven bool   `dynamodbav:"feedback_given"`
		Flat          string `dynamodbav:"flat_no"`
		ID            string `dynamodbav:"id"`
		ResidentID    string `dynamodbav:"resident_id"`
		ServiceType   string `dynamodbav:"service_type"`
		Status        string `dynamodbav:"status"`
		TimeSlot      string `dynamodbav:"time_slot"`
	}

	var request Request

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

func (r *FeedbackRepository) GetFeedbacksByID(residentID string) ([]model.Feedback, error) {
	query := `
		SELECT id, resident_id, rating, content, flat_no, username, request_id, assigned_to, service_type, date, time_slot
		FROM feedbacks
		WHERE resident_id = $1
	`

	rows, err := r.DB.Query(query, residentID)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}

	defer rows.Close()

	var feedbacks []model.Feedback
	for rows.Next() {
		var f model.Feedback
		if err := rows.Scan(&f.ID, &f.ResidentID, &f.Rating, &f.Content, &f.Flat, &f.ResidentName, &f.AssignedTo, &f.ServiceType, &f.Date, &f.TimeSlot); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}
	return feedbacks, nil
}

func (r *FeedbackRepository) GetAllFeedbacks() ([]model.Feedback, error) {
	// query := `
	// 	SELECT id, resident_id, rating, content, flat_no, username, request_id, assigned_to, service_type, date, time_slot
	// 	FROM feedbacks
	// `

	// rows, err := r.DB.Query(query)

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 	return nil, err
	// }
	// defer rows.Close()

	// var feedbacks []model.Feedback
	// for rows.Next() {
	// 	var f model.Feedback
	// 	if err := rows.Scan(&f.ID, &f.ResidentID, &f.Rating, &f.Content, &f.Flat, &f.ResidentName, &f.RequestID, &f.AssignedTo, &f.ServiceType, &f.Date, &f.TimeSlot); err != nil {
	// 		logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 		return nil, err
	// 	}
	// 	feedbacks = append(feedbacks, f)
	// }
	// return feedbacks, nil

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

	type Feedback struct {
		PK          string `dynamodbav:"PK"`
		SK          string `dynamodbav:"SK"`
		AssignedTo  string `dynamodbav:"assigned_to"`
		Content     string `dynamodbav:"content"`
		Date        string `dynamodbav:"date"`
		Flat        string `dynamodbav:"flat_no"`
		ID          string `dynamodbav:"id"`
		Rating      int32  `dynamodbav:"rating"`
		RequestID   string `dynamodbav:"request_id"`
		ResidentID  string `dynamodbav:"resident_id"`
		ServiceType string `dynamodbav:"service_type"`
		TimeSlot    string `dynamodbav:"time_slot"`
		UserName    string `dynamodbav:"username"`
	}

	var feedbackDetails Feedback

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

func (r *FeedbackRepository) IsFeedbackPresent(requestID uuid.UUID) (bool, error) {

	var exists bool

	query := `SELECT EXISTS (SELECT 1 FROM feedbacks WHERE request_id = $1);`

	err := r.DB.QueryRow(query, requestID).Scan(&exists)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return false, err
	}

	return exists, nil
}
