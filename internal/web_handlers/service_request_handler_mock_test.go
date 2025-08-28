package web_handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/mocks"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestGetAvailableTimeSlots_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	slots := []utils.TimeSlot{
		{
			Label:     "9:00 AM - 9:45 AM",
			StartTime: time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 1, 1, 9, 45, 0, 0, time.UTC),
		},
	}
	mockService.EXPECT().GetAvailableTimeSlots(model.Plumber).Return(slots)

	req := httptest.NewRequest(http.MethodGet, "/slots?serviceType=plumber", nil).WithContext(context.Background())
	w := httptest.NewRecorder()

	handler.GetAvailableTimeSlots(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp["message"] != "Time slot shown successfully" {
		t.Errorf("unexpected message: %v", resp["message"])
	}
	data := resp["data"].([]any)
	if len(data) != 1 {
		t.Errorf("expected 1 slot, got %d", len(data))
	}
}

func TestGetAvailableTimeSlots_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodPost, "/slots?serviceType=plumber", nil)
	w := httptest.NewRecorder()

	handler.GetAvailableTimeSlots(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestGetAvailableTimeSlots_InvalidServiceType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/slots?serviceType=carpenter", nil)
	w := httptest.NewRecorder()

	handler.GetAvailableTimeSlots(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetAvailableTimeSlots_EmptySlots(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	mockService.EXPECT().GetAvailableTimeSlots(model.Electrician).Return([]utils.TimeSlot{})

	req := httptest.NewRequest(http.MethodGet, "/slots?serviceType=electrician", nil)
	w := httptest.NewRecorder()

	handler.GetAvailableTimeSlots(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	data := resp["data"].([]any)
	if len(data) != 0 {
		t.Errorf("expected 0 slots, got %d", len(data))
	}
}

func TestGetAvailableTimeSlots_MultipleSlots(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	slots := []utils.TimeSlot{
		{
			Label:     "9:00 AM - 9:45 AM",
			StartTime: time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 1, 1, 9, 45, 0, 0, time.UTC),
		},
		{
			Label:     "9:45 AM - 10:30 AM",
			StartTime: time.Date(2025, 1, 1, 9, 45, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 1, 1, 10, 30, 0, 0, time.UTC),
		},
	}
	mockService.EXPECT().GetAvailableTimeSlots(model.Plumber).Return(slots)

	req := httptest.NewRequest(http.MethodGet, "/slots?serviceType=plumber", nil)
	w := httptest.NewRecorder()

	handler.GetAvailableTimeSlots(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	data := resp["data"].([]any)
	if len(data) != 2 {
		t.Errorf("expected 2 slots, got %d", len(data))
	}
}

func withUser(ctx context.Context, u *model.User) context.Context {
	ctx = context.WithValue(ctx, utils.UserIDKey, u.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, string(u.Role))
	ctx = context.WithValue(ctx, utils.UserEmailKey, u.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, u.Flat)
	return ctx
}

func TestBookServiceRequest_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "res1", Role: model.RoleResident, Email: "r", Flat: "101"}
	ctx := withUser(context.Background(), user)

	slot := utils.TimeSlot{
		Label:     "9:00 AM - 9:45 AM",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(45 * time.Minute),
	}

	body, _ := json.Marshal(map[string]any{"servicetype": "plumber", "slotid": 1})
	req := httptest.NewRequest(http.MethodPost, "/book", bytes.NewReader(body)).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().GetAvailableTimeSlots(model.ServiceType("plumber")).Return([]utils.TimeSlot{slot})
	mockService.EXPECT().BookServiceRequest(gomock.Any()).Return(nil)

	handler.BookServiceRequest(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestBookServiceRequest_InvalidMethod(t *testing.T) {
	handler := &ServiceRequestHandler{}
	req := httptest.NewRequest(http.MethodGet, "/book", nil)
	w := httptest.NewRecorder()

	handler.BookServiceRequest(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestBookServiceRequest_NoUserInContext(t *testing.T) {
	handler := &ServiceRequestHandler{}
	req := httptest.NewRequest(http.MethodPost, "/book", nil)
	w := httptest.NewRecorder()

	handler.BookServiceRequest(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestBookServiceRequest_WrongRole(t *testing.T) {
	user := &model.User{ID: "off1", Role: model.RoleOfficer, Email: "o", Flat: "NA"}
	ctx := withUser(context.Background(), user)

	handler := &ServiceRequestHandler{}
	req := httptest.NewRequest(http.MethodPost, "/book", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.BookServiceRequest(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestBookServiceRequest_InvalidJSON(t *testing.T) {
	user := &model.User{ID: "res1", Role: model.RoleResident, Email: "r", Flat: "101"}
	ctx := withUser(context.Background(), user)

	handler := &ServiceRequestHandler{}
	req := httptest.NewRequest(http.MethodPost, "/book", bytes.NewReader([]byte("{invalid json}"))).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.BookServiceRequest(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestBookServiceRequest_SlotIDTooLow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "res1", Role: model.RoleResident, Email: "r", Flat: "101"}
	ctx := withUser(context.Background(), user)

	slot := utils.TimeSlot{Label: "9:00 AM - 9:45 AM", StartTime: time.Now(), EndTime: time.Now().Add(45 * time.Minute)}

	body, _ := json.Marshal(map[string]any{"servicetype": "plumber", "slotid": 0})
	req := httptest.NewRequest(http.MethodPost, "/book", bytes.NewReader(body)).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().GetAvailableTimeSlots(model.ServiceType("plumber")).Return([]utils.TimeSlot{slot})

	handler.BookServiceRequest(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestBookServiceRequest_SlotIDTooHigh(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "res1", Role: model.RoleResident, Email: "r", Flat: "101"}
	ctx := withUser(context.Background(), user)

	slot := utils.TimeSlot{Label: "9:00 AM - 9:45 AM", StartTime: time.Now(), EndTime: time.Now().Add(45 * time.Minute)}

	body, _ := json.Marshal(map[string]any{"servicetype": "plumber", "slotid": 5})
	req := httptest.NewRequest(http.MethodPost, "/book", bytes.NewReader(body)).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().GetAvailableTimeSlots(model.ServiceType("plumber")).Return([]utils.TimeSlot{slot})

	handler.BookServiceRequest(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestBookServiceRequest_BookServiceFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "res1", Role: model.RoleResident, Email: "r", Flat: "101"}
	ctx := withUser(context.Background(), user)

	slot := utils.TimeSlot{Label: "9:00 AM - 9:45 AM", StartTime: time.Now(), EndTime: time.Now().Add(45 * time.Minute)}

	body, _ := json.Marshal(map[string]any{"servicetype": "plumber", "slotid": 1})
	req := httptest.NewRequest(http.MethodPost, "/book", bytes.NewReader(body)).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().GetAvailableTimeSlots(model.ServiceType("plumber")).Return([]utils.TimeSlot{slot})
	mockService.EXPECT().BookServiceRequest(gomock.Any()).Return(errors.New("failed"))

	handler.BookServiceRequest(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestReschedule_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := &ServiceRequestHandler{}

	req := httptest.NewRequest(http.MethodGet, "/service/reschedule/123", nil)
	w := httptest.NewRecorder()

	handler.RescheduleServiceRequest(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestReschedule_UserNotInContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := &ServiceRequestHandler{}

	req := httptest.NewRequest(http.MethodPatch, "/service/reschedule/123", nil)
	w := httptest.NewRecorder()

	handler.RescheduleServiceRequest(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestReschedule_InvalidPathFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := &ServiceRequestHandler{}

	user := &model.User{ID: "u1", Role: model.RoleResident}
	ctx := withUser(context.Background(), user)

	req := httptest.NewRequest(http.MethodPatch, "/service/foo/123", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.RescheduleServiceRequest(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestReschedule_InvalidRequestID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := &ServiceRequestHandler{}

	user := &model.User{ID: "u1", Role: model.RoleResident}
	ctx := withUser(context.Background(), user)

	req := httptest.NewRequest(http.MethodPatch, "/service/reschedule/abc123", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.RescheduleServiceRequest(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestReschedule_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := &ServiceRequestHandler{}

	user := &model.User{ID: "u1", Role: model.RoleResident}
	ctx := withUser(context.Background(), user)

	body := strings.NewReader("{slot:}") // malformed JSON
	req := httptest.NewRequest(http.MethodPatch, "/service/reschedule/"+uuid.New().String(), body).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.RescheduleServiceRequest(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestReschedule_RequestNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "u1", Role: model.RoleResident}
	ctx := withUser(context.Background(), user)

	reqID := uuid.New()
	body, _ := json.Marshal(map[string]int{"slot": 1})
	req := httptest.NewRequest(http.MethodPatch, "/service/reschedule/"+reqID.String(), bytes.NewReader(body)).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().GetServiceTypeByID(reqID).Return(model.ServiceType(""), errors.New("error"))

	handler.RescheduleServiceRequest(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("got %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestReschedule_InvalidSlot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "u1", Role: model.RoleResident}
	ctx := withUser(context.Background(), user)

	reqID := uuid.New()
	body, _ := json.Marshal(map[string]int{"slot": 5})
	req := httptest.NewRequest(http.MethodPatch, "/service/reschedule/"+reqID.String(), bytes.NewReader(body)).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().GetServiceTypeByID(reqID).Return(model.Plumber, nil)
	mockService.EXPECT().GetAvailableTimeSlots(model.Plumber).Return([]utils.TimeSlot{{Label: "9-10"}})

	handler.RescheduleServiceRequest(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestReschedule_ServiceFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "u1", Role: model.RoleResident}
	ctx := withUser(context.Background(), user)

	reqID := uuid.New()
	body, _ := json.Marshal(map[string]int{"slot": 1})
	req := httptest.NewRequest(http.MethodPatch, "/service/reschedule/"+reqID.String(), bytes.NewReader(body)).WithContext(ctx)
	w := httptest.NewRecorder()

	slot := utils.TimeSlot{Label: "9-10", StartTime: time.Now(), EndTime: time.Now().Add(45 * time.Minute)}

	mockService.EXPECT().GetServiceTypeByID(reqID).Return(model.Plumber, nil)
	mockService.EXPECT().GetAvailableTimeSlots(model.Plumber).Return([]utils.TimeSlot{slot})
	mockService.EXPECT().RescheduleServiceRequest(user.ID, reqID, slot, model.Plumber).Return(errors.New("error"))

	handler.RescheduleServiceRequest(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestReschedule_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "u1", Role: model.RoleResident}
	ctx := withUser(context.Background(), user)

	reqID := uuid.New()
	body, _ := json.Marshal(map[string]int{"slot": 1})
	req := httptest.NewRequest(http.MethodPatch, "/service/reschedule/"+reqID.String(), bytes.NewReader(body)).WithContext(ctx)
	w := httptest.NewRecorder()

	slot := utils.TimeSlot{Label: "9-10", StartTime: time.Now(), EndTime: time.Now().Add(45 * time.Minute)}

	mockService.EXPECT().GetServiceTypeByID(reqID).Return(model.Plumber, nil)
	mockService.EXPECT().GetAvailableTimeSlots(model.Plumber).Return([]utils.TimeSlot{slot})
	mockService.EXPECT().RescheduleServiceRequest(user.ID, reqID, slot, model.Plumber).Return(nil)

	handler.RescheduleServiceRequest(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCancelServiceRequest_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/service/cancel/123", nil) 
	w := httptest.NewRecorder()

	handler.CancelServiceRequest(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestCancelServiceRequest_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodDelete, "/service/cancel/123", nil)
	w := httptest.NewRecorder()

	handler.CancelServiceRequest(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestCancelServiceRequest_MissingRequestID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "u1", Role: model.RoleResident}
	ctx := withUser(context.Background(), user)

	req := httptest.NewRequest(http.MethodDelete, "/service/cancel", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.CancelServiceRequest(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCancelServiceRequest_InvalidRequestID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "u1", Role: model.RoleResident}
	ctx := withUser(context.Background(), user)

	req := httptest.NewRequest(http.MethodDelete, "/service/cancel/not-a-uuid", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.CancelServiceRequest(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCancelServiceRequest_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "u1", Role: model.RoleResident}
	ctx := withUser(context.Background(), user)

	reqID := uuid.New()
	req := httptest.NewRequest(http.MethodDelete, "/service/cancel/"+reqID.String(), nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().CancelServiceRequest(user.ID, reqID).Return(errors.New("db error"))

	handler.CancelServiceRequest(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestCancelServiceRequest_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "u1", Role: model.RoleResident}
	ctx := withUser(context.Background(), user)

	reqID := uuid.New()
	req := httptest.NewRequest(http.MethodDelete, "/service/cancel/"+reqID.String(), nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().CancelServiceRequest(user.ID, reqID).Return(nil)

	handler.CancelServiceRequest(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestApproveRequest_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := &ServiceRequestHandler{Service: nil}

	req := httptest.NewRequest(http.MethodGet, "/service/approve/123", nil)
	w := httptest.NewRecorder()

	handler.ApproveRequest(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestApproveRequest_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := &ServiceRequestHandler{Service: nil}

	req := httptest.NewRequest(http.MethodPatch, "/service/approve/123", nil)
	w := httptest.NewRecorder()

	handler.ApproveRequest(w, req) 

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestApproveRequest_UnauthorizedRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := &ServiceRequestHandler{Service: nil}

	user := &model.User{ID: "u1", Role: "resident"}
	ctx := withUser(context.Background(), user)

	req := httptest.NewRequest(http.MethodPatch, "/service/approve/123", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ApproveRequest(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestApproveRequest_InvalidPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := &ServiceRequestHandler{Service: nil}

	user := &model.User{ID: "u1", Role: "officer"}
	ctx := withUser(context.Background(), user)

	req := httptest.NewRequest(http.MethodPatch, "/service/invalid/123", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ApproveRequest(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestApproveRequest_InvalidRequestID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := &ServiceRequestHandler{Service: nil}

	user := &model.User{ID: "u1", Role: "admin"}
	ctx := withUser(context.Background(), user)

	req := httptest.NewRequest(http.MethodPatch, "/service/approve/not-a-uuid", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ApproveRequest(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestApproveRequest_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "u1", Role: "officer"}
	ctx := withUser(context.Background(), user)

	reqID := uuid.New()
	req := httptest.NewRequest(http.MethodPatch, "/service/approve/"+reqID.String(), nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().ApproveServiceRequest(reqID).Return(errors.New("db error"))

	handler.ApproveRequest(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestApproveRequest_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "u1", Role: "admin"}
	ctx := withUser(context.Background(), user)

	reqID := uuid.New()
	req := httptest.NewRequest(http.MethodPatch, "/service/approve/"+reqID.String(), nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().ApproveServiceRequest(reqID).Return(nil)

	handler.ApproveRequest(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetRequestsOfResident_InvalidMethod(t *testing.T) {
	handler := &ServiceRequestHandler{Service: nil}

	req := httptest.NewRequest(http.MethodPost, "/service/requests?status=pending", nil)
	w := httptest.NewRecorder()

	handler.GetRequestsOfResident(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestGetRequestsOfResident_UserNotFound(t *testing.T) {
	handler := &ServiceRequestHandler{Service: nil}

	req := httptest.NewRequest(http.MethodGet, "/service/requests?status=pending", nil)
	w := httptest.NewRecorder()

	handler.GetRequestsOfResident(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestGetRequestsOfResident_MissingStatus(t *testing.T) {
	handler := &ServiceRequestHandler{Service: nil}

	user := &model.User{ID: "u1", Role: model.RoleResident}
	ctx := withUser(context.Background(), user)

	req := httptest.NewRequest(http.MethodGet, "/service/requests", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetRequestsOfResident(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetRequestsOfResident_ResidentSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "u1", Role: model.RoleResident}
	ctx := withUser(context.Background(), user)

	reqs := []model.ServiceRequest{{RequestID: uuid.New(), Status: model.StatusPending}}
	mockService.EXPECT().GetServiceRequestsByStatus(user.ID, model.Status("pending")).Return(reqs)

	req := httptest.NewRequest(http.MethodGet, "/service/requests?status=pending", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetRequestsOfResident(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetRequestsOfResident_AdminMissingID(t *testing.T) {
	handler := &ServiceRequestHandler{Service: nil}

	user := &model.User{ID: "admin1", Role: "admin"}
	ctx := withUser(context.Background(), user)

	req := httptest.NewRequest(http.MethodGet, "/service/requests?status=pending", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetRequestsOfResident(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetRequestsOfResident_AdminSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: "admin"}
	ctx := withUser(context.Background(), user)

	targetResidentID := "res1"
	reqs := []model.ServiceRequest{{RequestID: uuid.New(), Status: model.StatusApproved}}
	mockService.EXPECT().GetServiceRequestsByStatus(targetResidentID, model.Status("approved")).Return(reqs)

	req := httptest.NewRequest(http.MethodGet, "/service/requests?status=approved&id="+targetResidentID, nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetRequestsOfResident(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetRequestsByServiceTypeAndStatus_InvalidMethod(t *testing.T) {
	handler := &ServiceRequestHandler{Service: nil}

	req := httptest.NewRequest(http.MethodPost, "/service/requests?status=pending&serviceType=plumber", nil)
	w := httptest.NewRecorder()

	handler.GetRequestsByServiceTypeAndStatus(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestGetRequestsByServiceTypeAndStatus_UserNotFound(t *testing.T) {
	handler := &ServiceRequestHandler{Service: nil}

	req := httptest.NewRequest(http.MethodGet, "/service/requests?status=pending&serviceType=plumber", nil)
	w := httptest.NewRecorder()

	handler.GetRequestsByServiceTypeAndStatus(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestGetRequestsByServiceTypeAndStatus_ForbiddenForResident(t *testing.T) {
	handler := &ServiceRequestHandler{Service: nil}

	user := &model.User{ID: "u1", Role: model.RoleResident}
	ctx := withUser(context.Background(), user)

	req := httptest.NewRequest(http.MethodGet, "/service/requests?status=pending&serviceType=plumber", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetRequestsByServiceTypeAndStatus(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestGetRequestsByServiceTypeAndStatus_MissingParams(t *testing.T) {
	handler := &ServiceRequestHandler{Service: nil}

	user := &model.User{ID: "admin1", Role: "admin"}
	ctx := withUser(context.Background(), user)

	req := httptest.NewRequest(http.MethodGet, "/service/requests?status=pending", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetRequestsByServiceTypeAndStatus(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetRequestsByServiceTypeAndStatus_PlumberPending(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: "admin"}
	ctx := withUser(context.Background(), user)

	expected := []model.ServiceRequest{{RequestID: uuid.New(), ServiceType: model.Plumber, Status: model.StatusPending}}
	mockService.EXPECT().GetPendingRequestsByServiceType(model.Plumber).Return(expected)

	req := httptest.NewRequest(http.MethodGet, "/service/requests?status=pending&serviceType=plumber", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetRequestsByServiceTypeAndStatus(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetRequestsByServiceTypeAndStatus_PlumberApproved(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: "admin"}
	ctx := withUser(context.Background(), user)

	expected := []model.ServiceRequest{{RequestID: uuid.New(), ServiceType: model.Plumber, Status: model.StatusApproved}}
	mockService.EXPECT().GetApprovedRequestsByServiceType(model.Plumber).Return(expected)

	req := httptest.NewRequest(http.MethodGet, "/service/requests?status=approved&serviceType=plumber", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetRequestsByServiceTypeAndStatus(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetRequestsByServiceTypeAndStatus_ElectricianPending(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: "admin"}
	ctx := withUser(context.Background(), user)

	expected := []model.ServiceRequest{{RequestID: uuid.New(), ServiceType: model.Electrician, Status: model.StatusPending}}
	mockService.EXPECT().GetPendingRequestsByServiceType(model.Electrician).Return(expected)

	req := httptest.NewRequest(http.MethodGet, "/service/requests?status=pending&serviceType=electrician", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetRequestsByServiceTypeAndStatus(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetRequestsByServiceTypeAndStatus_ElectricianApproved(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockServiceRequestServiceInterface(ctrl)
	handler := &ServiceRequestHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: "admin"}
	ctx := withUser(context.Background(), user)

	expected := []model.ServiceRequest{{RequestID: uuid.New(), ServiceType: model.Electrician, Status: model.StatusApproved}}
	mockService.EXPECT().GetApprovedRequestsByServiceType(model.Electrician).Return(expected)

	req := httptest.NewRequest(http.MethodGet, "/service/requests?status=approved&serviceType=electrician", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetRequestsByServiceTypeAndStatus(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}