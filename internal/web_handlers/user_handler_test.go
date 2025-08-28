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

	"github.com/MananLed/majorProjectSMS/internal/mocks"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"go.uber.org/mock/gomock"
)

func newRecorderAndPost(body string) (*httptest.ResponseRecorder, *http.Request) {
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	return rec, req
}

func validBody() string {
	return `{
		"firstname":"John",
		"lastname":"Doe",
		"middlename":"X",
		"email":"john.doe@example.com",
		"password":"Password123!",
		"mobile":"9876543210",
		"flat":"101"
	}`
}

func TestSignUp_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	us := mocks.NewMockUserServiceInterface(ctrl)
	h := NewUserHandler(us, nil)

	req := httptest.NewRequest(http.MethodGet, "/signup", nil)
	rec := httptest.NewRecorder()

	h.SignUp(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestSignUp_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	us := mocks.NewMockUserServiceInterface(ctrl)
	h := NewUserHandler(us, nil)

	rec, req := newRecorderAndPost("{this is not json}")
	h.SignUp(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestSignUp_ValidationFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	us := mocks.NewMockUserServiceInterface(ctrl)
	h := NewUserHandler(us, nil)

	body := `{
		"firstname":"John",
		"lastname":"Doe",
		"middlename":"X",
		"email":"not-an-email",
		"password":"Password123!",
		"mobile":"9876543210",
		"flat":"101"
	}`

	rec, req := newRecorderAndPost(body)
	h.SignUp(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestSignUp_PasswordNotUnique_ShouldNotCallSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	us := mocks.NewMockUserServiceInterface(ctrl)
	h := NewUserHandler(us, nil)

	us.EXPECT().IsPasswordUnique("Password123!").Return(false)

	rec, req := newRecorderAndPost(validBody())
	h.SignUp(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestSignUp_ServiceSignUpError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	us := mocks.NewMockUserServiceInterface(ctrl)
	h := NewUserHandler(us, nil)

	us.EXPECT().IsPasswordUnique("Password123!").Return(true)
	us.EXPECT().SignUp(gomock.Any()).Return(errors.New("db error"))

	rec, req := newRecorderAndPost(validBody())
	h.SignUp(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestSignUp_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	us := mocks.NewMockUserServiceInterface(ctrl)
	h := NewUserHandler(us, nil)

	us.EXPECT().IsPasswordUnique("Password123!").Return(true)
	us.EXPECT().SignUp(gomock.Any()).Return(nil)

	rec, req := newRecorderAndPost(validBody())
	h.SignUp(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestLogin_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockUserService}

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got status %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestLogin_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockUserService}

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{invalid json}`))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockUserService}

	mockUserService.EXPECT().
		Login("john.doe@example.com", "wrong").
		Return(nil, errors.New("invalid"))

	req := httptest.NewRequest(http.MethodPost, "/login",
		strings.NewReader(`{"email":"john.doe@example.com","password":"wrong"}`))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got status %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockUserService}

	user := &model.User{
		ID:    "123",
		Role:  "resident",
		Email: "john.doe@example.com",
		Flat:  "101",
	}

	mockUserService.EXPECT().
		Login("john.doe@example.com", "Password123!").
		Return(user, nil)

	req := httptest.NewRequest(http.MethodPost, "/login",
		strings.NewReader(`{"email":"john.doe@example.com","password":"Password123!"}`))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("got status %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestViewProfile_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockUserService}

	req := httptest.NewRequest(http.MethodPost, "/profile", nil) // wrong method
	w := httptest.NewRecorder()

	handler.ViewProfile(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestViewProfile_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockUserService}

	mockUserService.EXPECT().
		GetUserByID(gomock.Any()).
		Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	w := httptest.NewRecorder()

	handler.ViewProfile(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestViewProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockUserService}

	user := &model.User{
		ID:    "123",
		Email: "john.doe@example.com",
		Role:  "resident",
		Flat:  "101",
	}

	mockUserService.EXPECT().
		GetUserByID(gomock.Any()).
		Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	w := httptest.NewRecorder()

	handler.ViewProfile(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestUpdateProfile_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockUserService}

	req := httptest.NewRequest(http.MethodGet, "/profile/update", nil)
	w := httptest.NewRecorder()

	handler.UpdateProfile(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestUpdateProfile_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockUserService}

	mockUserService.EXPECT().
		GetUserByID(gomock.Any()).
		Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodPatch, "/profile/update", nil)
	w := httptest.NewRecorder()

	handler.UpdateProfile(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestUpdateProfile_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockUserService}

	user := &model.User{ID: "123", FirstName: "John"}
	mockUserService.EXPECT().
		GetUserByID(gomock.Any()).
		Return(user, nil)

	req := httptest.NewRequest(http.MethodPatch, "/profile/update", strings.NewReader(`{invalid json}`))
	w := httptest.NewRecorder()

	handler.UpdateProfile(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUpdateProfile_UpdateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockUserService}

	user := &model.User{ID: "123", FirstName: "John"}
	mockUserService.EXPECT().
		GetUserByID(gomock.Any()).
		Return(user, nil)

	updatedUser := model.User{ID: "123", FirstName: "Jane"}
	mockUserService.EXPECT().
		UpdateProfile(updatedUser).
		Return(errors.New("db error"))

	body := `{"firstName":"Jane"}`
	req := httptest.NewRequest(http.MethodPatch, "/profile/update", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.UpdateProfile(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestUpdateProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockUserService}

	user := &model.User{ID: "123", FirstName: "John"}
	mockUserService.EXPECT().
		GetUserByID(gomock.Any()).
		Return(user, nil)

	updatedUser := *user
	updatedUser.FirstName = "Jane"
	mockUserService.EXPECT().
		UpdateProfile(updatedUser).
		Return(nil)

	body := `{"firstName":"Jane"}`
	req := httptest.NewRequest(http.MethodPatch, "/profile/update", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.UpdateProfile(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err == nil {
		if resp["message"] != "User updated successfully" {
			t.Errorf("unexpected message: %v", resp["message"])
		}
	}
}

func TestChangePassword_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/change-password", nil)
	w := httptest.NewRecorder()

	handler.ChangePassword(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestChangePassword_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockService}

	mockService.EXPECT().
		GetUserByID(gomock.Any()).
		Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodPatch, "/change-password", nil)
	w := httptest.NewRecorder()

	handler.ChangePassword(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestChangePassword_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockService}

	user := &model.User{ID: "123", FirstName: "John"}
	mockService.EXPECT().
		GetUserByID(gomock.Any()).
		Return(user, nil)

	req := httptest.NewRequest(http.MethodPatch, "/change-password", strings.NewReader(`{bad json}`))
	w := httptest.NewRecorder()

	handler.ChangePassword(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestChangePassword_ChangeFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockService}

	user := &model.User{ID: "123", FirstName: "John"}
	mockService.EXPECT().
		GetUserByID(gomock.Any()).
		Return(user, nil)

	mockService.EXPECT().
		ChangePassword(user, "oldpass", "newpass").
		Return(errors.New("invalid password"))

	body := `{"oldPassword":"oldpass","newPassword":"newpass"}`
	req := httptest.NewRequest(http.MethodPatch, "/change-password", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.ChangePassword(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestChangePassword_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockService}

	user := &model.User{ID: "123", FirstName: "John"}
	mockService.EXPECT().
		GetUserByID(gomock.Any()).
		Return(user, nil)

	mockService.EXPECT().
		ChangePassword(user, "oldpass", "newpass").
		Return(nil)

	body := `{"oldPassword":"oldpass","newPassword":"newpass"}`
	req := httptest.NewRequest(http.MethodPatch, "/change-password", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.ChangePassword(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err == nil {
		if resp["message"] != "User password updated successfully" {
			t.Errorf("unexpected message: %v", resp["message"])
		}
	}
}

func TestDeleteProfile_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/delete", nil) // wrong method
	w := httptest.NewRecorder()

	handler.DeleteProfile(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestDeleteProfile_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockService}

	req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
	w := httptest.NewRecorder()

	handler.DeleteProfile(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestDeleteProfile_DeleteFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockService}

	user := &model.User{ID: "123", FirstName: "John"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "resident")
	req := httptest.NewRequest(http.MethodDelete, "/delete", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().
		DeleteProfile(ctx).
		Return(errors.New("db error"))

	handler.DeleteProfile(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestDeleteProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockService}

	user := &model.User{ID: "123", FirstName: "John"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "resident")

	req := httptest.NewRequest(http.MethodDelete, "/delete", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().
		DeleteProfile(ctx).
		Return(nil)

	handler.DeleteProfile(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCreateOfficer_InvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/officer", nil)
	w := httptest.NewRecorder()

	handler.CreateOfficer(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}


func TestCreateOfficer_UnauthorizedUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockService}

	req := httptest.NewRequest(http.MethodPost, "/officer", bytes.NewBufferString(`{"email":"o@example.com","password":"pass"}`))
	w := httptest.NewRecorder()

	handler.CreateOfficer(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want %d", w.Code, http.StatusForbidden)
	}
}


func TestCreateOfficer_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	req := httptest.NewRequest(http.MethodPost, "/officer", bytes.NewBufferString(`{invalid-json}`)).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.CreateOfficer(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateOfficer_SignUpFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	body := `{"email":"officer@example.com","password":"Password123!"}`
	req := httptest.NewRequest(http.MethodPost, "/officer", bytes.NewBufferString(body)).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().
		SignUp(gomock.Any()).
		Return(errors.New("db error"))

	handler.CreateOfficer(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestCreateOfficer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserServiceInterface(ctrl)
	handler := &UserHandler{UserService: mockService}

	user := &model.User{ID: "admin1", Role: model.RoleAdmin, Email: "dlfj", Flat: "xxx"}

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "admin")

	body := `{"email":"officer@example.com","password":"Password123!"}`
	req := httptest.NewRequest(http.MethodPost, "/officer", bytes.NewBufferString(body)).WithContext(ctx)
	w := httptest.NewRecorder()

	mockService.EXPECT().
		SignUp(gomock.Any()).
		Return(nil)

	handler.CreateOfficer(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}