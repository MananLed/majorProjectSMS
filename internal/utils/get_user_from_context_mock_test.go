package utils

import (
	"context"
	"testing"

	"github.com/MananLed/majorProjectSMS/internal/model"
)

func TestGetUserFromContext_Success(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, UserIDKey, "user-123")
	ctx = context.WithValue(ctx, UserRoleKey, string(model.RoleResident))
	ctx = context.WithValue(ctx, UserEmailKey, "test@example.com")
	ctx = context.WithValue(ctx, UserFlatKey, "A-101")

	user, err := GetUserFromContext(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != "user-123" || user.Email != "test@example.com" || user.Flat != "A-101" {
		t.Errorf("unexpected user values: %+v", user)
	}
	if user.Role != model.RoleResident {
		t.Errorf("expected role Resident, got %s", user.Role)
	}
}

func TestGetUserFromContext_MissingID(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, UserRoleKey, string(model.RoleResident))
	ctx = context.WithValue(ctx, UserEmailKey, "test@example.com")
	ctx = context.WithValue(ctx, UserFlatKey, "A-101")

	_, err := GetUserFromContext(ctx)
	if err == nil {
		t.Fatal("expected error for missing ID, got nil")
	}
}

func TestGetUserFromContext_MissingRole(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, UserIDKey, "user-123")
	ctx = context.WithValue(ctx, UserEmailKey, "test@example.com")
	ctx = context.WithValue(ctx, UserFlatKey, "A-101")

	_, err := GetUserFromContext(ctx)
	if err == nil {
		t.Fatal("expected error for missing role, got nil")
	}
}

func TestGetUserFromContext_MissingEmail(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, UserIDKey, "user-123")
	ctx = context.WithValue(ctx, UserRoleKey, string(model.RoleResident))
	ctx = context.WithValue(ctx, UserFlatKey, "A-101")

	_, err := GetUserFromContext(ctx)
	if err == nil {
		t.Fatal("expected error for missing email, got nil")
	}
}

func TestGetUserFromContext_MissingFlat(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, UserIDKey, "user-123")
	ctx = context.WithValue(ctx, UserRoleKey, string(model.RoleResident))
	ctx = context.WithValue(ctx, UserEmailKey, "test@example.com")

	_, err := GetUserFromContext(ctx)
	if err == nil {
		t.Fatal("expected error for missing flat, got nil")
	}
}
