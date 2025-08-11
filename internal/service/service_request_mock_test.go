package service

import (
	"context"
	"testing"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
)

type MockServiceRequestRepo struct {
	requests map[string]model.ServiceRequest
}

func (m *MockServiceRequestRepo) LoadRequests() ([]model.ServiceRequest, error) {
	var serviceRequest []model.ServiceRequest

	for i := range m.requests {
		serviceRequest = append(serviceRequest, m.requests[i])
	}

	return serviceRequest, nil
}

func (m *MockServiceRequestRepo) SaveRequests(serviceRequests []model.ServiceRequest) error {
	mp := make(map[string]model.ServiceRequest)

	for i := range serviceRequests {
		mp[serviceRequests[i].RequestID] = serviceRequests[i]
	}
	m.requests = mp
	return nil
}

func TestBookServiceRequest(t *testing.T) {
	mockRepo := &MockServiceRequestRepo{requests: make(map[string]model.ServiceRequest)}
	service := NewServiceRequestService(mockRepo)

	request := model.ServiceRequest{
		RequestID:   "req-123",
		ResidentID:  "res-1",
		ServiceType: model.Plumber,
		StartTime:   time.Date(0, 0, 0, 9, 0, 0, 0, time.Local),
		EndTime:     time.Date(0, 0, 0, 10, 0, 0, 0, time.Local),
		Status:      model.StatusPending,
	}

	err := service.BookServiceRequest(request)
	if err != nil {
		t.Errorf("expected booking to succeed, got error: %v", err)
	}

	err = service.BookServiceRequest(request)
	if err == nil || err.Error() != "time slot already booked" {
		t.Errorf("expected booking conflict error, got: %v", err)
	}
}

func TestCancelServiceRequest(t *testing.T) {
	mockRepo := &MockServiceRequestRepo{requests: make(map[string]model.ServiceRequest)}
	service := NewServiceRequestService(mockRepo)

	req := model.ServiceRequest{
		RequestID:   "cancel-1",
		ResidentID:  "res-1",
		ServiceType: model.Electrician,
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		Status:      model.StatusPending,
	}

	mockRepo.requests[req.RequestID] = req

	err := service.CancelServiceRequest("res-1", "cancel-1")
	if err != nil {
		t.Errorf("expected cancel to succeed, got: %v", err)
	}

	err = service.CancelServiceRequest("res-1", "cancel-1")
	if err == nil || err.Error() != "request already cancelled" {
		t.Errorf("expected already cancelled error, got: %v", err)
	}
}

func TestDeleteServiceRequestByID(t *testing.T) {
	mockRepo := &MockServiceRequestRepo{requests: make(map[string]model.ServiceRequest)}
	service := NewServiceRequestService(mockRepo)

	req1 := model.ServiceRequest{RequestID: "req-1", ResidentID: "res-1"}
	req2 := model.ServiceRequest{RequestID: "req-2", ResidentID: "res-2"}

	mockRepo.requests[req1.RequestID] = req1
	mockRepo.requests[req2.RequestID] = req2

	ctx := context.WithValue(context.Background(), utils.UserIDKey, req1.RequestID)

	err := service.DeleteServiceRequestByID(ctx)
	if err != nil {
		t.Errorf("expected delete to succeed, got: %v", err)
	}

	if _, exists := mockRepo.requests["req-1"]; exists {
		t.Errorf("expected req-1 to be deleted")
	}
}
