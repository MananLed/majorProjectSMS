package repository

import (

	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MananLed/majorProjectSMS/internal/model"
)

func newSocietyRepo(t *testing.T) (*SocietyRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	repo := NewSocietyRepository(db)
	cleanup := func() { db.Close() }
	return repo, mock, cleanup
}


func TestSocietyRepository_GetAllResidents_Success(t *testing.T) {
	repo, mock, cleanup := newSocietyRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id", "first_name", "middle_name", "last_name", "mobile_number", "email", "password", "role",
	}).AddRow("res-1", "John", "M", "Doe", "1234567890", "john@example.com", "secret", model.RoleResident)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, first_name, middle_name, last_name, mobile_number, email, password, role 
		 FROM users WHERE role = $1`,
	)).WithArgs(string(model.RoleResident)).WillReturnRows(rows)

	users, err := repo.GetAllResidents()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 resident, got %d", len(users))
	}
}

func TestSocietyRepository_GetAllResidents_Empty(t *testing.T) {
	repo, mock, cleanup := newSocietyRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id", "first_name", "middle_name", "last_name", "mobile_number", "email", "password", "role",
	})

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, first_name, middle_name, last_name, mobile_number, email, password, role 
		 FROM users WHERE role = $1`,
	)).WithArgs(string(model.RoleResident)).WillReturnRows(rows)

	users, err := repo.GetAllResidents()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("expected 0 residents, got %d", len(users))
	}
}

func TestSocietyRepository_GetAllResidents_QueryError(t *testing.T) {
	repo, mock, cleanup := newSocietyRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, first_name, middle_name, last_name, mobile_number, email, password, role 
		 FROM users WHERE role = $1`,
	)).WithArgs(string(model.RoleResident)).WillReturnError(errors.New("query failed"))

	_, err := repo.GetAllResidents()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestSocietyRepository_GetAllResidents_ScanError(t *testing.T) {
	repo, mock, cleanup := newSocietyRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id", "first_name",
	}).AddRow("res-1", "Bob")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, first_name, middle_name, last_name, mobile_number, email, password, role 
		 FROM users WHERE role = $1`,
	)).WithArgs(string(model.RoleResident)).WillReturnRows(rows)

	_, err := repo.GetAllResidents()
	if err == nil {
		t.Fatalf("expected scan error, got nil")
	}
}

func TestSocietyRepository_GetAllOfficers_Success(t *testing.T) {
	repo, mock, cleanup := newSocietyRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id", "first_name", "middle_name", "last_name", "mobile_number", "email", "password", "role",
	}).AddRow("off-1", "Alice", "", "Smith", "9876543210", "alice@example.com", "topsecret", model.RoleOfficer)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, first_name, middle_name, last_name, mobile_number, email, password, role 
		 FROM users WHERE role = $1`,
	)).WithArgs(string(model.RoleOfficer)).WillReturnRows(rows)

	users, err := repo.GetAllOfficers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 officer, got %d", len(users))
	}
}

func TestSocietyRepository_GetAllOfficers_Empty(t *testing.T) {
	repo, mock, cleanup := newSocietyRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id", "first_name", "middle_name", "last_name", "mobile_number", "email", "password", "role",
	})

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, first_name, middle_name, last_name, mobile_number, email, password, role 
		 FROM users WHERE role = $1`,
	)).WithArgs(string(model.RoleOfficer)).WillReturnRows(rows)

	users, err := repo.GetAllOfficers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("expected 0 officers, got %d", len(users))
	}
}

func TestSocietyRepository_GetAllOfficers_QueryError(t *testing.T) {
	repo, mock, cleanup := newSocietyRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, first_name, middle_name, last_name, mobile_number, email, password, role 
		 FROM users WHERE role = $1`,
	)).WithArgs(string(model.RoleOfficer)).WillReturnError(errors.New("query failed"))

	_, err := repo.GetAllOfficers()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestSocietyRepository_GetAllOfficers_ScanError(t *testing.T) {
	repo, mock, cleanup := newSocietyRepo(t)
	defer cleanup()

	// Provide only 3 columns instead of 8
	rows := sqlmock.NewRows([]string{
		"id", "first_name", "last_name",
	}).AddRow("off-1", "Alice", "Smith")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, first_name, middle_name, last_name, mobile_number, email, password, role 
		 FROM users WHERE role = $1`,
	)).WithArgs(string(model.RoleOfficer)).WillReturnRows(rows)

	_, err := repo.GetAllOfficers()
	if err == nil {
		t.Fatalf("expected scan error, got nil")
	}
}
