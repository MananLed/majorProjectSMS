package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/MananLed/majorProjectSMS/internal/dto"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type ServiceRequestRepository struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

type ServiceRequestRepositoryInterface interface {
	CreateRequest(req *model.ServiceRequest) error
	UpdateRequest(req *model.ServiceRequest) error
	DeleteRequest(requestID uuid.UUID) error
	DeleteRequestsByResidentID(residentID string) error
	GetServiceRequestsByStatus(userID string, status model.Status) ([]model.ServiceRequest, error)
	GetServiceTypeByID(requestID uuid.UUID) (model.ServiceType, error)
	GetRequestsByServiceTypeAndStatus(user *model.User, serviceType model.ServiceType, status model.Status) ([]model.ServiceRequest, error)
}

func NewServiceRequestRepository(ddbClient *dynamodb.Client, tableName string) *ServiceRequestRepository {
	return &ServiceRequestRepository{DynamoDbClient: ddbClient, TableName: tableName}
}

func (r *ServiceRequestRepository) CreateRequest(req *model.ServiceRequest) error {

	fetchUserStatement := "SELECT * FROM " + r.TableName + " WHERE PK = ? AND begins_with(SK, ?)"

	result, err := r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fetchUserStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "REQUESTS"},
			&types.AttributeValueMemberS{Value: string(model.StatusPending) + "#" + string(req.ServiceType) + "#" + req.ResidentID + "#" + req.Date + "#"},
		},
	})

	if err != nil {
		return err
	}

	if len(result.Items) > 0 {
		return errors.New("user already has a booked request")
	}

	result, err = r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fetchUserStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "REQUESTS"},
			&types.AttributeValueMemberS{Value: string(model.StatusApproved) + "#" + string(req.ServiceType) + "#" + req.ResidentID + "#" + req.Date + "#"},
		},
	})

	if err != nil {
		return err
	}

	if len(result.Items) > 0 {
		return errors.New("user already has a booked request")
	}
	result, err = r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fetchUserStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "REQUESTS"},
			&types.AttributeValueMemberS{Value: string(model.StatusCompleted) + "#" + string(req.ServiceType) + "#" + req.ResidentID + "#" + req.Date + "#"},
		},
	})

	if err != nil {
		return err
	}

	if len(result.Items) > 0 {
		return errors.New("user already has a booked request")
	}

	statement := "INSERT INTO " + r.TableName + " VALUE {'PK': ?, 'SK': ?, 'assigned_to': ?, 'date': ?, 'feedback_given': ?, 'flat_no': ?, 'id': ?, 'resident_id': ?, 'service_type': ?, 'status': ?, 'time_slot': ?}"

	_, err = r.DynamoDbClient.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "REQUESTS"},
			&types.AttributeValueMemberS{Value: (string(req.Status) + "#" + string(req.ServiceType) + "#" + req.ResidentID + "#" + req.Date + "#" + req.RequestID.String())},
			&types.AttributeValueMemberS{Value: req.AssignedTo},
			&types.AttributeValueMemberS{Value: req.Date},
			&types.AttributeValueMemberBOOL{Value: req.FeedbackGiven},
			&types.AttributeValueMemberS{Value: req.Flat},
			&types.AttributeValueMemberS{Value: req.RequestID.String()},
			&types.AttributeValueMemberS{Value: req.ResidentID},
			&types.AttributeValueMemberS{Value: string(req.ServiceType)},
			&types.AttributeValueMemberS{Value: string(req.Status)},
			&types.AttributeValueMemberS{Value: req.TimeSlot},
		},
	})

	if err != nil {
		return err
	}

	_, err = r.DynamoDbClient.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: &statement,
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: req.RequestID.String()},
			&types.AttributeValueMemberS{Value: req.RequestID.String()},
			&types.AttributeValueMemberS{Value: req.AssignedTo},
			&types.AttributeValueMemberS{Value: req.Date},
			&types.AttributeValueMemberBOOL{Value: req.FeedbackGiven},
			&types.AttributeValueMemberS{Value: req.Flat},
			&types.AttributeValueMemberS{Value: req.RequestID.String()},
			&types.AttributeValueMemberS{Value: req.ResidentID},
			&types.AttributeValueMemberS{Value: string(req.ServiceType)},
			&types.AttributeValueMemberS{Value: string(req.Status)},
			&types.AttributeValueMemberS{Value: req.TimeSlot},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *ServiceRequestRepository) UpdateRequest(req *model.ServiceRequest) error {

	fetchRequestStatement := "SELECT * FROM " + r.TableName + " WHERE PK = ? AND SK = ?"

	result, err := r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fetchRequestStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: req.RequestID.String()},
			&types.AttributeValueMemberS{Value: req.RequestID.String()},
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

	if req.Status == model.StatusCompleted {
		req.AssignedTo = request.AssignedTo
	}

	updateRequestStatement := "UPDATE " + r.TableName + " SET time_slot = ?, status = ?, assigned_to = ? WHERE PK = ? AND SK = ?"

	if req.TimeSlot != request.TimeSlot && req.TimeSlot != "" {
		_, err = r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
			Statement: aws.String(updateRequestStatement),
			Parameters: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: req.TimeSlot},
				&types.AttributeValueMemberS{Value: string(req.Status)},
				&types.AttributeValueMemberS{Value: req.AssignedTo},
				&types.AttributeValueMemberS{Value: req.RequestID.String()},
				&types.AttributeValueMemberS{Value: req.RequestID.String()},
			},
		})

		if err != nil {
			return err
		}

		_, err = r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
			Statement: aws.String(updateRequestStatement),
			Parameters: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: req.TimeSlot},
				&types.AttributeValueMemberS{Value: string(req.Status)},
				&types.AttributeValueMemberS{Value: req.AssignedTo},
				&types.AttributeValueMemberS{Value: "REQUESTS"},
				&types.AttributeValueMemberS{Value: request.Status + "#" + request.ServiceType + "#" + request.ResidentID + "#" + request.Date + "#" + request.ID},
			},
		})

		if err != nil {
			return err
		}
	} else {
		deleteRequestStatement := "DELETE FROM " + r.TableName + " WHERE PK = ? AND SK = ?"

		_, err = r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
			Statement: aws.String(deleteRequestStatement),
			Parameters: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: req.RequestID.String()},
				&types.AttributeValueMemberS{Value: req.RequestID.String()},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to delete request: %v", err)
		}

		_, err = r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
			Statement: aws.String(deleteRequestStatement),
			Parameters: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "REQUESTS"},
				&types.AttributeValueMemberS{Value: request.Status + "#" + request.ServiceType + "#" + request.ResidentID + "#" + request.Date + "#" + request.ID},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to delete request: %v", err)
		}

		statement := "INSERT INTO " + r.TableName + " VALUE {'PK': ?, 'SK': ?, 'assigned_to': ?, 'date': ?, 'feedback_given': ?, 'flat_no': ?, 'id': ?, 'resident_id': ?, 'service_type': ?, 'status': ?, 'time_slot': ?}"

		_, err = r.DynamoDbClient.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
			Statement: &statement,
			Parameters: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "REQUESTS"},
				&types.AttributeValueMemberS{Value: (string(req.Status) + "#" + string(request.ServiceType) + "#" + request.ResidentID + "#" + request.Date + "#" + request.ID)},
				&types.AttributeValueMemberS{Value: req.AssignedTo},
				&types.AttributeValueMemberS{Value: request.Date},
				&types.AttributeValueMemberBOOL{Value: request.FeedbackGiven},
				&types.AttributeValueMemberS{Value: request.Flat},
				&types.AttributeValueMemberS{Value: request.ID},
				&types.AttributeValueMemberS{Value: request.ResidentID},
				&types.AttributeValueMemberS{Value: request.ServiceType},
				&types.AttributeValueMemberS{Value: string(req.Status)},
				&types.AttributeValueMemberS{Value: request.TimeSlot},
			},
		})

		if err != nil {
			return err
		}

		_, err = r.DynamoDbClient.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
			Statement: &statement,
			Parameters: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: request.ID},
				&types.AttributeValueMemberS{Value: request.ID},
				&types.AttributeValueMemberS{Value: req.AssignedTo},
				&types.AttributeValueMemberS{Value: request.Date},
				&types.AttributeValueMemberBOOL{Value: request.FeedbackGiven},
				&types.AttributeValueMemberS{Value: request.Flat},
				&types.AttributeValueMemberS{Value: request.ID},
				&types.AttributeValueMemberS{Value: request.ResidentID},
				&types.AttributeValueMemberS{Value: request.ServiceType},
				&types.AttributeValueMemberS{Value: string(req.Status)},
				&types.AttributeValueMemberS{Value: request.TimeSlot},
			},
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ServiceRequestRepository) DeleteRequest(requestID uuid.UUID) error {

	fetchRequestStatement := "SELECT * FROM " + r.TableName + " WHERE PK = ? AND SK = ?"

	result, err := r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fetchRequestStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: requestID.String()},
			&types.AttributeValueMemberS{Value: requestID.String()},
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

	if request.Status != string(model.StatusPending) {
		return errors.New("request cannot be cancelled")
	}

	deleteRequestStatement := "DELETE FROM " + r.TableName + " WHERE PK = ? AND SK = ?"

	_, err = r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(deleteRequestStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: requestID.String()},
			&types.AttributeValueMemberS{Value: requestID.String()},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete request: %v", err)
	}

	_, err = r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(deleteRequestStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "REQUESTS"},
			&types.AttributeValueMemberS{Value: request.Status + "#" + request.ServiceType + "#" + request.ResidentID + "#" + request.Date + "#" + request.ID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete request: %v", err)
	}

	return nil
}

