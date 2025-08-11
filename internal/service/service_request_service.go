package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/repository"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
)

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
	allRequests, err := s.Repo.LoadRequests()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	req.StartTime = s.normalizeTime(req.StartTime)
	req.EndTime = s.normalizeTime(req.EndTime)

	for _, r := range allRequests {
		if r.ServiceType == req.ServiceType &&
			s.normalizeTime(r.StartTime).Equal(req.StartTime) &&
			s.normalizeTime(r.EndTime).Equal(req.EndTime) &&
			r.Status != model.StatusCancelled {
			return errors.New("time slot already booked")
		}
	}

	allRequests = append(allRequests, req)
	return s.Repo.SaveRequests(allRequests)
}

func (s *ServiceRequestService) RescheduleServiceRequest(userID, requestID string, newSlot utils.TimeSlot, newService model.ServiceType) error {
	allRequests, err := s.Repo.LoadRequests()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	newStart := s.normalizeTime(newSlot.StartTime)
	newEnd := s.normalizeTime(newSlot.EndTime)

	updated := false
	for i, r := range allRequests {
		if r.RequestID == requestID && r.ResidentID == userID {
			if r.Status == model.StatusCancelled || r.Status == model.StatusApproved {
				return errors.New("cannot reschedule a cancelled or approved request")
			}

			allRequests[i].TimeSlot = fmt.Sprintf("%s - %s", newStart.Format("3:04 PM"), newEnd.Format("3:04 PM"))
			allRequests[i].StartTime = newStart
			allRequests[i].EndTime = newEnd
			allRequests[i].ServiceType = newService
			allRequests[i].Status = model.StatusPending
			updated = true
			break
		}
	}

	if !updated {
		return errors.New("request not found")
	}

	return s.Repo.SaveRequests(allRequests)
}

func (s *ServiceRequestService) CancelServiceRequest(userID, requestID string) error {
	allRequests, err := s.Repo.LoadRequests()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	cancelled := false
	for i, r := range allRequests {
		if r.RequestID == requestID && r.ResidentID == userID {
			if r.Status == model.StatusCancelled {
				return errors.New("request already cancelled")
			} else if r.Status == model.StatusApproved {
				return errors.New("request is already approved")
			}
			allRequests[i].Status = model.StatusCancelled
			cancelled = true
			break
		}
	}

	if !cancelled {
		return errors.New("request not found")
	}

	return s.Repo.SaveRequests(allRequests)
}

func (s *ServiceRequestService) GetServiceRequestsByStatus(userID string, status model.Status) []model.ServiceRequest {
	requests, err := s.Repo.LoadRequests()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil
	}

	var filtered []model.ServiceRequest
	for _, r := range requests {
		if r.ResidentID == userID && r.Status == status {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func (s *ServiceRequestService) GetAvailableTimeSlots(service model.ServiceType) []utils.TimeSlot {
	booked := map[string]bool{}

	requests, err := s.Repo.LoadRequests()
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

func (s *ServiceRequestService) GetServiceTypeByID(requestID string) (model.ServiceType, error) {
	requests, err := s.Repo.LoadRequests()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return "", errors.New("request with such ID is not present")
	}

	for _, request := range requests {
		if request.RequestID == requestID {
			fmt.Println(request.RequestID, requestID)
			return request.ServiceType, nil
		}
	}
	return "", errors.New("request with such ID is not present")
}

func (s *ServiceRequestService) GetPendingRequestsByServiceType(serviceType model.ServiceType) []model.ServiceRequest {
	requests, err := s.Repo.LoadRequests()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil
	}

	var filtered []model.ServiceRequest
	for _, r := range requests {
		if r.ServiceType == serviceType && r.Status == model.StatusPending {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func (s *ServiceRequestService) GetApprovedRequestsByServiceType(serviceType model.ServiceType) []model.ServiceRequest {
	requests, err := s.Repo.LoadRequests()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil
	}

	var filtered []model.ServiceRequest
	for _, r := range requests {
		if r.ServiceType == serviceType && r.Status == model.StatusApproved {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func (s *ServiceRequestService) ApproveServiceRequest(requestID string) error {
	requests, err := s.Repo.LoadRequests()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	updated := false
	for i, r := range requests {
		if r.RequestID == requestID {
			if r.Status != model.StatusPending {
				return errors.New("only pending requests can be approved")
			}
			requests[i].Status = model.StatusApproved
			updated = true
			break
		}
	}

	if !updated {
		return errors.New("request ID not found")
	}

	return s.Repo.SaveRequests(requests)
}

func (s *ServiceRequestService) DeleteServiceRequestByID(ctx context.Context) error {
	requestID := ctx.Value(utils.UserIDKey).(string)
	requests, err := s.Repo.LoadRequests()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}
	var validRequests []model.ServiceRequest
	for _, r := range requests {
		if r.RequestID == requestID {
			continue
		}
		validRequests = append(validRequests, r)
	}

	return s.Repo.SaveRequests(validRequests)
}
