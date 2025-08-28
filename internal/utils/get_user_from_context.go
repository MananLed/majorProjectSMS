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

	role, ok := ctx.Value(UserRoleKey).(string)
	if !ok {
		return nil, errors.New("user role not found in context")
	}

	email, ok := ctx.Value(UserEmailKey).(string)
	if !ok {
		return nil, errors.New("user email not found in context")
	}

	flat, ok := ctx.Value(UserFlatKey).(string)
	if !ok {
		return nil, errors.New("user flat not found in context")
	}

	return &model.User{
		ID:       id,
		Role:     model.UserRole(role),
		Email:    email,
		Flat:     flat,
	}, nil
}
