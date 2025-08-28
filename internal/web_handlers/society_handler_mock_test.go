package web_handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MananLed/majorProjectSMS/internal/mocks"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"go.uber.org/mock/gomock"
)

func TestViewResidents_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockSocietyServiceInterface(ctrl)
	handler := &SocietyHandler{SocietyService: mockService}

	req := httptest.NewRequest(http.MethodPost, "/residents", nil)
	w := httptest.NewRecorder()

	handler.ViewResidents(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestViewResidents_UnauthorizedUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockSocietyServiceInterface(ctrl)
	handler := &SocietyHandler{SocietyService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/residents", nil)
	w := httptest.NewRecorder()

	handler.ViewResidents(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestViewResidents_ServiceFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockSocietyServiceInterface(ctrl)
	handler := &SocietyHandler{SocietyService: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	req := httptest.NewRequest(http.MethodGet, "/residents", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().
		GetAllResidents(ctx).
		Return(nil, errors.New("db error"))

	handler.ViewResidents(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestViewResidents_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockSocietyServiceInterface(ctrl)
	handler := &SocietyHandler{SocietyService: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	req := httptest.NewRequest(http.MethodGet, "/residents", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	expectedResidents := []model.User{
		{ID: "r1", FirstName: "John"},
		{ID: "r2", FirstName: "Jane"},
	}

	mockService.EXPECT().
		GetAllResidents(ctx).
		Return(expectedResidents, nil)

	handler.ViewResidents(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestViewOfficers_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockSocietyServiceInterface(ctrl)
	handler := &SocietyHandler{SocietyService: mockService}

	req := httptest.NewRequest(http.MethodPost, "/officers", nil)
	w := httptest.NewRecorder()

	handler.ViewOfficers(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestViewOfficers_UnauthorizedUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockSocietyServiceInterface(ctrl)
	handler := &SocietyHandler{SocietyService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/officers", nil)
	w := httptest.NewRecorder()

	handler.ViewOfficers(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestViewOfficers_ServiceFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockSocietyServiceInterface(ctrl)
	handler := &SocietyHandler{SocietyService: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	req := httptest.NewRequest(http.MethodGet, "/officers", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().
		GetAllOfficers(ctx).
		Return(nil, errors.New("db error"))

	handler.ViewOfficers(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestViewOfficers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockSocietyServiceInterface(ctrl)
	handler := &SocietyHandler{SocietyService: mockService}


	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	req := httptest.NewRequest(http.MethodGet, "/officers", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	expectedOfficers := []model.User{
		{ID: "o1", FirstName: "Alice"},
		{ID: "o2", FirstName: "Bob"},
	}

	mockService.EXPECT().
		GetAllOfficers(ctx).
		Return(expectedOfficers, nil)

	handler.ViewOfficers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}