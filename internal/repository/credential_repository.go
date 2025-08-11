package repository

import (
	"errors"
	"fmt"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
)

type CredentialRepository struct {
	UserRepository
}

type CredentialRepositoryInterface interface {
	DeleteUserByIDAndRole(id string, role model.UserRole) error
}

func (r *CredentialRepository) DeleteUserByIDAndRole(id string, role model.UserRole) error {
	users, err := r.LoadUsers()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	var updated []model.User
	found := false

	for _, u := range users {
		if u.ID == id && u.Role == role {
			found = true
			continue
		}
		updated = append(updated, u)
	}

	if !found {
		return errors.New("user not found")
	}

	return r.SaveUsers(updated)
}