func (r *ServiceRequestRepository) DeleteRequestsByResidentID(residentID string) error {

	fetchRequestStatement := "SELECT * FROM " + r.TableName + " WHERE PK = ? AND begins_with(SK, ?)"

	result, err := r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fetchRequestStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "REQUESTS"},
			&types.AttributeValueMemberS{Value: string(model.StatusPending) + "#" + string(model.Electrician) + "#" + residentID},
		},
	})

	if err != nil {
		return err
	}

	result1, err := r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fetchRequestStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "REQUESTS"},
			&types.AttributeValueMemberS{Value: string(model.StatusApproved) + "#" + string(model.Electrician) + "#" + residentID},
		},
	})

	if err != nil {
		return err
	}

	result2, err := r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fetchRequestStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "REQUESTS"},
			&types.AttributeValueMemberS{Value: string(model.StatusCompleted) + "#" + string(model.Electrician) + "#" + residentID},
		},
	})

	if err != nil {
		return err
	}

	result3, err := r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fetchRequestStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "REQUESTS"},
			&types.AttributeValueMemberS{Value: string(model.StatusPending) + "#" + string(model.Plumber) + "#" + residentID},
		},
	})

	if err != nil {
		return err
	}

	result4, err := r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fetchRequestStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "REQUESTS"},
			&types.AttributeValueMemberS{Value: string(model.StatusApproved) + "#" + string(model.Plumber) + "#" + residentID},
		},
	})

	if err != nil {
		return err
	}

	result5, err := r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fetchRequestStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "REQUESTS"},
			&types.AttributeValueMemberS{Value: string(model.StatusCompleted) + "#" + string(model.Plumber) + "#" + residentID},
		},
	})

	if err != nil {
		return err
	}

	result.Items = append(result.Items, result1.Items...)
	result.Items = append(result.Items, result2.Items...)
	result.Items = append(result.Items, result3.Items...)
	result.Items = append(result.Items, result4.Items...)
	result.Items = append(result.Items, result5.Items...)

	if len(result.Items) == 0 {
		return errors.New("no such request exist")
	}

	var request dto.Request

	deleteRequestStatement := "DELETE FROM " + r.TableName + " WHERE PK = ? AND SK = ?"

	for _, u := range result.Items {
		err = attributevalue.UnmarshalMap(u, &request)
		if err != nil {
			return err
		}
		_, err = r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
			Statement: aws.String(deleteRequestStatement),
			Parameters: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: request.ID},
				&types.AttributeValueMemberS{Value: request.ID},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to delete request: %v", err)
		}

		_, err = r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
			Statement: aws.String(deleteRequestStatement),
			Parameters: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "REQUESTS"},
				&types.AttributeValueMemberS{Value: request.Status + "#" + request.ServiceType + "#" + request.ResidentID + "#" + request.Date + "#" + request.ID},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to delete request: %v", err)
		}
	}
	return nil
}

