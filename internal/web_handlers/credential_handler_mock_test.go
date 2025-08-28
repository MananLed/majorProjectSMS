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

func TestDeleteOfficer_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockCredentialServiceInterface(ctrl)
	handler := &CredentialHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/delete-officer", nil)
	w := httptest.NewRecorder()

	handler.DeleteOfficer(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestDeleteOfficer_UnauthorizedUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockCredentialServiceInterface(ctrl)
	handler := &CredentialHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodDelete, "/delete-officer?id=123", nil)
	w := httptest.NewRecorder()

	handler.DeleteOfficer(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestDeleteOfficer_MissingID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockCredentialServiceInterface(ctrl)
	handler := &CredentialHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	req := httptest.NewRequest(http.MethodDelete, "/delete-officer", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.DeleteOfficer(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestDeleteOfficer_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockCredentialServiceInterface(ctrl)
	handler := &CredentialHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	req := httptest.NewRequest(http.MethodDelete, "/delete-officer?id=123", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().
		DeleteOfficerCredentials(ctx, "123").
		Return(errors.New("db error"))

	handler.DeleteOfficer(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestDeleteOfficer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockCredentialServiceInterface(ctrl)
	handler := &CredentialHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	req := httptest.NewRequest(http.MethodDelete, "/delete-officer?id=123", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().
		DeleteOfficerCredentials(ctx, "123").
		Return(nil)

	handler.DeleteOfficer(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestDeleteResident_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockCredentialServiceInterface(ctrl)
	handler := &CredentialHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/delete-resident", nil) 
	w := httptest.NewRecorder()

	handler.DeleteResident(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestDeleteResident_UnauthorizedUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockCredentialServiceInterface(ctrl)
	handler := &CredentialHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodDelete, "/delete-resident?id=123", nil)
	w := httptest.NewRecorder()

	handler.DeleteResident(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want %d", w.Code, http.StatusForbidden)
	}
}


func TestDeleteResident_MissingID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockCredentialServiceInterface(ctrl)
	handler := &CredentialHandler{Service: mockService}


	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	req := httptest.NewRequest(http.MethodDelete, "/delete-resident", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.DeleteResident(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}


func TestDeleteResident_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockCredentialServiceInterface(ctrl)
	handler := &CredentialHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	req := httptest.NewRequest(http.MethodDelete, "/delete-resident?id=123", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().
		DeleteResidentCredentials(ctx, "123").
		Return(errors.New("db error"))

	handler.DeleteResident(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestDeleteResident_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockCredentialServiceInterface(ctrl)
	handler := &CredentialHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	req := httptest.NewRequest(http.MethodDelete, "/delete-resident?id=123", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().
		DeleteResidentCredentials(ctx, "123").
		Return(nil)

	handler.DeleteResident(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}