package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MananLed/majorProjectSMS/internal/model"
)

func setupMock(t *testing.T) (*CredentialRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	repo := NewCredentialRepository(db)
	cleanup := func() { db.Close() }
	return repo, mock, cleanup
}

func TestDeleteUserByIDAndRole_Success(t *testing.T) {
	repo, mock, cleanup := setupMock(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND role = $2)`)).
		WithArgs("123", string(model.RoleResident)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1 AND role = $2`)).
		WithArgs("123", string(model.RoleResident)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteUserByIDAndRole("123", model.RoleResident)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestDeleteUserByIDAndRole_UserNotFound(t *testing.T) {
	repo, mock, cleanup := setupMock(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND role = $2)`)).
		WithArgs("123", string(model.RoleResident)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	err := repo.DeleteUserByIDAndRole("123", model.RoleResident)
	if err == nil || err.Error() != "user not found" {
		t.Errorf("expected 'user not found', got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestDeleteUserByIDAndRole_CheckExistenceError(t *testing.T) {
	repo, mock, cleanup := setupMock(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND role = $2)`)).
		WithArgs("123", string(model.RoleResident)).
		WillReturnError(errors.New("db error"))

	err := repo.DeleteUserByIDAndRole("123", model.RoleResident)
	if err == nil || err.Error() != "failed to check user existence: db error" {
		t.Errorf("expected error 'failed to check user existence', got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestDeleteUserByIDAndRole_DeleteError(t *testing.T) {
	repo, mock, cleanup := setupMock(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND role = $2)`)).
		WithArgs("123", string(model.RoleResident)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1 AND role = $2`)).
		WithArgs("123", string(model.RoleResident)).
		WillReturnError(errors.New("delete error"))

	err := repo.DeleteUserByIDAndRole("123", model.RoleResident)
	if err == nil || err.Error() != "failed to delete user: delete error" {
		t.Errorf("expected error 'failed to delete user', got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}