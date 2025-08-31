package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"golang.org/x/crypto/bcrypt"
)


func newUserRepo(t *testing.T) (*UserRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock db: %v", err)
	}
	repo := NewUserRepository(db)
	return repo, mock, func() { db.Close() }
}


func TestUserRepository_AddUser_Success(t *testing.T) {
	repo, mock, cleanup := newUserRepo(t)
	defer cleanup()

	user := model.User{
		ID: "u-1", FirstName: "John", Email: "john@example.com",
		Password: "hashedpass", Role: model.RoleResident, Flat: "A-101",
	}

	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO users (id, first_name, middle_name, last_name, mobile_number, email, password, role, flat_no)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`)).
		WithArgs(user.ID, user.FirstName, user.MiddleName, user.LastName, user.MobileNumber,
			user.Email, user.Password, user.Role, user.Flat).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.AddUser(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUserRepository_AddUser_Error(t *testing.T) {
	repo, mock, cleanup := newUserRepo(t)
	defer cleanup()

	user := model.User{ID: "u-1"}
	mock.ExpectExec("INSERT INTO users").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("insert failed"))

	err := repo.AddUser(user)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestUserRepository_GetUserByIDAndPassword_Success(t *testing.T) {
	repo, mock, cleanup := newUserRepo(t)
	defer cleanup()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	rows := sqlmock.NewRows([]string{
		"id", "first_name", "middle_name", "last_name", "mobile_number", "email", "password", "role", "flat_no",
	}).AddRow("u-1", "John", "", "Doe", "1234567890", "john@example.com", string(hashed), model.RoleResident, "A-101")

	mock.ExpectQuery("SELECT id, first_name").
		WithArgs("john@example.com").
		WillReturnRows(rows)

	user, err := repo.GetUserByIDAndPassword("john@example.com", "secret")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID != "u-1" {
		t.Errorf("expected id u-1, got %v", user.ID)
	}
}

func TestUserRepository_GetUserByIDAndPassword_WrongPassword(t *testing.T) {
	repo, mock, cleanup := newUserRepo(t)
	defer cleanup()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	rows := sqlmock.NewRows([]string{
		"id", "first_name", "middle_name", "last_name", "mobile_number", "email", "password", "role", "flat_no",
	}).AddRow("u-1", "John", "", "Doe", "1234567890", "john@example.com", string(hashed), model.RoleResident, "A-101")

	mock.ExpectQuery("SELECT id, first_name").
		WithArgs("john@example.com").
		WillReturnRows(rows)

	_, err := repo.GetUserByIDAndPassword("john@example.com", "wrongpass")
	if err == nil {
		t.Fatalf("expected error for wrong password, got nil")
	}
}

func TestUserRepository_GetUserByIDAndPassword_QueryError(t *testing.T) {
	repo, mock, cleanup := newUserRepo(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, first_name").
		WithArgs("john@example.com").
		WillReturnError(errors.New("db error"))

	_, err := repo.GetUserByIDAndPassword("john@example.com", "secret")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestUserRepository_UpdateUser_Success(t *testing.T) {
	repo, mock, cleanup := newUserRepo(t)
	defer cleanup()

	user := model.User{ID: "u-1", FirstName: "John"}
	mock.ExpectExec("UPDATE users").
		WithArgs(user.FirstName, user.MiddleName, user.LastName, user.MobileNumber,
			user.Email, user.Password, user.Role, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdateUser(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUserRepository_UpdateUser_NoRows(t *testing.T) {
	repo, mock, cleanup := newUserRepo(t)
	defer cleanup()

	user := model.User{ID: "u-1"}
	mock.ExpectExec("UPDATE users").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 0))

	err := repo.UpdateUser(user)
	if err == nil {
		t.Fatalf("expected user not found error, got nil")
	}
}

func TestUserRepository_ChangePassword_Success(t *testing.T) {
	repo, mock, cleanup := newUserRepo(t)
	defer cleanup()

	mock.ExpectExec("UPDATE users").
		WithArgs("newhash", "u-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.ChangePassword("u-1", "newhash")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUserRepository_ChangePassword_NoRows(t *testing.T) {
	repo, mock, cleanup := newUserRepo(t)
	defer cleanup()

	mock.ExpectExec("UPDATE users").
		WithArgs("newhash", "u-1").
		WillReturnResult(sqlmock.NewResult(1, 0))

	err := repo.ChangePassword("u-1", "newhash")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestUserRepository_IsPasswordUnique_False(t *testing.T) {
	repo, mock, cleanup := newUserRepo(t)
	defer cleanup()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	rows := sqlmock.NewRows([]string{"password"}).AddRow(string(hashed))
	mock.ExpectQuery("SELECT password FROM users").WillReturnRows(rows)

	ok := repo.IsPasswordUnique("secret")
	if ok {
		t.Fatalf("expected false, got true")
	}
}

func TestUserRepository_IsPasswordUnique_True(t *testing.T) {
	repo, mock, cleanup := newUserRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"password"}).AddRow("some-other-hash")
	mock.ExpectQuery("SELECT password FROM users").WillReturnRows(rows)

	ok := repo.IsPasswordUnique("secret")
	if !ok {
		t.Fatalf("expected true, got false")
	}
}

func TestUserRepository_DeleteUserByID_Success(t *testing.T) {
	repo, mock, cleanup := newUserRepo(t)
	defer cleanup()

	mock.ExpectExec("DELETE FROM users").
		WithArgs("u-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteUserByID("u-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUserRepository_DeleteUserByID_NoRows(t *testing.T) {
	repo, mock, cleanup := newUserRepo(t)
	defer cleanup()

	mock.ExpectExec("DELETE FROM users").
		WithArgs("u-1").
		WillReturnResult(sqlmock.NewResult(1, 0))

	err := repo.DeleteUserByID("u-1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestUserRepository_GetUserByID_Success(t *testing.T) {
	repo, mock, cleanup := newUserRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id", "first_name", "middle_name", "last_name", "mobile_number", "email", "password", "role", "flat_no",
	}).AddRow("u-1", "John", "", "Doe", "1234567890", "john@example.com", "hash", model.RoleResident, "A-101")

	mock.ExpectQuery("SELECT id, first_name").
		WithArgs("u-1").
		WillReturnRows(rows)

	user, err := repo.GetUserByID("u-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID != "u-1" {
		t.Errorf("expected u-1, got %v", user.ID)
	}
}

func TestUserRepository_GetUserByID_NotFound(t *testing.T) {
	repo, mock, cleanup := newUserRepo(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, first_name").
		WithArgs("u-1").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetUserByID("u-1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
