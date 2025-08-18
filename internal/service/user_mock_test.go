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

func (m *MockUserRepo) GetUserByID(id string) (*model.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
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
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("nin"), bcrypt.DefaultCost)
	user := model.User{
	ID:       "man",
	Password: string(hashedPassword),
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

	user := model.User{ID: "man", Password: "nin"}

	_ = service.SignUp(user)

	_, err := service.Login("man", "sin")

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

	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, user.Role)
	ctx = context.WithValue(ctx, utils.UserPassKey, user.Password)

	newPassword := "new123"

	err := service.ChangePassword(ctx, rawPassword, newPassword)
	if err != nil {
		t.Errorf("Change Password failed: %v", err)
	}

	updatedUser, _ := mockRepo.GetUserByID(user.ID)
	err = bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte(newPassword))
	if err != nil {
		t.Errorf("Password was not updated correctly")
	}
}

