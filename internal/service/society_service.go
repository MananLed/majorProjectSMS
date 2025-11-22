//go:generate mockgen -source=society_service.go -destination=../mocks/society_mock_service.go -package=mocks
package service

import (
	"context"
	"errors"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/repository"
	"github.com/MananLed/majorProjectSMS/internal/utils"
)

type SocietyServiceInterface interface{
	GetAllResidents(ctx context.Context) ([]model.User, error)
	GetAllOfficers(ctx context.Context) ([]model.User, error)
}

type SocietyService struct {
	SocietyRepo repository.SocietyRepositoryInterface
}

func NewSocietyService(repo repository.SocietyRepositoryInterface) *SocietyService {
	return &SocietyService{SocietyRepo: repo}
}

func (s *SocietyService) GetAllResidents(ctx context.Context) ([]model.User, error) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if user.Role == model.RoleResident {
		return nil, errors.New("unauthorized access: only admins can view residents")
	}
	return s.SocietyRepo.GetAllResidents()
}

func (s *SocietyService) GetAllOfficers(ctx context.Context) ([]model.User, error) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if user.Role == model.RoleResident {
		return nil, errors.New("unauthorized access: only admins can view officers")
	}
	return s.SocietyRepo.GetAllOfficers()
}
