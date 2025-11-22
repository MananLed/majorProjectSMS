//go:generate mockgen -source=credential_service.go -destination=../mocks/credential_mock_service.go -package=mocks
package service

import (
	"context"
	"errors"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/repository"
	"github.com/MananLed/majorProjectSMS/internal/utils"
)

type CredentialServiceInterface interface {
	DeleteOfficerCredentials(ctx context.Context, officerID string) error
	DeleteResidentCredentials(ctx context.Context, residentID string) error
}

type CredentialService struct {
	Repo repository.CredentialRepositoryInterface
}

func NewCredentialService(r repository.CredentialRepositoryInterface) *CredentialService {
	return &CredentialService{Repo: r}
}

func (s *CredentialService) DeleteOfficerCredentials(ctx context.Context, officerID string) error {

	user, err := utils.GetUserFromContext(ctx)

	if err != nil || user.Role != model.RoleAdmin {
		return errors.New("unauthorized: only admin can delete credentials")
	}

	return s.Repo.DeleteUserByIDAndRole(officerID, model.RoleOfficer)
}

func (s *CredentialService) DeleteResidentCredentials(ctx context.Context, residentID string) error {

	user, err := utils.GetUserFromContext(ctx)

	if err != nil || user.Role != model.RoleAdmin {
		return errors.New("unauthorized: only admin can delete credentials")
	}

	return s.Repo.DeleteUserByIDAndRole(residentID, model.RoleResident)
}
