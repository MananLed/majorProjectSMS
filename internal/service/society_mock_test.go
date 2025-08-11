package service

import (
	"context"
	"testing"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
)

type MockSocietyRepo struct {
	Users []model.User
}

func (m *MockSocietyRepo) GetAllResidents() ([]model.User, error) {
	var residents []model.User
	for _, u := range m.Users {
		if u.Role == model.RoleResident {
			residents = append(residents, u)
		}
	}
	return residents, nil
}

func (m *MockSocietyRepo) GetAllOfficers() ([]model.User, error) {
	var officers []model.User
	for _, u := range m.Users {
		if u.Role == model.RoleOfficer {
			officers = append(officers, u)
		}
	}
	return officers, nil
}

func ctxWithUser(user *model.User) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, user.Role)
	return ctx
}

func TestGetAllResidentsAsAdmin(t *testing.T) {
	mockRepo := &MockSocietyRepo{
		Users: []model.User{
			{
				ID:    "r1",
				Role:  model.RoleResident,
				Email: "resident@example.com",
			},
			{
				ID:    "o1",
				Role:  model.RoleOfficer,
				Email: "officer@example.com",
			},
		},
	}

	service := NewSocietyService(mockRepo)

	admin := &model.User{ID: "a1", Role: model.RoleAdmin}
	ctx := ctxWithUser(admin)

	residents, err := service.GetAllResidents(ctx)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(residents) != 1 || residents[0].ID != "r1" {
		t.Errorf("unexpected residents list: %v", residents)
	}
}

func TestGetAllOfficersAsAdmin(t *testing.T) {
	mockRepo := &MockSocietyRepo{
		Users: []model.User{
			{ID: "r1", Role: model.RoleResident},
			{ID: "o1", Role: model.RoleOfficer},
		},
	}

	service := NewSocietyService(mockRepo)

	admin := &model.User{ID: "a1",Password: "", Role: model.RoleAdmin}
	ctx := ctxWithUser(admin)

	officers, err := service.GetAllOfficers(ctx)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(officers) != 1 || officers[0].ID != "o1" {
		t.Errorf("unexpected officers list: %v", officers)
	}
}

func TestGetAllResidentsAsNonAdmin(t *testing.T) {
	mockRepo := &MockSocietyRepo{}
	service := NewSocietyService(mockRepo)

	officer := &model.User{ID: "o1",Password: "", Role: model.RoleOfficer}
	ctx := ctxWithUser(officer)

	_, err := service.GetAllResidents(ctx)
	if err == nil {
		t.Error("expected error for non-admin user")
	}
}

func TestGetAllOfficersAsNonAdmin(t *testing.T) {
	mockRepo := &MockSocietyRepo{}
	service := NewSocietyService(mockRepo)

	resident := &model.User{ID: "r1",Password: "", Role: model.RoleResident}
	ctx := ctxWithUser(resident)

	_, err := service.GetAllOfficers(ctx)
	if err == nil {
		t.Error("expected error for non-admin user")
	}
}

func TestMissingUserInContext(t *testing.T) {
	mockRepo := &MockSocietyRepo{}
	service := NewSocietyService(mockRepo)

	ctx := context.Background()

	_, err := service.GetAllResidents(ctx)
	if err == nil {
		t.Error("expected error due to missing user in context")
	}
}
