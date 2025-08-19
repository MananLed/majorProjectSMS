package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/repository"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
	"github.com/google/uuid"
)

type ServiceRequestServiceInterface interface {
	BookServiceRequest(req model.ServiceRequest) error
	RescheduleServiceRequest(userID string, requestID uuid.UUID, newSlot utils.TimeSlot, newService model.ServiceType) error
	CancelServiceRequest(userID string, requestID uuid.UUID) error
	GetServiceRequestsByStatus(userID string, status model.Status) []model.ServiceRequest
	GetAvailableTimeSlots(service model.ServiceType) []utils.TimeSlot
	GetServiceTypeByID(requestID uuid.UUID) (model.ServiceType, error)
	GetPendingRequestsByServiceType(serviceType model.ServiceType) []model.ServiceRequest
	GetApprovedRequestsByServiceType(serviceType model.ServiceType) []model.ServiceRequest
	ApproveServiceRequest(requestID uuid.UUID) error
	DeleteServiceRequestByID(ctx context.Context) error
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

func (s *ServiceRequestService) GetServiceRequestsByStatus(userID string, status model.Status) []model.ServiceRequest {
	return s.Repo.GetServiceRequestsByStatus(userID, status)
}

func (s *ServiceRequestService) GetAvailableTimeSlots(service model.ServiceType) []utils.TimeSlot {
	booked := map[string]bool{}

	requests, err := s.Repo.GetAllRequests()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return utils.GenerateTimeSlots()
	}

	for _, r := range requests {
		if r.ServiceType == service && r.Status != model.StatusCancelled {
			booked[r.TimeSlot] = true
		}
	}

	var available []utils.TimeSlot
	for _, slot := range utils.GenerateTimeSlots() {
		if !booked[slot.Label] {
			available = append(available, slot)
		}
	}
	return available
}

func (s *ServiceRequestService) GetServiceTypeByID(requestID uuid.UUID) (model.ServiceType, error) {
	return s.Repo.GetServiceTypeByID(requestID)
}

func (s *ServiceRequestService) GetPendingRequestsByServiceType(serviceType model.ServiceType) []model.ServiceRequest {
	return s.Repo.GetPendingRequestsByServiceType(serviceType)
}

func (s *ServiceRequestService) GetApprovedRequestsByServiceType(serviceType model.ServiceType) []model.ServiceRequest {
	return s.Repo.GetApprovedRequestsByServiceType(serviceType)
}

func (s *ServiceRequestService) ApproveServiceRequest(requestID uuid.UUID) error {
	req, err := s.Repo.GetRequestByID(requestID)
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	req.Status = model.StatusApproved

	return s.Repo.UpdateRequest(req)
}

func (s *ServiceRequestService) DeleteServiceRequestByID(ctx context.Context) error {
	residentID := ctx.Value(utils.UserIDKey).(string)
	return s.Repo.DeleteRequestsByResidentID(residentID)
}