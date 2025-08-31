package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/google/uuid"
)

func newServiceRequestRepo(t *testing.T) (*ServiceRequestRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	repo := NewServiceRequestRepository(db)
	cleanup := func() { db.Close() }
	return repo, mock, cleanup
}

func TestServiceRequestRepository_CreateRequest(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	req := &model.ServiceRequest{
		RequestID:   uuid.New(),
		ResidentID:  "resident-123",
		Status:      model.StatusPending,
		TimeSlot:    "10:00-10:45",
		StartTime:   time.Date(2025, time.January, 1, 10, 0, 0, 0, time.UTC),
		EndTime:     time.Date(2025, time.January, 1, 10, 45, 0, 0, time.UTC),
		ServiceType: "Plumber",
		Flat:        "A-101",
	}

	t.Run("error in checking existing request", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS (`)).
			WithArgs(req.ResidentID, req.ServiceType).
			WillReturnError(sqlmock.ErrCancelled)

		err := repo.CreateRequest(req)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("already has request", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS (`)).
			WithArgs(req.ResidentID, req.ServiceType).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		err := repo.CreateRequest(req)
		if err == nil || err.Error() != "user already has a booked request" {
			t.Errorf("expected user already has a booked request error, got %v", err)
		}
	})

	t.Run("insert fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS (`)).
			WithArgs(req.ResidentID, req.ServiceType).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO service_requests`)).
			WithArgs(req.RequestID, req.ResidentID, req.Status, req.TimeSlot,
				req.StartTime, req.EndTime, req.ServiceType, req.Flat).
			WillReturnError(sqlmock.ErrCancelled)

		err := repo.CreateRequest(req)
		if err == nil {
			t.Errorf("expected error from insert, got nil")
		}
	})

	t.Run("insert success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS (`)).
			WithArgs(req.ResidentID, req.ServiceType).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO service_requests`)).
			WithArgs(req.RequestID, req.ResidentID, req.Status, req.TimeSlot,
				req.StartTime, req.EndTime, req.ServiceType, req.Flat).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateRequest(req)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}

///

func TestServiceRequestRepository_GetAllRequests_Success(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"request_id", "resident_id", "status", "time_slot", "start_time", "end_time", "service_type", "flat_no",
	}).
		AddRow(uuid.New(), "res-123", model.StatusPending, "09:00-09:45", time.Now(), time.Now().Add(45*time.Minute), model.Electrician, "A-101").
		AddRow(uuid.New(), "res-456", model.StatusApproved, "10:00-10:45", time.Now(), time.Now().Add(45*time.Minute), model.Plumber, "B-202")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		 FROM service_requests`,
	)).WillReturnRows(rows)

	requests, err := repo.GetAllRequests()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(requests))
	}
}

func TestServiceRequestRepository_GetAllRequests_Empty(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"request_id", "resident_id", "status", "time_slot", "start_time", "end_time", "service_type", "flat_no",
	})

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		 FROM service_requests`,
	)).WillReturnRows(rows)

	requests, err := repo.GetAllRequests()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(requests) != 0 {
		t.Fatalf("expected 0 requests, got %d", len(requests))
	}
}

func TestServiceRequestRepository_GetAllRequests_QueryError(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		 FROM service_requests`,
	)).WillReturnError(errors.New("query failed"))

	_, err := repo.GetAllRequests()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestServiceRequestRepository_GetAllRequests_ScanError(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"request_id", "resident_id", "status", "time_slot", "start_time", "end_time", "service_type", "flat_no",
	}).AddRow("not-a-uuid", "res-123", model.StatusPending, "09:00-09:45", time.Now(), time.Now().Add(45*time.Minute), model.Electrician, "A-101")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		 FROM service_requests`,
	)).WillReturnRows(rows)

	_, err := repo.GetAllRequests()
	if err == nil {
		t.Fatalf("expected scan error, got nil")
	}
}



func TestServiceRequestRepository_GetRequestByID_Success(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	id := uuid.New()
	start := time.Now()
	end := start.Add(45 * time.Minute)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		FROM service_requests
		WHERE request_id = $1
	`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{
			"request_id", "resident_id", "status", "time_slot", "start_time", "end_time", "service_type", "flat_no",
		}).AddRow(id, "resident-123", model.StatusPending, "10:00-10:45", start, end, model.ServiceType("Plumber"), "A-101"))

	req, err := repo.GetRequestByID(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if req == nil {
		t.Fatalf("expected a request, got nil")
	}
	if req.RequestID != id {
		t.Errorf("expected id %v, got %v", id, req.RequestID)
	}
	if req.ResidentID != "resident-123" {
		t.Errorf("expected resident-123, got %s", req.ResidentID)
	}
}

