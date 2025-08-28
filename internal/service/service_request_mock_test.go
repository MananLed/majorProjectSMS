package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/google/uuid"
)


type MockServiceRequestRepo struct {
	requests map[uuid.UUID]model.ServiceRequest
	saveErr  error
	getErr   error
}

func (m *MockServiceRequestRepo) CreateRequest(req *model.ServiceRequest) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	if m.requests == nil {
		m.requests = make(map[uuid.UUID]model.ServiceRequest)
	}
	m.requests[req.RequestID] = *req
	return nil
}

func (m *MockServiceRequestRepo) GetAllRequests() ([]model.ServiceRequest, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var result []model.ServiceRequest
	for _, v := range m.requests {
		result = append(result, v)
	}
	return result, nil
}

func (m *MockServiceRequestRepo) GetRequestByID(requestID uuid.UUID) (*model.ServiceRequest, error) {
	req, ok := m.requests[requestID]
	if !ok {
		return nil, errors.New("not found")
	}
	return &req, nil
}

func (m *MockServiceRequestRepo) UpdateRequest(req *model.ServiceRequest) error {
	if _, ok := m.requests[req.RequestID]; !ok {
		return errors.New("not found")
	}
	m.requests[req.RequestID] = *req
	return nil
}

func (m *MockServiceRequestRepo) DeleteRequest(requestID uuid.UUID) error {
	delete(m.requests, requestID)
	return nil
}

func (m *MockServiceRequestRepo) DeleteRequestsByResidentID(residentID string) error {
	for id, r := range m.requests {
		if r.ResidentID == residentID {
			delete(m.requests, id)
		}
	}
	return nil
}

func (m *MockServiceRequestRepo) GetServiceRequestsByStatus(userID string, status model.Status) []model.ServiceRequest {
	var result []model.ServiceRequest
	for _, r := range m.requests {
		if r.ResidentID == userID && r.Status == status {
			result = append(result, r)
		}
	}
	return result
}

func (m *MockServiceRequestRepo) GetServiceTypeByID(requestID uuid.UUID) (model.ServiceType, error) {
	req, ok := m.requests[requestID]
	if !ok {
		return "", errors.New("not found")
	}
	return req.ServiceType, nil
}

func (m *MockServiceRequestRepo) GetPendingRequestsByServiceType(serviceType model.ServiceType) []model.ServiceRequest {
	var result []model.ServiceRequest
	for _, r := range m.requests {
		if r.ServiceType == serviceType && r.Status == model.StatusPending {
			result = append(result, r)
		}
	}
	return result
}

func (m *MockServiceRequestRepo) GetApprovedRequestsByServiceType(serviceType model.ServiceType) []model.ServiceRequest {
	var result []model.ServiceRequest
	for _, r := range m.requests {
		if r.ServiceType == serviceType && r.Status == model.StatusApproved {
			result = append(result, r)
		}
	}
	return result
}


