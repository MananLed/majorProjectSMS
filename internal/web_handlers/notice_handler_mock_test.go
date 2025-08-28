package web_handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/mocks"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestIssueNotice_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockNoticeServiceInterface(ctrl)
	handler := &NoticeHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/issue-notice", nil)
	w := httptest.NewRecorder()

	handler.IssueNotice(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestIssueNotice_UnauthorizedUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockNoticeServiceInterface(ctrl)
	handler := &NoticeHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "resident")

	body := bytes.NewBufferString(`{"content":"Test Notice"}`)
	req := httptest.NewRequest(http.MethodPost, "/issue-notice", body).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.IssueNotice(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestIssueNotice_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockNoticeServiceInterface(ctrl)
	handler := &NoticeHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	body := bytes.NewBufferString(`{"bad_json"}`)
	req := httptest.NewRequest(http.MethodPost, "/issue-notice", body).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.IssueNotice(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}

	body2 := bytes.NewBufferString(`{"content":""}`)
	req2 := httptest.NewRequest(http.MethodPost, "/issue-notice", body2).WithContext(ctx)
	w2 := httptest.NewRecorder()

	handler.IssueNotice(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w2.Code, http.StatusBadRequest)
	}
}

func TestIssueNotice_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockNoticeServiceInterface(ctrl)
	handler := &NoticeHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "officer")

	body := bytes.NewBufferString(`{"content":"Important Notice"}`)
	req := httptest.NewRequest(http.MethodPost, "/issue-notice", body).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().
		IssueNotice("Important Notice", gomock.Any(), gomock.Any()).
		Return(errors.New("db error"))

	handler.IssueNotice(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestIssueNotice_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockNoticeServiceInterface(ctrl)
	handler := &NoticeHandler{Service: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	body := bytes.NewBufferString(`{"content":"Maintenance Work Scheduled"}`)
	req := httptest.NewRequest(http.MethodPost, "/issue-notice", body).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().
		IssueNotice("Maintenance Work Scheduled", gomock.Any(), gomock.Any()).
		Return(nil)

	handler.IssueNotice(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetNotices_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockNoticeServiceInterface(ctrl)
	handler := &NoticeHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/notices", nil)
	w := httptest.NewRecorder()

	mockService.EXPECT().
		GetNotices().
		Return(nil, errors.New("db error"))

	handler.GetNotices(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestGetNotices_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockNoticeServiceInterface(ctrl)
	handler := &NoticeHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/notices", nil)
	w := httptest.NewRecorder()

	expectedNotices := []model.Notice{
		{ID: uuid.New(), Content: "Water Supply Disruption"},
		{ID: uuid.New(), Content: "Fire Drill Scheduled"},
	}

	mockService.EXPECT().
		GetNotices().
		Return(expectedNotices, nil)

	handler.GetNotices(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetNoticesByMonthYear_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockNoticeServiceInterface(ctrl)
	handler := &NoticeHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodPost, "/notices?year=2024", nil)
	w := httptest.NewRecorder()

	handler.GetNoticesByMonthYear(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestGetNoticesByMonthYear_MissingYearAndMonth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := &NoticeHandler{Service: mocks.NewMockNoticeServiceInterface(ctrl)}

	req := httptest.NewRequest(http.MethodGet, "/notices", nil)
	w := httptest.NewRecorder()

	handler.GetNoticesByMonthYear(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetNoticesByMonthYear_MonthWithoutYear(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := &NoticeHandler{Service: mocks.NewMockNoticeServiceInterface(ctrl)}

	req := httptest.NewRequest(http.MethodGet, "/notices?month=5", nil)
	w := httptest.NewRecorder()

	handler.GetNoticesByMonthYear(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetNoticesByMonthYear_InvalidYearFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := &NoticeHandler{Service: mocks.NewMockNoticeServiceInterface(ctrl)}

	req := httptest.NewRequest(http.MethodGet, "/notices?year=abcd", nil)
	w := httptest.NewRecorder()

	handler.GetNoticesByMonthYear(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetNoticesByMonthYear_YearServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockNoticeServiceInterface(ctrl)
	handler := &NoticeHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/notices?year=2024", nil)
	w := httptest.NewRecorder()

	mockService.EXPECT().
		GetNoticesByYear(2024).
		Return(nil, errors.New("not found"))

	handler.GetNoticesByMonthYear(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("got %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestGetNoticesByMonthYear_YearSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockNoticeServiceInterface(ctrl)
	handler := &NoticeHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/notices?year=2024", nil)
	w := httptest.NewRecorder()

	expected := []model.Notice{{ID: uuid.New(), Content: "Year Notice"}}
	mockService.EXPECT().
		GetNoticesByYear(2024).
		Return(expected, nil)

	handler.GetNoticesByMonthYear(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetNoticesByMonthYear_InvalidMonthFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := &NoticeHandler{Service: mocks.NewMockNoticeServiceInterface(ctrl)}

	req := httptest.NewRequest(http.MethodGet, "/notices?year=2024&month=abc", nil)
	w := httptest.NewRecorder()

	handler.GetNoticesByMonthYear(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetNoticesByMonthYear_MonthOutOfRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := &NoticeHandler{Service: mocks.NewMockNoticeServiceInterface(ctrl)}

	req := httptest.NewRequest(http.MethodGet, "/notices?year=2024&month=13", nil)
	w := httptest.NewRecorder()

	handler.GetNoticesByMonthYear(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetNoticesByMonthYear_MonthYearServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockNoticeServiceInterface(ctrl)
	handler := &NoticeHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/notices?year=2024&month=5", nil)
	w := httptest.NewRecorder()

	mockService.EXPECT().
		GetNoticesByMonthYear(time.Month(5), 2024).
		Return(nil, errors.New("not found"))

	handler.GetNoticesByMonthYear(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("got %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestGetNoticesByMonthYear_MonthYearSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockNoticeServiceInterface(ctrl)
	handler := &NoticeHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/notices?year=2024&month=5", nil)
	w := httptest.NewRecorder()

	expected := []model.Notice{{ID: uuid.New(), Content: "May Notice"}}
	mockService.EXPECT().
		GetNoticesByMonthYear(time.Month(5), 2024).
		Return(expected, nil)

	handler.GetNoticesByMonthYear(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}