func TestServiceRequestRepository_GetRequestByID_NotFound(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		FROM service_requests
		WHERE request_id = $1
	`)).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	req, err := repo.GetRequestByID(id)
	if err != nil {
		t.Fatalf("expected nil error for not found, got %v", err)
	}
	if req != nil {
		t.Fatalf("expected nil result for not found, got %+v", req)
	}
}

func TestServiceRequestRepository_GetRequestByID_QueryError(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		FROM service_requests
		WHERE request_id = $1
	`)).
		WithArgs(id).
		WillReturnError(errors.New("db error"))

	_, err := repo.GetRequestByID(id)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestServiceRequestRepository_GetRequestByID_ScanError(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	id := uuid.New()
	start := time.Now()
	end := start.Add(45 * time.Minute)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		FROM service_requests
		WHERE request_id = $1
	`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{
			"request_id", "resident_id", "status", "time_slot", "start_time", "end_time", "service_type", "flat_no",
		}).AddRow("not-a-uuid", "resident-123", model.StatusPending, "10:00-10:45", start, end, model.ServiceType("Plumber"), "A-101"))

	_, err := repo.GetRequestByID(id)
	if err == nil {
		t.Fatalf("expected scan error, got nil")
	}
}


func TestServiceRequestRepository_UpdateRequest_Success(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	req := &model.ServiceRequest{
		RequestID:  uuid.New(),
		Status:     model.StatusApproved,
		TimeSlot:   "11:00-11:45",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(45 * time.Minute),
		ServiceType: model.ServiceType("Electrician"),
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE service_requests
		SET status = $1, time_slot = $2, start_time = $3, end_time = $4, service_type = $5
		WHERE request_id = $6
	`)).
		WithArgs(req.Status, req.TimeSlot, req.StartTime, req.EndTime, req.ServiceType, req.RequestID).
		WillReturnResult(sqlmock.NewResult(1, 1)) 

	err := repo.UpdateRequest(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestServiceRequestRepository_UpdateRequest_ExecError(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	req := &model.ServiceRequest{
		RequestID:  uuid.New(),
		Status:     model.StatusPending,
		TimeSlot:   "09:00-09:45",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(45 * time.Minute),
		ServiceType: model.ServiceType("Plumber"),
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE service_requests
		SET status = $1, time_slot = $2, start_time = $3, end_time = $4, service_type = $5
		WHERE request_id = $6
	`)).
		WithArgs(req.Status, req.TimeSlot, req.StartTime, req.EndTime, req.ServiceType, req.RequestID).
		WillReturnError(errors.New("update failed"))

	err := repo.UpdateRequest(req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestServiceRequestRepository_UpdateRequest_NoRowsAffected(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	req := &model.ServiceRequest{
		RequestID:  uuid.New(),
		Status:     model.StatusCancelled,
		TimeSlot:   "15:00-15:45",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(45 * time.Minute),
		ServiceType: model.ServiceType("Electrician"),
	}

	// return 0 rows affected
	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE service_requests
		SET status = $1, time_slot = $2, start_time = $3, end_time = $4, service_type = $5
		WHERE request_id = $6
	`)).
		WithArgs(req.Status, req.TimeSlot, req.StartTime, req.EndTime, req.ServiceType, req.RequestID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	err := repo.UpdateRequest(req)
	if err == nil {
		t.Fatalf("expected error for no rows affected, got nil")
	}
}


func TestServiceRequestRepository_DeleteRequest_Success(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	requestID := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM service_requests WHERE request_id = $1`,
	)).
		WithArgs(requestID).
		WillReturnResult(sqlmock.NewResult(1, 1)) 

	err := repo.DeleteRequest(requestID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestServiceRequestRepository_DeleteRequest_ExecError(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	requestID := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM service_requests WHERE request_id = $1`,
	)).
		WithArgs(requestID).
		WillReturnError(errors.New("delete failed"))

	err := repo.DeleteRequest(requestID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}


func TestServiceRequestRepository_DeleteRequestsByResidentID_Success(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	residentID := "resident-123"

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM service_requests WHERE resident_id = $1`,
	)).
		WithArgs(residentID).
		WillReturnResult(sqlmock.NewResult(1, 2)) 

	err := repo.DeleteRequestsByResidentID(residentID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestServiceRequestRepository_DeleteRequestsByResidentID_ExecError(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	residentID := "resident-456"

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM service_requests WHERE resident_id = $1`,
	)).
		WithArgs(residentID).
		WillReturnError(errors.New("delete failed"))

	err := repo.DeleteRequestsByResidentID(residentID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}


func TestServiceRequestRepository_GetServiceRequestsByStatus_Success(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	userID := "res-123"
	status := model.StatusPending
	requestID := uuid.New()

	rows := sqlmock.NewRows([]string{
		"request_id", "resident_id", "status", "time_slot", "start_time", "end_time", "service_type", "flat_no",
	}).AddRow(
		requestID, userID, status, "slot-1", time.Now(), time.Now().Add(45*time.Minute), model.Electrician, "A-101",
	)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		 FROM service_requests
		 WHERE status = $1 and resident_id = $2`,
	)).
		WithArgs(status, userID).
		WillReturnRows(rows)

	requests := repo.GetServiceRequestsByStatus(userID, status)
	if len(requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(requests))
	}
}

func TestServiceRequestRepository_GetServiceRequestsByStatus_QueryError(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	userID := "res-404"
	status := model.StatusPending

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		 FROM service_requests
		 WHERE status = $1 and resident_id = $2`,
	)).
		WithArgs(status, userID).
		WillReturnError(errors.New("query failed"))

	requests := repo.GetServiceRequestsByStatus(userID, status)
	if len(requests) != 0 {
		t.Fatalf("expected 0 requests on query error, got %d", len(requests))
	}
}

func TestServiceRequestRepository_GetServiceRequestsByStatus_ScanError(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	userID := "res-123"
	status := model.StatusPending


	rows := sqlmock.NewRows([]string{
		"request_id", "resident_id", "status", "time_slot", "start_time", "end_time", "service_type", "flat_no",
	}).AddRow("not-a-uuid", userID, status, "slot-1", time.Now(), time.Now().Add(45*time.Minute), model.Electrician, "A-101")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		 FROM service_requests
		 WHERE status = $1 and resident_id = $2`,
	)).
		WithArgs(status, userID).
		WillReturnRows(rows)

	requests := repo.GetServiceRequestsByStatus(userID, status)
	if len(requests) != 0 {
		t.Fatalf("expected 0 valid requests due to scan error, got %d", len(requests))
	}
}


func TestServiceRequestRepository_GetServiceTypeByID_Success(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	requestID := uuid.New()
	expectedType := model.Plumber

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT service_type 
		 FROM service_requests
		 WHERE request_id = $1`,
	)).
		WithArgs(requestID).
		WillReturnRows(sqlmock.NewRows([]string{"service_type"}).AddRow(expectedType))

	st, err := repo.GetServiceTypeByID(requestID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if st != expectedType {
		t.Fatalf("expected %v, got %v", expectedType, st)
	}
}

func TestServiceRequestRepository_GetServiceTypeByID_NoRows(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	requestID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT service_type 
		 FROM service_requests
		 WHERE request_id = $1`,
	)).
		WithArgs(requestID).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetServiceTypeByID(requestID)
	if err == nil || !strings.Contains(err.Error(), "no request with such id exist") {
		t.Fatalf("expected 'no request with such id exist' error, got %v", err)
	}
}

func TestServiceRequestRepository_GetServiceTypeByID_DBError(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	requestID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT service_type 
		 FROM service_requests
		 WHERE request_id = $1`,
	)).
		WithArgs(requestID).
		WillReturnError(errors.New("db failure"))

	_, err := repo.GetServiceTypeByID(requestID)
	if err == nil || !strings.Contains(err.Error(), "db failure") {
		t.Fatalf("expected db failure error, got %v", err)
	}
}

func TestServiceRequestRepository_GetPendingRequestsByServiceType_Success(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	serviceType := model.Electrician
	rows := sqlmock.NewRows([]string{
		"request_id", "resident_id", "status", "time_slot", "start_time", "end_time", "service_type", "flat_no",
	}).AddRow(uuid.New(), "resident1", model.StatusPending, 1, time.Now(), time.Now().Add(45*time.Minute), serviceType, "A-101").
		AddRow(uuid.New(), "resident2", model.StatusPending, 2, time.Now(), time.Now().Add(45*time.Minute), serviceType, "B-202")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		 FROM service_requests
		 WHERE status = $1 and service_type = $2`,
	)).
		WithArgs(model.StatusPending, serviceType).
		WillReturnRows(rows)

	requests := repo.GetPendingRequestsByServiceType(serviceType)
	if len(requests) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(requests))
	}
}

func TestServiceRequestRepository_GetPendingRequestsByServiceType_Empty(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	serviceType := model.Electrician
	rows := sqlmock.NewRows([]string{
		"request_id", "resident_id", "status", "time_slot", "start_time", "end_time", "service_type", "flat_no",
	})

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		 FROM service_requests
		 WHERE status = $1 and service_type = $2`,
	)).
		WithArgs(model.StatusPending, serviceType).
		WillReturnRows(rows)

	requests := repo.GetPendingRequestsByServiceType(serviceType)
	if len(requests) != 0 {
		t.Fatalf("expected 0 requests, got %d", len(requests))
	}
}

func TestServiceRequestRepository_GetPendingRequestsByServiceType_DBError(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	serviceType := model.Electrician
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		 FROM service_requests
		 WHERE status = $1 and service_type = $2`,
	)).
		WithArgs(model.StatusPending, serviceType).
		WillReturnError(errors.New("db failure"))

	requests := repo.GetPendingRequestsByServiceType(serviceType)
	if requests != nil {
		t.Fatalf("expected nil due to db failure, got %+v", requests)
	}
}

func TestServiceRequestRepository_GetPendingRequestsByServiceType_ScanError(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	serviceType := model.Electrician

	rows := sqlmock.NewRows([]string{
		"request_id", "resident_id", "status", "time_slot", "start_time", "end_time", "service_type", "flat_no",
	}).AddRow(123, "resident1", model.StatusPending, 1, time.Now(), time.Now().Add(45*time.Minute), serviceType, "A-101")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		 FROM service_requests
		 WHERE status = $1 and service_type = $2`,
	)).
		WithArgs(model.StatusPending, serviceType).
		WillReturnRows(rows)

	requests := repo.GetPendingRequestsByServiceType(serviceType)
	if requests != nil {
		t.Fatalf("expected nil due to scan error, got %+v", requests)
	}
}

func TestServiceRequestRepository_GetApprovedRequestsByServiceType_Success(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	serviceType := model.Plumber
	rows := sqlmock.NewRows([]string{
		"request_id", "resident_id", "status", "time_slot", "start_time", "end_time", "service_type", "flat_no",
	}).AddRow(uuid.New(), "resident1", model.StatusApproved, 3, time.Now(), time.Now().Add(45*time.Minute), serviceType, "C-303").
		AddRow(uuid.New(), "resident2", model.StatusApproved, 4, time.Now(), time.Now().Add(45*time.Minute), serviceType, "D-404")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		 FROM service_requests
		 WHERE status = $1 and service_type = $2`,
	)).
		WithArgs(model.StatusApproved, serviceType).
		WillReturnRows(rows)

	requests := repo.GetApprovedRequestsByServiceType(serviceType)
	if len(requests) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(requests))
	}
}

func TestServiceRequestRepository_GetApprovedRequestsByServiceType_Empty(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	serviceType := model.Plumber
	rows := sqlmock.NewRows([]string{
		"request_id", "resident_id", "status", "time_slot", "start_time", "end_time", "service_type", "flat_no",
	})

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		 FROM service_requests
		 WHERE status = $1 and service_type = $2`,
	)).
		WithArgs(model.StatusApproved, serviceType).
		WillReturnRows(rows)

	requests := repo.GetApprovedRequestsByServiceType(serviceType)
	if len(requests) != 0 {
		t.Fatalf("expected 0 requests, got %d", len(requests))
	}
}

func TestServiceRequestRepository_GetApprovedRequestsByServiceType_DBError(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	serviceType := model.Plumber
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		 FROM service_requests
		 WHERE status = $1 and service_type = $2`,
	)).
		WithArgs(model.StatusApproved, serviceType).
		WillReturnError(errors.New("db failure"))

	requests := repo.GetApprovedRequestsByServiceType(serviceType)
	if requests != nil {
		t.Fatalf("expected nil due to db failure, got %+v", requests)
	}
}

func TestServiceRequestRepository_GetApprovedRequestsByServiceType_ScanError(t *testing.T) {
	repo, mock, cleanup := newServiceRequestRepo(t)
	defer cleanup()

	serviceType := model.Plumber
	rows := sqlmock.NewRows([]string{
		"request_id", "resident_id", "status", "time_slot", "start_time", "end_time", "service_type", "flat_no",
	}).AddRow(123, "resident1", model.StatusApproved, 1, time.Now(), time.Now().Add(45*time.Minute), serviceType, "C-505")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT request_id, resident_id, status, time_slot, start_time, end_time, service_type, flat_no
		 FROM service_requests
		 WHERE status = $1 and service_type = $2`,
	)).
		WithArgs(model.StatusApproved, serviceType).
		WillReturnRows(rows)

	requests := repo.GetApprovedRequestsByServiceType(serviceType)
	if requests != nil {
		t.Fatalf("expected nil due to scan error, got %+v", requests)
	}
}