func TestBookServiceRequest_Success(t *testing.T) {
	mockRepo := &MockServiceRequestRepo{requests: make(map[uuid.UUID]model.ServiceRequest)}
	service := NewServiceRequestService(mockRepo)

	req := model.ServiceRequest{
		RequestID:  uuid.New(),
		ResidentID: "resident1",
		Status:     model.StatusPending,
		TimeSlot:   "10:00 - 10:45",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(45 * time.Minute),
		ServiceType: model.Electrician,
		Flat:       "A-101",
	}

	err := service.BookServiceRequest(req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(mockRepo.requests) != 1 {
		t.Errorf("expected 1 request saved, got %d", len(mockRepo.requests))
	}
}

func TestRescheduleServiceRequest(t *testing.T) {
	mockRepo := &MockServiceRequestRepo{requests: make(map[uuid.UUID]model.ServiceRequest)}
	service := NewServiceRequestService(mockRepo)

	reqID := uuid.New()
	mockRepo.requests[reqID] = model.ServiceRequest{
		RequestID:  reqID,
		ResidentID: "resident1",
		Status:     model.StatusPending,
		ServiceType: model.Plumber,
	}

	newSlot := utils.TimeSlot{
		Label:     "11:00 - 11:45",
		StartTime: time.Date(1, 1, 1, 11, 0, 0, 0, time.Local),
		EndTime:   time.Date(1, 1, 1, 11, 45, 0, 0, time.Local),
	}

	err := service.RescheduleServiceRequest("resident1", reqID, newSlot, model.Plumber)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	updated := mockRepo.requests[reqID]
	if updated.TimeSlot != "11:00 AM - 11:45 AM" {
		t.Errorf("expected timeslot to be '11:00 - 11:45', got %v", updated.TimeSlot)
	}
}

func TestCancelServiceRequest(t *testing.T) {
	mockRepo := &MockServiceRequestRepo{requests: make(map[uuid.UUID]model.ServiceRequest)}
	service := NewServiceRequestService(mockRepo)

	reqID := uuid.New()
	mockRepo.requests[reqID] = model.ServiceRequest{
		RequestID:  reqID,
		ResidentID: "resident1",
		Status:     model.StatusPending,
	}

	err := service.CancelServiceRequest("resident1", reqID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if _, exists := mockRepo.requests[reqID]; exists {
		t.Errorf("expected request to be deleted")
	}
}

func TestApproveServiceRequest(t *testing.T) {
	mockRepo := &MockServiceRequestRepo{requests: make(map[uuid.UUID]model.ServiceRequest)}
	service := NewServiceRequestService(mockRepo)

	reqID := uuid.New()
	mockRepo.requests[reqID] = model.ServiceRequest{
		RequestID:  reqID,
		ResidentID: "resident1",
		Status:     model.StatusPending,
		ServiceType: model.Plumber,
	}

	err := service.ApproveServiceRequest(reqID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if mockRepo.requests[reqID].Status != model.StatusApproved {
		t.Errorf("expected status Approved, got %v", mockRepo.requests[reqID].Status)
	}
}

func TestDeleteServiceRequestByID(t *testing.T) {
	mockRepo := &MockServiceRequestRepo{requests: make(map[uuid.UUID]model.ServiceRequest)}
	service := NewServiceRequestService(mockRepo)

	reqID := uuid.New()
	mockRepo.requests[reqID] = model.ServiceRequest{
		RequestID:  reqID,
		ResidentID: "resident1",
		Status:     model.StatusPending,
	}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, "resident1")

	err := service.DeleteServiceRequestByID(ctx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(mockRepo.requests) != 0 {
		t.Errorf("expected requests map to be empty")
	}
}

func TestBookServiceRequest_Failure(t *testing.T) {
    mockRepo := &MockServiceRequestRepo{saveErr: errors.New("db error")}
    service := NewServiceRequestService(mockRepo)

    req := model.ServiceRequest{RequestID: uuid.New()}
    err := service.BookServiceRequest(req)
    if err == nil {
        t.Errorf("expected error, got nil")
    }
}

func TestRescheduleServiceRequest_NotFound(t *testing.T) {
    mockRepo := &MockServiceRequestRepo{requests: make(map[uuid.UUID]model.ServiceRequest)}
    service := NewServiceRequestService(mockRepo)

    newSlot := utils.TimeSlot{
        Label:     "11:00 - 11:45",
        StartTime: time.Date(1, 1, 1, 11, 0, 0, 0, time.Local),
        EndTime:   time.Date(1, 1, 1, 11, 45, 0, 0, time.Local),
    }

    err := service.RescheduleServiceRequest("resident1", uuid.New(), newSlot, model.Plumber)
    if err == nil {
        t.Errorf("expected error, got nil")
    }
}

func TestGetAvailableTimeSlots_ErrorFromRepo(t *testing.T) {
    mockRepo := &MockServiceRequestRepo{getErr: errors.New("db error")}
    service := NewServiceRequestService(mockRepo)

    slots := service.GetAvailableTimeSlots(model.Plumber)
    if len(slots) == 0 {
        t.Errorf("expected fallback slots, got none")
    }
}

func TestGetAvailableTimeSlots_BookedSlotExcluded(t *testing.T) {
    mockRepo := &MockServiceRequestRepo{requests: make(map[uuid.UUID]model.ServiceRequest)}
    service := NewServiceRequestService(mockRepo)

    bookedSlot := utils.GenerateTimeSlots()[0]
    reqID := uuid.New()
    mockRepo.requests[reqID] = model.ServiceRequest{
        RequestID:   reqID,
        ServiceType: model.Plumber,
        Status:      model.StatusPending,
        TimeSlot:    bookedSlot.Label,
    }

    slots := service.GetAvailableTimeSlots(model.Plumber)
    for _, s := range slots {
        if s.Label == bookedSlot.Label {
            t.Errorf("expected booked slot to be excluded")
        }
    }
}

func TestApproveServiceRequest_NotFound(t *testing.T) {
    mockRepo := &MockServiceRequestRepo{requests: make(map[uuid.UUID]model.ServiceRequest)}
    service := NewServiceRequestService(mockRepo)

    err := service.ApproveServiceRequest(uuid.New())
    if err == nil {
        t.Errorf("expected error, got nil")
    }
}

func TestGetServiceTypeByID(t *testing.T) {
    mockRepo := &MockServiceRequestRepo{requests: make(map[uuid.UUID]model.ServiceRequest)}
    service := NewServiceRequestService(mockRepo)

    reqID := uuid.New()
    mockRepo.requests[reqID] = model.ServiceRequest{RequestID: reqID, ServiceType: model.Electrician}

    st, err := service.GetServiceTypeByID(reqID)
    if err != nil || st != model.Electrician {
        t.Errorf("expected Electrician, got %v err %v", st, err)
    }

    _, err = service.GetServiceTypeByID(uuid.New())
    if err == nil {
        t.Errorf("expected error for missing request")
    }
}
