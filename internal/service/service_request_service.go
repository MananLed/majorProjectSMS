//go:generate mockgen -source=service_request_service.go -destination=../mocks/service_request_mock_service.go -package=mocks
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/repository"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/google/uuid"
)

type ServiceRequestServiceInterface interface {
	BookServiceRequest(req model.ServiceRequest) error
	RescheduleServiceRequest(userID string, requestID uuid.UUID, newSlot utils.TimeSlot, newService model.ServiceType) error
	CancelServiceRequest(userID string, requestID uuid.UUID) error
	GetServiceRequestsByStatus(userID string, status model.Status) ([]model.ServiceRequest, error)
	GetAvailableTimeSlots(service model.ServiceType) ([]utils.TimeSlot, error)
	GetServiceTypeByID(requestID uuid.UUID) (model.ServiceType, error)
	GetPendingRequestsByServiceType(serviceType model.ServiceType) ([]model.ServiceRequest, error)
	GetApprovedRequestsByServiceType(serviceType model.ServiceType) ([]model.ServiceRequest, error)
	ApproveServiceRequest(requestID uuid.UUID, assignedTo string) error
	DeleteServiceRequestByID(ctx context.Context) error
	CompleteServiceRequest(requestID uuid.UUID) error
	GetCompletedRequestsByServiceType(serviceType model.ServiceType) ([]model.ServiceRequest, error)
}

type ServiceRequestService struct {
	Repo repository.ServiceRequestRepositoryInterface
}

func NewServiceRequestService(repo repository.ServiceRequestRepositoryInterface) *ServiceRequestService {
	return &ServiceRequestService{Repo: repo}
}

func (s *ServiceRequestService) normalizeTime(t time.Time) time.Time {
	return time.Date(1, 1, 1, t.Hour(), t.Minute(), 0, 0, time.Local)
}

func (s *ServiceRequestService) BookServiceRequest(req model.ServiceRequest) error {
	req.StartTime = s.normalizeTime(req.StartTime)
	req.EndTime = s.normalizeTime(req.EndTime)
	return s.Repo.CreateRequest(&req)
}

func (s *ServiceRequestService) RescheduleServiceRequest(userID string, requestID uuid.UUID, newSlot utils.TimeSlot, newService model.ServiceType) error {

	newStart := s.normalizeTime(newSlot.StartTime)
	newEnd := s.normalizeTime(newSlot.EndTime)

	var updatedRequest model.ServiceRequest

	updatedRequest.TimeSlot = fmt.Sprintf("%s - %s", newStart.Format("3:04 PM"), newEnd.Format("3:04 PM"))
	updatedRequest.StartTime = newStart
	updatedRequest.EndTime = newEnd
	updatedRequest.ServiceType = newService
	updatedRequest.Status = model.StatusPending
	updatedRequest.ResidentID = userID
	updatedRequest.RequestID = requestID

	return s.Repo.UpdateRequest(&updatedRequest)
}

func (s *ServiceRequestService) CancelServiceRequest(userID string, requestID uuid.UUID) error {
	return s.Repo.DeleteRequest(requestID)
}

func (s *ServiceRequestService) GetServiceRequestsByStatus(userID string, status model.Status) ([]model.ServiceRequest, error) {
	return s.Repo.GetServiceRequestsByStatus(userID, status)
}

func (s *ServiceRequestService) GetAvailableTimeSlots(service model.ServiceType) ([]utils.TimeSlot, error) {
	booked := map[string]bool{}
	now := time.Now()
	formattedDate := now.Format("02-01-2006")

	requests, err := s.Repo.GetRequestsByServiceTypeAndStatus(nil, service, model.StatusPending)
	if err != nil {
		return nil, err
	}

	requestsApproved, err := s.Repo.GetRequestsByServiceTypeAndStatus(nil, service, model.StatusApproved)
	if err != nil {
		return nil, err
	}
	requests = append(requests, requestsApproved...)

	requestsCompleted, err := s.Repo.GetRequestsByServiceTypeAndStatus(nil, service, model.StatusCompleted)
	if err != nil {
		return nil, err
	}
	requests = append(requests, requestsCompleted...)

	for _, r := range requests {
		if r.Date == formattedDate{
			booked[r.TimeSlot] = true
		}
	}

	var available []utils.TimeSlot
	for _, slot := range utils.GenerateTimeSlots() {
		if !booked[slot.Label] {
			available = append(available, slot)
		}
	}
	return available, nil
}

func (s *ServiceRequestService) GetServiceTypeByID(requestID uuid.UUID) (model.ServiceType, error) {
	return s.Repo.GetServiceTypeByID(requestID)
}

func (s *ServiceRequestService) GetPendingRequestsByServiceType(user *model.User, serviceType model.ServiceType) ([]model.ServiceRequest, error) {
	return s.Repo.GetRequestsByServiceTypeAndStatus(user, serviceType, model.StatusPending)
}

func (s *ServiceRequestService) GetApprovedRequestsByServiceType(user *model.User, serviceType model.ServiceType) ([]model.ServiceRequest, error) {
	return s.Repo.GetRequestsByServiceTypeAndStatus(user, serviceType, model.StatusApproved)
}

func (s *ServiceRequestService) GetCompletedRequestsByServiceType(user *model.User, serviceType model.ServiceType) ([]model.ServiceRequest, error) {
	return s.Repo.GetRequestsByServiceTypeAndStatus(user, serviceType, model.StatusCompleted)
}

func (s *ServiceRequestService) ApproveServiceRequest(requestID uuid.UUID, assignedTo string) error {

	req := &model.ServiceRequest{
		RequestID: requestID,
		Status: model.StatusApproved,
		AssignedTo: assignedTo,
	}

	return s.Repo.UpdateRequest(req)
}

func (s *ServiceRequestService) CompleteServiceRequest(requestID uuid.UUID) error{
	req := &model.ServiceRequest{
		RequestID: requestID,
		Status: model.StatusCompleted,
	}
	return s.Repo.UpdateRequest(req)
}

func (s *ServiceRequestService) DeleteServiceRequestByID(ctx context.Context) error {
	residentID := ctx.Value(utils.UserIDKey).(string)
	return s.Repo.DeleteRequestsByResidentID(residentID)
}