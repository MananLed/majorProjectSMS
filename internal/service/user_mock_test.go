package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepo struct {
	users map[string]model.User
}

func (m *MockUserRepo) AddUser(user model.User) error {
	if _, exists := m.users[user.ID]; exists {
		return errors.New("user already exists")
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepo) GetUserByIDAndPassword(id string, password string) (*model.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	if user.Password != password {
		return nil, errors.New("invalid credentials")
	}
	return &user, nil
}

func (m *MockUserRepo) UpdateUser(user model.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepo) ChangePassword(id string, newHashedPassword string) error {
	user, exists := m.users[id]
	if !exists {
		return errors.New("user not found")
	}
	user.Password = newHashedPassword
	m.users[id] = user
	return nil
}

func (m *MockUserRepo) IsPasswordUnique(hashedPassword string) bool {
	for _, user := range m.users {
		if user.Password == hashedPassword {
			return false
		}
	}
	return true
}

func (m *MockUserRepo) DeleteUserByID(id string) error {
	if _, exists := m.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(m.users, id)
	return nil
}

func (m *MockUserRepo) GetUserByID(id string) (*model.User, error) {
	if _, exists := m.users[id]; !exists {
		return nil, errors.New("user not found")
	}
	user := m.users[id]
	return &user, nil
}

//Tests

func TestSignUp(t *testing.T) {
	mockRepo := &MockUserRepo{users: make(map[string]model.User)}
	service := NewUserService(mockRepo)

	user := model.User{
		ID:       "man",
		Password: "nin",
		Role:     model.RoleResident,
	}

	err := service.SignUp(user)

	if err != nil {
		t.Errorf("expected signup to succeed, got error: %v", err)
	}
}

func TestLogin(t *testing.T) {
	mockRepo := &MockUserRepo{users: make(map[string]model.User)}
	service := NewUserService(mockRepo)
	user := model.User{
		ID:       "man",
		Password: "nin",
	}

	_ = service.SignUp(user)

	loggedInUser, err := service.Login("man", "nin")

	if err != nil || loggedInUser.ID != "man" {
		t.Errorf("login failed: %v", err)
	}
}

func TestLoginFailWrongPassword(t *testing.T) {
	mockRepo := &MockUserRepo{users: make(map[string]model.User)}
	service := NewUserService(mockRepo)

	user := model.User{ID: "mansi", Password: "nint"}
	_ = service.SignUp(user)

	_, err := service.Login("mansi", "sint")

	if err == nil {
		t.Error("expected login to fail due to wrong password")
	}
}

func TestChangePassword(t *testing.T) {
	mockRepo := &MockUserRepo{users: make(map[string]model.User)}
	service := NewUserService(mockRepo)

	rawPassword := "old123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)

	user := model.User{
		ID:       "man@example.com",
		Password: string(hashed),
		Role:     model.RoleResident,
	}

	mockRepo.AddUser(user)

	newPassword := "new123"

	err := service.ChangePassword(&user, rawPassword, newPassword)
	if err != nil {
		t.Errorf("Change Password failed: %v", err)
	}

	updatedUser, _ := mockRepo.GetUserByID(user.ID)
	err = bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte(newPassword))
	if err != nil {
		t.Errorf("Password was not updated correctly")
	}
}

func TestSignUp_Failure(t *testing.T) {
	mockRepo := &MockUserRepo{users: make(map[string]model.User)}
	service := NewUserService(mockRepo)

	user := model.User{ID: "", Password: ""}

	err := service.SignUp(user)
	if err == nil {
		t.Error("expected signup to fail for empty ID/password")
	}
}

func TestLoginFailUserNotFound(t *testing.T) {
	mockRepo := &MockUserRepo{users: make(map[string]model.User)}
	service := NewUserService(mockRepo)

	_, err := service.Login("ghost", "pass")
	if err == nil {
		t.Error("expected login to fail for non-existent user")
	}
}

func TestChangePassword_WrongCurrentPassword(t *testing.T) {
	mockRepo := &MockUserRepo{users: make(map[string]model.User)}
	service := NewUserService(mockRepo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correct123"), bcrypt.DefaultCost)
	user := model.User{ID: "man@example.com", Password: string(hashed)}
	mockRepo.AddUser(user)

	err := service.ChangePassword(&user, "wrong123", "new123")
	if err == nil {
		t.Error("expected error for wrong current password")
	}
}

func TestChangePassword_DuplicatePassword(t *testing.T) {
	mockRepo := &MockUserRepo{users: make(map[string]model.User)}
	service := NewUserService(mockRepo)

	existing := model.User{ID: "exist", Password: "duplicate"}
	mockRepo.AddUser(existing)

	user := model.User{ID: "u1", Password: "old"}
	mockRepo.AddUser(user)

	err := service.ChangePassword(&user, "old", "duplicate")
	if err == nil {
		t.Error("expected error for duplicate password")
	}
}

func TestUpdateProfile_Success(t *testing.T) {
	mockRepo := &MockUserRepo{users: make(map[string]model.User)}
	service := NewUserService(mockRepo)

	user := model.User{ID: "man", Password: "nin"}
	mockRepo.AddUser(user)

	updated := model.User{ID: "man", Password: "newpass"}
	err := service.UpdateProfile(updated)
	if err != nil {
		t.Errorf("expected update to succeed, got %v", err)
	}
}

func TestUpdateProfile_Failure(t *testing.T) {
	mockRepo := &MockUserRepo{users: make(map[string]model.User)}
	service := NewUserService(mockRepo)

	user := model.User{ID: "ghost", Password: "xxx"}
	err := service.UpdateProfile(user)
	if err == nil {
		t.Error("expected update to fail for non-existent user")
	}
}

func TestDeleteProfile_Success(t *testing.T) {
	mockRepo := &MockUserRepo{users: make(map[string]model.User)}
	service := NewUserService(mockRepo)

	user := model.User{ID: "u1", Password: "p1", Email: "dlfj", Flat: "sldjf", Role: model.RoleResident}
	mockRepo.AddUser(user)

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "resident")
	err := service.DeleteProfile(ctx)
	if err != nil {
		t.Errorf("expected delete to succeed, got %v", err)
	}
}

func TestDeleteProfile_Failure(t *testing.T) {
	mockRepo := &MockUserRepo{users: make(map[string]model.User)}
	service := NewUserService(mockRepo)

	ctx := context.Background()
	err := service.DeleteProfile(ctx)
	if err == nil {
		t.Error("expected delete to fail due to missing user in context")
	}
}

func TestGetUserByID_Success(t *testing.T) {
    mockRepo := &MockUserRepo{users: make(map[string]model.User)}
    service := NewUserService(mockRepo)

    user := model.User{ID: "u1", Password: "p1", Email: "dlf", Flat: "djf"}
    mockRepo.AddUser(user)

    ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)
	ctx = context.WithValue(ctx, utils.UserRoleKey, "resident")
    got, err := service.GetUserByID(ctx)
    if err != nil || got.ID != "u1" {
        t.Errorf("expected to get user u1, got %v, err: %v", got, err)
    }
}

func TestGetUserByID_Failure(t *testing.T) {
    mockRepo := &MockUserRepo{users: make(map[string]model.User)}
    service := NewUserService(mockRepo)

    ctx := context.Background()
    _, err := service.GetUserByID(ctx)
    if err == nil {
        t.Error("expected failure when no user in context")
    }
}