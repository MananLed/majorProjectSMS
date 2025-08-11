package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/MananLed/majorProjectSMS/constants"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
)

type ServiceRequestRepository struct {
	mu sync.Mutex
}

type ServiceRequestRepositoryInterface interface {
	LoadRequests() ([]model.ServiceRequest, error)
	SaveRequests([]model.ServiceRequest) error
}

func NewServiceRequestRepository() *ServiceRequestRepository {
	return &ServiceRequestRepository{}
}

func (r *ServiceRequestRepository) LoadRequests() ([]model.ServiceRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(string(constants.ServiceRequestDataPath))
	if err != nil {
		if os.IsNotExist(err) {
			return []model.ServiceRequest{}, nil
		}
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}

	var requests []model.ServiceRequest
	if err := json.Unmarshal(data, &requests); err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}

	return requests, nil
}

func (r *ServiceRequestRepository) SaveRequests(requests []model.ServiceRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := json.MarshalIndent(requests, "", "  ")
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	return os.WriteFile(string(constants.ServiceRequestDataPath), data, 0644)
}
