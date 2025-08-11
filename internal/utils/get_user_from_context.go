package utils

import (
	"context"
	"errors"

	"github.com/MananLed/majorProjectSMS/internal/model"
)

func GetUserFromContext(ctx context.Context) (*model.User, error) {
	id, ok := ctx.Value(UserIDKey).(string)
	if !ok || id == "" {
		return nil, errors.New("user ID not found in context")
	}

	role, ok := ctx.Value(UserRoleKey).(model.UserRole)
	if !ok {
		return nil, errors.New("user role not found in context")
	}

	password, ok := ctx.Value(UserPassKey).(string)
	if !ok {
		return nil, errors.New("user password not found in context")
	}

	return &model.User{
		ID:       id,
		Role:     role,
		Password: string(password),
	}, nil
}