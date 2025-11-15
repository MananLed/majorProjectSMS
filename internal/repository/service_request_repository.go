package repository

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
	"github.com/google/uuid"
)

type ServiceRequestRepository struct {
	mu sync.Mutex
	db *sql.DB
}

type ServiceRequestRepositoryInterface interface {
	CreateRequest(req *model.ServiceRequest) error
	GetAllRequests() ([]model.ServiceRequest, error)
	GetRequestByID(requestID uuid.UUID) (*model.ServiceRequest, error)
	UpdateRequest(req *model.ServiceRequest) error
	DeleteRequest(requestID uuid.UUID) error
	DeleteRequestsByResidentID(residentID string) error
	GetServiceRequestsByStatus(userID string, status model.Status) []model.ServiceRequest
	GetServiceTypeByID(requestID uuid.UUID) (model.ServiceType, error)
	GetPendingRequestsByServiceType(serviceType model.ServiceType) []model.ServiceRequest
	GetApprovedRequestsByServiceType(serviceType model.ServiceType) []model.ServiceRequest
	GetCompletedRequestsByServiceType(serviceType model.ServiceType) []model.ServiceRequest
}

func NewServiceRequestRepository(db *sql.DB) *ServiceRequestRepository {
	return &ServiceRequestRepository{db: db}
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

	res , err := r.db.Exec(query,
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
//***********************************************************************************
func (r *ServiceRequestRepository) DeleteRequest(requestID uuid.UUID) error {
	query := `DELETE FROM service_requests WHERE request_id = $1`


	_, err := r.db.Exec(query, requestID)


	if err != nil {
		logger.LogToFile(fmt.Sprintf("error : %v", err))
		return fmt.Errorf("failed to delete service request: %v", err)
	}
	return nil
}
//**************************************************************************************
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
func (r *ServiceRequestRepository) GetServiceRequestsByStatus(userID string, status model.Status) []model.ServiceRequest {
	query := `
		SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no, date, assigned_to, feedback_given
		FROM service_requests
		WHERE status = $1 and resident_id = $2
	`

	rows, err := r.db.Query(query, status, userID)


	if err != nil {
		logger.LogToFile(fmt.Sprintf("error : %v", err))
		return nil
	}
	defer rows.Close()

	var requests []model.ServiceRequest
	for rows.Next() {
		var r model.ServiceRequest
		if err := rows.Scan(&r.RequestID, &r.ResidentID, &r.Status, &r.TimeSlot, &r.StartTime, &r.EndTime, &r.ServiceType, &r.Flat, &r.Date, &r.AssignedTo, &r.FeedbackGiven); err != nil {
			logger.LogToFile(fmt.Sprintf("error: %v", err))
			return nil
		}
		requests = append(requests, r)
	}
	return requests
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
// XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
func (r *ServiceRequestRepository) GetPendingRequestsByServiceType(serviceType model.ServiceType) []model.ServiceRequest {

	query := `
		SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no, date, assigned_to, feedback_given
		FROM service_requests
		WHERE status = $1 and service_type = $2
	`

	rows, err := r.db.Query(query, model.StatusPending, serviceType)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil
	}
	defer rows.Close()

	var requests []model.ServiceRequest
	for rows.Next() {
		var req model.ServiceRequest
		err := rows.Scan(&req.RequestID, &req.ResidentID, &req.Status, &req.TimeSlot, &req.StartTime, &req.EndTime, &req.ServiceType, &req.Flat, &req.Date, &req.AssignedTo, &req.FeedbackGiven)
		if err != nil {
			logger.LogToFile(fmt.Sprintf("error: %v", err))
			return nil
		}
		requests = append(requests, req)
	}
	return requests
}

func (r *ServiceRequestRepository) GetApprovedRequestsByServiceType(serviceType model.ServiceType) []model.ServiceRequest {

	query := `
		SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no, date, assigned_to, feedback_given
		FROM service_requests
		WHERE status = $1 and service_type = $2
	`


	rows, err := r.db.Query(query, model.StatusApproved, serviceType)


	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil
	}
	defer rows.Close()

	var requests []model.ServiceRequest
	for rows.Next() {
		var req model.ServiceRequest
		err := rows.Scan(&req.RequestID, &req.ResidentID, &req.Status, &req.TimeSlot, &req.StartTime, &req.EndTime, &req.ServiceType, &req.Flat, &req.Date, &req.AssignedTo, &req.FeedbackGiven)
		if err != nil {
			logger.LogToFile(fmt.Sprintf("error: %v", err))
			return nil
		}
		requests = append(requests, req)
	}
	return requests
}

func (r *ServiceRequestRepository) GetCompletedRequestsByServiceType(serviceType model.ServiceType) []model.ServiceRequest {

	query := `
		SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no, date, assigned_to, feedback_given
		FROM service_requests
		WHERE status = $1 and service_type = $2
	`


	rows, err := r.db.Query(query, model.StatusCompleted, serviceType)


	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil
	}
	defer rows.Close()

	var requests []model.ServiceRequest
	for rows.Next() {
		var req model.ServiceRequest
		err := rows.Scan(&req.RequestID, &req.ResidentID, &req.Status, &req.TimeSlot, &req.StartTime, &req.EndTime, &req.ServiceType, &req.Flat, &req.Date, &req.AssignedTo, &req.FeedbackGiven)
		if err != nil {
			logger.LogToFile(fmt.Sprintf("error: %v", err))
			return nil
		}
		requests = append(requests, req)
	}
	return requests
}
//XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX