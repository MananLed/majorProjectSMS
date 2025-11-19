package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type ServiceRequestRepository struct {
	mu             sync.Mutex
	db             *sql.DB
	DynamoDbClient *dynamodb.Client
	TableName      string
}

type ServiceRequestRepositoryInterface interface {
	CreateRequest(req *model.ServiceRequest) error
	GetAllRequests() ([]model.ServiceRequest, error)
	GetRequestByID(requestID uuid.UUID) (*model.ServiceRequest, error)
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
	// query := `
	// 	INSERT INTO service_requests(request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no, date, assigned_to)
	// 	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	// `

	// var exists bool
	// err := r.db.QueryRow(
	// 	`SELECT EXISTS (
	// 	SELECT 1 FROM service_requests
	// 	WHERE resident_id = $1 AND service_type = $2 AND date = $3)`, req.ResidentID, req.ServiceType, req.Date).Scan(&exists)

	// if err != nil {
	// 	logger.LogToFile("error checking existing request: " + err.Error())
	// 	return err
	// }

	// if exists {
	// 	logger.LogToFile("user already has a booked request")
	// 	return fmt.Errorf("user already has a booked request")
	// }

	// _, err = r.db.Exec(query, req.RequestID, req.ResidentID, req.Status, req.TimeSlot, req.StartTime, req.EndTime, req.ServiceType, req.Flat, req.Date, req.AssignedTo)

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 	return fmt.Errorf("failed to insert request: %v", err)
	// }
	// return nil

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

func (r *ServiceRequestRepository) GetAllRequests() ([]model.ServiceRequest, error) {
	query := `
		SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no, date, assigned_to, feedback_given
		FROM service_requests
	`

	rows, err := r.db.Query(query)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, fmt.Errorf("failed to load service requests: %v", err)
	}
	defer rows.Close()

	var requests []model.ServiceRequest
	for rows.Next() {
		var req model.ServiceRequest
		err := rows.Scan(&req.RequestID, &req.ResidentID, &req.Status, &req.TimeSlot, &req.StartTime, &req.EndTime, &req.ServiceType, &req.Flat, &req.Date, &req.AssignedTo, &req.FeedbackGiven)
		if err != nil {
			logger.LogToFile(fmt.Sprintf("error: %v", err))
			return nil, fmt.Errorf("failed to retrieve service request: %v", err)
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (r *ServiceRequestRepository) GetRequestByID(requestID uuid.UUID) (*model.ServiceRequest, error) {
	query := `
		SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no, date, assigned_to, feedback_given
		FROM service_requests
		WHERE request_id = $1
	`

	row := r.db.QueryRow(query, requestID)

	var req model.ServiceRequest
	err := row.Scan(&req.RequestID, &req.ResidentID, &req.Status, &req.TimeSlot, &req.StartTime, &req.EndTime, &req.ServiceType, &req.Flat, &req.Date, &req.AssignedTo, &req.FeedbackGiven)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, fmt.Errorf("failed to fetch service request: %v", err)
	}
	return &req, nil
}

func (r *ServiceRequestRepository) UpdateRequest(req *model.ServiceRequest) error {
	// query := `
	// 	UPDATE service_requests
	// 	SET status = $1, time_slot = $2, start_time = $3, end_time = $4, service_type = $5, assigned_to = $6
	// 	WHERE request_id = $7
	// `

	// res, err := r.db.Exec(query,
	// 	req.Status, req.TimeSlot, req.StartTime, req.EndTime, req.ServiceType, req.AssignedTo, req.RequestID,
	// )

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 	return fmt.Errorf("failed to update service request: %v", err)
	// }

	// rowsAffected, err := res.RowsAffected()

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 	return fmt.Errorf("failed to check update result: %v", err)
	// }

	// if rowsAffected == 0 {
	// 	return fmt.Errorf("no service request found with ID %v", req.RequestID)
	// }

	// return nil

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

// ***********************************************************************************
func (r *ServiceRequestRepository) DeleteRequest(requestID uuid.UUID) error {
	// query := `DELETE FROM service_requests WHERE request_id = $1`

	// _, err := r.db.Exec(query, requestID)

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error : %v", err))
	// 	return fmt.Errorf("failed to delete service request: %v", err)
	// }
	// return nil

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

// **************************************************************************************
func (r *ServiceRequestRepository) DeleteRequestsByResidentID(residentID string) error {
	// query := `DELETE FROM service_requests WHERE resident_id = $1`

	// _, err := r.db.Exec(query, residentID)

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error : %v", err))
	// 	return fmt.Errorf("failed to delete service request: %v", err)
	// }
	// return nil

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

// ***************************************************************************************
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

	var serviceRequests []model.ServiceRequest

	for _, r := range response.Items {
		var request Request
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
		var request Request
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

	// query := `
	// 	SELECT service_type
	// 	FROM service_requests
	// 	WHERE request_id = $1
	// `
	// var servicetype model.ServiceType

	// err := r.db.QueryRow(query, requestID).Scan(&servicetype)

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error : %v", err))
	// 	return "", fmt.Errorf("no request with such id exist: %v", err)
	// }

	// return servicetype, nil

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

	var serviceRequests []model.ServiceRequest

	for _, r := range response.Items {
		var request Request
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
