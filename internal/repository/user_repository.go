package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/MananLed/majorProjectSMS/constants"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryInterface interface {
	AddUser(user model.User) error
	GetUserByID(id string) (*model.User, error)
	UpdateUser(user model.User) error
	ChangePassword(id string, newHashedPassword string) error
	IsPasswordUnique(hashedPassword string) bool
	DeleteUserByID(id string) error
}

type UserRepository struct {
	mu sync.Mutex
}

func (r *UserRepository) LoadUsers() ([]model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(string(constants.UserDataPath))

	if err != nil {
		if os.IsNotExist(err) {
			return []model.User{}, nil
		}
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}

	var users []model.User
	err = json.Unmarshal(data, &users)
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) SaveUsers(users []model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := json.MarshalIndent(users, "", " ")
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}
	return os.WriteFile(string(constants.UserDataPath), data, 0644)
}

func (r *UserRepository) GetUserByID(id string) (*model.User, error) {
	users, err := r.LoadUsers()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}

	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *UserRepository) AddUser(newUser model.User) error {
	users, err := r.LoadUsers()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	for _, u := range users {
		if u.ID == newUser.ID {
			return errors.New("user with this ID already exists")
		}
	}

	users = append(users, newUser)
	return r.SaveUsers(users)
}

func (r *UserRepository) UpdateUser(updatedUser model.User) error {
	users, err := r.LoadUsers()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	updated := false
	for i, u := range users {
		if u.Password == updatedUser.Password {
			users[i] = updatedUser
			updated = true
			break
		}
	}

	if !updated {
		return errors.New("user not found")
	}

	return r.SaveUsers(users)
}

func (r *UserRepository) ChangePassword(id string, newHashedPassword string) error {
	users, err := r.LoadUsers()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	updated := false
	for i, u := range users {
		if u.ID == id {
			users[i].Password = newHashedPassword
			updated = true
			break
		}
	}

	if !updated {
		return errors.New("user not found")
	}

	return r.SaveUsers(users)
}

func (r *UserRepository) IsPasswordUnique(Password string) bool {
	users, err := r.LoadUsers()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return false
	}

	for _, user := range users {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Password))
		if err == nil {
			return false
		}
	}
	return true
}

func (r *UserRepository) DeleteUserByID(id string) error {
	users, err := r.LoadUsers()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	updatedUsers := []model.User{}
	found := false
	for _, user := range users {
		if user.ID != id {
			updatedUsers = append(updatedUsers, user)
		} else {
			found = true
		}
	}

	if !found {
		return errors.New("user not found")
	}

	return r.SaveUsers(updatedUsers)
}
