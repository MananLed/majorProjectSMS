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
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestGetFeedbacks_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodPost, "/feedbacks", nil)
	w := httptest.NewRecorder()

	handler.GetFeedbacks(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestGetFeedbacks_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/feedbacks", nil)
	w := httptest.NewRecorder()

	handler.GetFeedbacks(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestGetFeedbacks_ResidentSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	user := &model.User{ID: "res1", Role: model.RoleResident, Email: "x", Flat: "101"}
	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, string(user.Role))
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)

	req := httptest.NewRequest(http.MethodGet, "/feedbacks", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	expected := []model.Feedback{{ID: uuid.New(), Content: "Good"}}
	mockService.EXPECT().GetFeedbackByID("res1").Return(expected, nil)

	handler.GetFeedbacks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetFeedbacks_ResidentError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	user := &model.User{ID: "res1", Role: model.RoleResident, Email: "x", Flat: "101"}
	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, string(user.Role))
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)

	req := httptest.NewRequest(http.MethodGet, "/feedbacks", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().GetFeedbackByID("res1").Return(nil, errors.New("db error"))

	handler.GetFeedbacks(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestGetFeedbacks_AdminAllSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "a", Flat: "NA"}
	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, string(user.Role))
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)

	req := httptest.NewRequest(http.MethodGet, "/feedbacks", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	expected := []model.Feedback{{ID: uuid.New(), Content: "Great"}}
	mockService.EXPECT().GetFeedbacks().Return(expected, nil)

	handler.GetFeedbacks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetFeedbacks_AdminByResidentIDSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "a", Flat: "NA"}
	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, string(user.Role))
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)

	req := httptest.NewRequest(http.MethodGet, "/feedbacks?residentId=res1", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	expected := []model.Feedback{{ID: uuid.New(), Content: "Ok"}}
	mockService.EXPECT().GetFeedbackByID("res1").Return(expected, nil)

	handler.GetFeedbacks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetFeedbacks_AdminAllError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "a", Flat: "NA"}
	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, string(user.Role))
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)

	req := httptest.NewRequest(http.MethodGet, "/feedbacks", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().GetFeedbacks().Return(nil, errors.New("db error"))

	handler.GetFeedbacks(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestGetFeedbacks_OfficerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	user := &model.User{ID: "off1", Role: "officer", Email: "o", Flat: "NA"}
	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, string(user.Role))
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)

	req := httptest.NewRequest(http.MethodGet, "/feedbacks", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	expected := []model.Feedback{{ID: uuid.New(), Content: "Ok"}}
	mockService.EXPECT().GetFeedbacks().Return(expected, nil)

	handler.GetFeedbacks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetFeedbacks_UnauthorizedRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	user := &model.User{ID: "guest1", Role: "guest", Email: "g", Flat: "NA"}
	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, string(user.Role))
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)

	req := httptest.NewRequest(http.MethodGet, "/feedbacks", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetFeedbacks(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestGiveFeedback_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/feedback", nil)
	w := httptest.NewRecorder()

	handler.GiveFeedback(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestGiveFeedback_UserNotInContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodPost, "/feedback", nil)
	w := httptest.NewRecorder()

	handler.GiveFeedback(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestGiveFeedback_UnauthorizedRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: "admin", Email: "a", Flat: "NA"}
	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, string(user.Role))
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)

	req := httptest.NewRequest(http.MethodPost, "/feedback", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GiveFeedback(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestGiveFeedback_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	user := &model.User{ID: "res1", Role: "resident", Email: "r", Flat: "101"}
	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, string(user.Role))
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)

	req := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(`{invalid}`)).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GiveFeedback(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGiveFeedback_InvalidRatingLow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	user := &model.User{ID: "res1", Role: "resident", Email: "r", Flat: "101"}
	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, string(user.Role))
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)

	body := strings.NewReader(`{"rating":0,"content":"bad"}`)
	req := httptest.NewRequest(http.MethodPost, "/feedback", body).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GiveFeedback(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGiveFeedback_InvalidRatingHigh(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	user := &model.User{ID: "res1", Role: "resident", Email: "r", Flat: "101"}
	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, string(user.Role))
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)

	body := strings.NewReader(`{"rating":6,"content":"too good"}`)
	req := httptest.NewRequest(http.MethodPost, "/feedback", body).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GiveFeedback(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGiveFeedback_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	user := &model.User{ID: "res1", Role: "resident", Email: "r", Flat: "101"}
	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, string(user.Role))
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)

	body := strings.NewReader(`{"rating":4,"content":"ok"}`)
	req := httptest.NewRequest(http.MethodPost, "/feedback", body).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().IssueFeedback("ok", user.ID, user.Flat, int32(4)).Return(errors.New("failed"))

	handler.GiveFeedback(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestGiveFeedback_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockFeedbackServiceInterface(ctrl)
	handler := &FeedbackHandler{Service: mockService}

	user := &model.User{ID: "res1", Role: "resident", Email: "r", Flat: "101"}
	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, string(user.Role))
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)

	body := strings.NewReader(`{"rating":5,"content":"great service"}`)
	req := httptest.NewRequest(http.MethodPost, "/feedback", body).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().IssueFeedback("great service", user.ID, user.Flat, int32(5)).Return(nil)

	handler.GiveFeedback(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("failed to parse response JSON: %v", err)
	}
	if resp["message"] != "Feedback issued successfully" {
		t.Errorf("unexpected message: %v", resp["message"])
	}
}
