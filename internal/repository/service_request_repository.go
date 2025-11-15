package repository

import (
	"context"
	"database/sql"
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
	GetRequestsByServiceTypeAndStatus(serviceType model.ServiceType, status model.Status) ([]model.ServiceRequest, error)
}

func NewServiceRequestRepository(ddbClient *dynamodb.Client, tableName string) *ServiceRequestRepository {
	return &ServiceRequestRepository{DynamoDbClient: ddbClient, TableName: tableName}
}

func (r *ServiceRequestRepository) CreateRequest(req *model.ServiceRequest) error {
	query := `
		INSERT INTO service_requests(request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no, date, assigned_to)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS (
		SELECT 1 FROM service_requests 
		WHERE resident_id = $1 AND service_type = $2 AND date = $3)`, req.ResidentID, req.ServiceType, req.Date).Scan(&exists)

	if err != nil {
		logger.LogToFile("error checking existing request: " + err.Error())
		return err
	}

	if exists {
		logger.LogToFile("user already has a booked request")
		return fmt.Errorf("user already has a booked request")
	}

	_, err = r.db.Exec(query, req.RequestID, req.ResidentID, req.Status, req.TimeSlot, req.StartTime, req.EndTime, req.ServiceType, req.Flat, req.Date, req.AssignedTo)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return fmt.Errorf("failed to insert request: %v", err)
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
	query := `
		UPDATE service_requests
		SET status = $1, time_slot = $2, start_time = $3, end_time = $4, service_type = $5, assigned_to = $6
		WHERE request_id = $7
	`

	res, err := r.db.Exec(query,
		req.Status, req.TimeSlot, req.StartTime, req.EndTime, req.ServiceType, req.AssignedTo, req.RequestID,
	)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return fmt.Errorf("failed to update service request: %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return fmt.Errorf("failed to check update result: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no service request found with ID %v", req.RequestID)
	}

	return nil
}

// ***********************************************************************************
func (r *ServiceRequestRepository) DeleteRequest(requestID uuid.UUID) error {
	query := `DELETE FROM service_requests WHERE request_id = $1`

	_, err := r.db.Exec(query, requestID)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error : %v", err))
		return fmt.Errorf("failed to delete service request: %v", err)
	}
	return nil
}

// **************************************************************************************
func (r *ServiceRequestRepository) DeleteRequestsByResidentID(residentID string) error {
	query := `DELETE FROM service_requests WHERE resident_id = $1`

	_, err := r.db.Exec(query, residentID)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error : %v", err))
		return fmt.Errorf("failed to delete service request: %v", err)
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
		AssignedTo    string `dyamodbav:"assigned_to"`
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

	query := `
		SELECT service_type 
		FROM service_requests
		WHERE request_id = $1
	`
	var servicetype model.ServiceType

	err := r.db.QueryRow(query, requestID).Scan(&servicetype)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error : %v", err))
		return "", fmt.Errorf("no request with such id exist: %v", err)
	}

	return servicetype, nil

}

func (r *ServiceRequestRepository) GetRequestsByServiceTypeAndStatus(serviceType model.ServiceType, status model.Status) ([]model.ServiceRequest, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pkValue AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkValue":  &types.AttributeValueMemberS{Value: "REQUESTS"},
			":skPrefix": &types.AttributeValueMemberS{Value: string(status) + "#" + string(serviceType)},
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
		AssignedTo    string `dyamodbav:"assigned_to"`
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
