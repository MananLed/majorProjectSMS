package repository 

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/google/uuid"
)

func setupFeedbackMock(t *testing.T) (*FeedbackRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	repo := NewFeedbackRepository(db)
	cleanup := func() { db.Close() }
	return repo, mock, cleanup
}

func TestSaveFeedback_Success(t *testing.T) {
	repo, mock, cleanup := setupFeedbackMock(t)
	defer cleanup()

	f := model.Feedback{
		ID:         uuid.New(),
		ResidentID: "r1",
		Rating:     5,
		Content:    "Great service",
		Flat:       "101A",
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO feedbacks (id, resident_id, rating, content, flat_no)
		VALUES ($1, $2, $3, $4, $5)`)).
		WithArgs(f.ID, f.ResidentID, f.Rating, f.Content, f.Flat).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.SaveFeedback(f)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSaveFeedback_InsertError(t *testing.T) {
	repo, mock, cleanup := setupFeedbackMock(t)
	defer cleanup()

	f := model.Feedback{
		ID:         uuid.New(),
		ResidentID: "r1",
		Rating:     3,
		Content:    "Ok service",
		Flat:       "102",
	}

	mock.ExpectExec("INSERT INTO feedbacks").
		WithArgs(f.ID, f.ResidentID, f.Rating, f.Content, f.Flat).
		WillReturnError(errors.New("insert failed"))

	err := repo.SaveFeedback(f)
	if err == nil || err.Error() != "insert failed" {
		t.Errorf("expected insert failed error, got %v", err)
	}
}

func TestGetFeedbacksByID_Success(t *testing.T) {
	repo, mock, cleanup := setupFeedbackMock(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "resident_id", "rating", "content", "flat_no"}).
		AddRow(uuid.New(), "r1", 4, "Good", "101").
		AddRow(uuid.New(), "r1", 5, "Excellent", "101")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, resident_id, rating, content, flat_no
		FROM feedbacks
		WHERE resident_id = $1`)).
		WithArgs("r1").
		WillReturnRows(rows)

	result, err := repo.GetFeedbacksByID("r1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 feedbacks, got %d", len(result))
	}
}

func TestGetFeedbacksByID_NoRows(t *testing.T) {
	repo, mock, cleanup := setupFeedbackMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, resident_id, rating, content, flat_no").
		WithArgs("r1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "resident_id", "rating", "content", "flat_no"}))

	result, err := repo.GetFeedbacksByID("r1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 feedbacks, got %d", len(result))
	}
}

func TestGetFeedbacksByID_QueryError(t *testing.T) {
	repo, mock, cleanup := setupFeedbackMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, resident_id, rating, content, flat_no").
		WithArgs("r1").
		WillReturnError(errors.New("db error"))

	_, err := repo.GetFeedbacksByID("r1")
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}

func TestGetFeedbacksByID_ScanError(t *testing.T) {
	repo, mock, cleanup := setupFeedbackMock(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "resident_id", "rating", "content", "flat_no"}).
		AddRow(uuid.New(), "r1", "bad-int", "Nice", "103")

	mock.ExpectQuery("SELECT id, resident_id, rating, content, flat_no").
		WithArgs("r1").
		WillReturnRows(rows)

	_, err := repo.GetFeedbacksByID("r1")
	if err == nil {
		t.Errorf("expected scan error, got nil")
	}
}

func TestGetAllFeedbacks_Success(t *testing.T) {
	repo, mock, cleanup := setupFeedbackMock(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "resident_id", "rating", "content", "flat_no"}).
		AddRow(uuid.New(), "r2", 2, "Not good", "104A").
		AddRow(uuid.New(), "r3", 5, "Perfect", "105B")

	mock.ExpectQuery("SELECT id, resident_id, rating, content, flat_no").
		WillReturnRows(rows)

	result, err := repo.GetAllFeedbacks()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 feedbacks, got %d", len(result))
	}
}

func TestGetAllFeedbacks_NoRows(t *testing.T) {
	repo, mock, cleanup := setupFeedbackMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, resident_id, rating, content, flat_no").
		WillReturnRows(sqlmock.NewRows([]string{"id", "resident_id", "rating", "content", "flat_no"}))

	result, err := repo.GetAllFeedbacks()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 feedbacks, got %d", len(result))
	}
}

func TestGetAllFeedbacks_QueryError(t *testing.T) {
	repo, mock, cleanup := setupFeedbackMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, resident_id, rating, content, flat_no").
		WillReturnError(errors.New("db error"))

	_, err := repo.GetAllFeedbacks()
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}

func TestGetAllFeedbacks_ScanError(t *testing.T) {
	repo, mock, cleanup := setupFeedbackMock(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "resident_id", "rating", "content", "flat_no"}).
		AddRow(uuid.New(), "r2", "bad-int", "Text", "106C")

	mock.ExpectQuery("SELECT id, resident_id, rating, content, flat_no").
		WillReturnRows(rows)

	_, err := repo.GetAllFeedbacks()
	if err == nil {
		t.Errorf("expected scan error, got nil")
	}
}