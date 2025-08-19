package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/repository"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
	"github.com/fatih/color"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface{
	SignUp(user model.User) error
	Login(id string, password string) (*model.User, error)
	UpdateProfile(user model.User) error
	ChangePassword(ctx context.Context, currentPassword string, newPassword string) error
	IsPasswordUnique(Password string) bool
	DeleteProfile(ctx context.Context) error
}

type UserService struct {
	UserRepo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{
		UserRepo: repo,
	}
}

func (us *UserService) SignUp(user model.User) error {
	if strings.TrimSpace(user.ID) == "" || strings.TrimSpace(user.Password) == "" {
		return errors.New("ID and Password cannot be empty")
	}
	return us.UserRepo.AddUser(user)
}

func (us *UserService) Login(id string, password string) (*model.User, error) {
	user, err := us.UserRepo.GetUserByID(id)
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		color.Red("Incorrect password.")
		return nil, err
	}
	return user, nil
}

func (us *UserService) UpdateProfile(user model.User) error {
	return us.UserRepo.UpdateUser(user)
}

func (us *UserService) ChangePassword(ctx context.Context, currentPassword string, newPassword string) error {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		color.Red("Unauthorized access.")
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	if !us.UserRepo.IsPasswordUnique(string(hashedPassword)) {
		return errors.New("password is already used by another user")
	}

	return us.UserRepo.ChangePassword(user.ID, string(hashedPassword))
}

func (us *UserService) IsPasswordUnique(Password string) bool {
	return us.UserRepo.IsPasswordUnique(Password)
}

func (us *UserService) DeleteProfile(ctx context.Context) error {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}
	return us.UserRepo.DeleteUserByID(user.ID)
}