func (r *ServiceRequestRepository) GetServiceRequestsByStatus(userID string, status model.Status) ([]model.ServiceRequest, error) {

	inputElectrician := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pkValue AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkValue":  &types.AttributeValueMemberS{Value: "REQUESTS"},
			":skPrefix": &types.AttributeValueMemberS{Value: (string(status) + "#" + string(model.Electrician) + "#" + userID)},
		},
	}

	ctx := context.TODO()
	response, err := r.DynamoDbClient.Query(ctx, inputElectrician)
	if err != nil {
		return nil, err
	}

	var serviceRequests []model.ServiceRequest

	for _, r := range response.Items {
		var request dto.Request
		var serviceRequest model.ServiceRequest

		err = attributevalue.UnmarshalMap(r, &request)
		if err != nil {
			return nil, err
		}

		serviceRequest.RequestID, _ = uuid.Parse(request.ID)
		serviceRequest.ResidentID = request.ResidentID
		serviceRequest.Flat = request.Flat
		serviceRequest.Status = model.Status(request.Status)
		serviceRequest.ServiceType = model.ServiceType(request.ServiceType)
		serviceRequest.TimeSlot = request.TimeSlot
		serviceRequest.Date = request.Date
		serviceRequest.AssignedTo = request.AssignedTo
		serviceRequest.FeedbackGiven = request.FeedbackGiven

		serviceRequests = append(serviceRequests, serviceRequest)
	}

	inputPlumber := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pkValue AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkValue":  &types.AttributeValueMemberS{Value: "REQUESTS"},
			":skPrefix": &types.AttributeValueMemberS{Value: (string(status) + "#" + string(model.Plumber) + "#" + userID)},
		},
	}

	response, err = r.DynamoDbClient.Query(ctx, inputPlumber)
	if err != nil {
		return nil, err
	}

	for _, r := range response.Items {
		var request dto.Request
		var serviceRequest model.ServiceRequest

		err = attributevalue.UnmarshalMap(r, &request)
		if err != nil {
			return nil, err
		}

		serviceRequest.RequestID, _ = uuid.Parse(request.ID)
		serviceRequest.ResidentID = request.ResidentID
		serviceRequest.Flat = request.Flat
		serviceRequest.Status = model.Status(request.Status)
		serviceRequest.ServiceType = model.ServiceType(request.ServiceType)
		serviceRequest.TimeSlot = request.TimeSlot
		serviceRequest.Date = request.Date
		serviceRequest.AssignedTo = request.AssignedTo
		serviceRequest.FeedbackGiven = request.FeedbackGiven

		serviceRequests = append(serviceRequests, serviceRequest)
	}

	return serviceRequests, nil
}

