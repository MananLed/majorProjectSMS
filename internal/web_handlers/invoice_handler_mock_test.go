package web_handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MananLed/majorProjectSMS/internal/mocks"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"go.uber.org/mock/gomock"
)

func TestIssueInvoice_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockInvoiceServiceInterface(ctrl)
	handler := &InvoiceHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/invoice", nil)
	w := httptest.NewRecorder()

	handler.IssueInvoice(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestIssueInvoice_NoUserInContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockInvoiceServiceInterface(ctrl)
	handler := &InvoiceHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodPost, "/invoice", nil)
	w := httptest.NewRecorder()

	handler.IssueInvoice(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestIssueInvoice_UnauthorizedRoleResident(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockInvoiceServiceInterface(ctrl)
	handler := &InvoiceHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "resident")

	req := httptest.NewRequest(http.MethodPost, "/invoice", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.IssueInvoice(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestIssueInvoice_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockInvoiceServiceInterface(ctrl)
	handler := &InvoiceHandler{Service: mockService}

	body := strings.NewReader(`{amount: }`) 
	req := httptest.NewRequest(http.MethodPost, "/invoice", body)

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}
	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	req = httptest.NewRequest(http.MethodPost, "/invoice", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.IssueInvoice(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestIssueInvoice_InvalidAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockInvoiceServiceInterface(ctrl)
	handler := &InvoiceHandler{Service: mockService}

	body := strings.NewReader(`{"amount": 0}`)
	req := httptest.NewRequest(http.MethodPost, "/invoice", body)

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}
	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	req = httptest.NewRequest(http.MethodPost, "/invoice", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.IssueInvoice(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestIssueInvoice_ServiceError(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockService := mocks.NewMockInvoiceServiceInterface(ctrl)
    handler := &InvoiceHandler{Service: mockService}

    body := strings.NewReader(`{"amount": 5000}`)
    req := httptest.NewRequest(http.MethodPost, "/invoice", body)


    ctx := context.WithValue(req.Context(), utils.UserIDKey, "admin1")
    ctx = context.WithValue(ctx, utils.UserRoleKey, string(model.RoleAdmin)) 
    ctx = context.WithValue(ctx, utils.UserEmailKey, "dsfs")
    ctx = context.WithValue(ctx, utils.UserFlatKey, "xxx")

    req = req.WithContext(ctx)
    w := httptest.NewRecorder()

    mockService.EXPECT().
        GenerateInvoice(5000.0, gomock.Any(), gomock.Any()).
        Return(errors.New("failed to generate"))

    handler.IssueInvoice(w, req)

    if w.Code != http.StatusInternalServerError {
        t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
    }
}

func TestIssueInvoice_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockService := mocks.NewMockInvoiceServiceInterface(ctrl)
    handler := &InvoiceHandler{Service: mockService}

    body := strings.NewReader(`{"amount": 5000}`)
    req := httptest.NewRequest(http.MethodPost, "/invoice", body)

    ctx := context.WithValue(req.Context(), utils.UserIDKey, "admin1")
    ctx = context.WithValue(ctx, utils.UserRoleKey, string(model.RoleAdmin))
    ctx = context.WithValue(ctx, utils.UserEmailKey, "dlfj")
    ctx = context.WithValue(ctx, utils.UserFlatKey, "xxx")

    req = req.WithContext(ctx) 
    w := httptest.NewRecorder()

    mockService.EXPECT().
        GenerateInvoice(5000.0, gomock.Any(), gomock.Any()).
        Return(nil)

    handler.IssueInvoice(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("got %d, want %d", w.Code, http.StatusOK)
    }

    var resp map[string]any
    if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
        t.Errorf("failed to parse response JSON: %v", err)
    }
    if resp["message"] != "Invoice issued successfully" {
        t.Errorf("unexpected message: %v", resp["message"])
    }
}
