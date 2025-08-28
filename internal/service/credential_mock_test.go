package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
)

type MockCredentialRepo struct {
	users map[string]model.User
}

func (m *MockCredentialRepo) DeleteUserByIDAndRole(id string, role model.UserRole) error {
	user, exists := m.users[id]
	if !exists {
		return errors.New("user not found")
	}
	if user.Role != role {
		return errors.New("role mismatch")
	}

	delete(m.users, id)
	return nil
}

func TestDeleteResident(t *testing.T) {
	mockRepo := &MockCredentialRepo{
		users: make(map[string]model.User),
	}

	admin := model.User{
		ID:       "admin@gmail.com",
		Password: "admin",
		Role:     model.RoleAdmin,
	}
	resident := model.User{
		ID:       "resident@gmail.com",
		Password: "resident",
		Role:     model.RoleResident,
	}

	mockRepo.users[admin.ID] = admin
	mockRepo.users[resident.ID] = resident

	service := NewCredentialService(mockRepo)

	ctx := context.WithValue(context.Background(), utils.UserRoleKey, string(model.RoleAdmin))
	ctx = context.WithValue(ctx, utils.UserEmailKey, "dsfsl")
	ctx = context.WithValue(ctx, utils.UserFlatKey, "xxx")
	ctx = context.WithValue(ctx, utils.UserIDKey, "dfjld")

	err := service.DeleteResidentCredentials(ctx, resident.ID)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if _, exists := mockRepo.users[resident.ID]; exists {
		t.Errorf("resident user should have been deleted")
	}
}

func TestDeleteOfficerCredentials(t *testing.T) {
	mockRepo := &MockCredentialRepo{
		users: make(map[string]model.User),
	}

	admin := model.User{
		ID:    "admin@gmail.com",
		Password: "admin",
		Role:  model.RoleAdmin,
	}
	officer := model.User{
		ID:    "officer@gmail.com",
		Password: "officer",
		Role:  model.RoleOfficer,
	}

	mockRepo.users[admin.ID] = admin
	mockRepo.users[officer.ID] = officer

	service := NewCredentialService(mockRepo)

	ctx := context.WithValue(context.Background(), utils.UserRoleKey, string(model.RoleAdmin))
	ctx = context.WithValue(ctx, utils.UserEmailKey, "dsfsl")
	ctx = context.WithValue(ctx, utils.UserFlatKey, "xxx")
	ctx = context.WithValue(ctx, utils.UserIDKey, "dfjld")

	err := service.DeleteOfficerCredentials(ctx, officer.ID)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if _, exists := mockRepo.users[officer.ID]; exists {
		t.Errorf("officer user should have been deleted")
	}
}

func TestDeleteResidentCredentialsUnauthorized(t *testing.T) {
	mockRepo := &MockCredentialRepo{users: make(map[string]model.User)}

	mockRepo.users["resident@gmail.com"] = model.User{
		ID:   "resident@gmail.com",
		Password: "resident",
		Role: model.RoleResident,
	}

	service := NewCredentialService(mockRepo)

	ctx := context.WithValue(context.Background(), utils.UserRoleKey, string(model.RoleResident))
	ctx = context.WithValue(ctx, utils.UserEmailKey, "dsfsl")
	ctx = context.WithValue(ctx, utils.UserFlatKey, "xxx")
	ctx = context.WithValue(ctx, utils.UserIDKey, "dfjld")

	err := service.DeleteResidentCredentials(ctx, "resident@gmail.com")
	if err == nil || err.Error() != "unauthorized: only admin can delete credentials" {
		t.Errorf("expected unauthorized error, got: %v", err)
	}
}