func (r *ServiceRequestRepository) GetServiceTypeByID(requestID uuid.UUID) (model.ServiceType, error) {

	fetchRequestStatement := "SELECT * FROM " + r.TableName + " WHERE PK = ? AND SK = ?"

	result, err := r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fetchRequestStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: requestID.String()},
			&types.AttributeValueMemberS{Value: requestID.String()},
		},
	})

	if err != nil {
		return "", err
	}

	if len(result.Items) == 0 {
		return "", errors.New("no such request exist")
	}

	var request dto.Request

	err = attributevalue.UnmarshalMap(result.Items[0], &request)
	if err != nil {
		return "", err
	}

	if request.Status != string(model.StatusPending) {
		return "", errors.New("request cannot be rescheduled")
	}

	return model.ServiceType(request.ServiceType), nil

}

func (r *ServiceRequestRepository) GetRequestsByServiceTypeAndStatus(user *model.User, serviceType model.ServiceType, status model.Status) ([]model.ServiceRequest, error) {

	var sortKey string
	if user != nil && user.Role == model.RoleResident {
		sortKey = string(status) + "#" + string(serviceType) + "#" + user.ID
	} else {
		sortKey = string(status) + "#" + string(serviceType)
	}

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pkValue AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkValue":  &types.AttributeValueMemberS{Value: "REQUESTS"},
			":skPrefix": &types.AttributeValueMemberS{Value: sortKey},
		},
	}

	ctx := context.TODO()
	response, err := r.DynamoDbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var serviceRequests []model.ServiceRequest

	for _, r := range response.Items {
		var request dto.Request
		var serviceRequest model.ServiceRequest

		err = attributevalue.UnmarshalMap(r, &request)
		if err != nil {
			return nil, err
		}

		serviceRequest.RequestID, _ = uuid.Parse(request.ID)
		serviceRequest.ResidentID = request.ResidentID
		serviceRequest.Flat = request.Flat
		serviceRequest.Status = model.Status(request.Status)
		serviceRequest.ServiceType = model.ServiceType(request.ServiceType)
		serviceRequest.TimeSlot = request.TimeSlot
		serviceRequest.Date = request.Date
		serviceRequest.AssignedTo = request.AssignedTo
		serviceRequest.FeedbackGiven = request.FeedbackGiven

		serviceRequests = append(serviceRequests, serviceRequest)
	}

	return serviceRequests, nil
}
