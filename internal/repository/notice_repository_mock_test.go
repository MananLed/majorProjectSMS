package repository 

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/google/uuid"
)

func setupNoticeMock(t *testing.T) (*NoticeRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	repo := NewNoticeRepository(db)
	cleanup := func() { db.Close() }
	return repo, mock, cleanup
}

func TestSaveNotice_Success(t *testing.T) {
	repo, mock, cleanup := setupNoticeMock(t)
	defer cleanup()

	notice := model.Notice{DateIssued: time.Now(), Content: "Test", Month: 8, Year: 2025}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO notices (id, date_issued, content, month, year)
		VALUES ($1, $2, $3, $4, $5)`)).
		WithArgs(sqlmock.AnyArg(), notice.DateIssued, notice.Content, notice.Month, notice.Year).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.SaveNotice(notice)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSaveNotice_Error(t *testing.T) {
	repo, mock, cleanup := setupNoticeMock(t)
	defer cleanup()

	notice := model.Notice{DateIssued: time.Now(), Content: "BadInsert", Month: 9, Year: 2025}

	mock.ExpectExec("INSERT INTO notices").
		WillReturnError(errors.New("insert failed"))

	err := repo.SaveNotice(notice)
	if err == nil || err.Error() != "insert failed" {
		t.Errorf("expected insert failed, got %v", err)
	}
}

func TestGetAllNotices_Success(t *testing.T) {
	repo, mock, cleanup := setupNoticeMock(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "date_issued", "content", "month", "year"}).
		AddRow(uuid.New(), time.Now(), "Notice 1", 8, 2025).
		AddRow(uuid.New(), time.Now(), "Notice 2", 8, 2025)

	mock.ExpectQuery("SELECT id, date_issued, content, month, year").
		WillReturnRows(rows)

	results, err := repo.GetAllNotices()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 notices, got %d", len(results))
	}
}

func TestGetAllNotices_NoRows(t *testing.T) {
	repo, mock, cleanup := setupNoticeMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, date_issued, content, month, year").
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_issued", "content", "month", "year"}))

	results, err := repo.GetAllNotices()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 notices, got %d", len(results))
	}
}

func TestGetAllNotices_QueryError(t *testing.T) {
	repo, mock, cleanup := setupNoticeMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, date_issued, content, month, year").
		WillReturnError(errors.New("db error"))

	_, err := repo.GetAllNotices()
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}

func TestGetAllNotices_ScanError(t *testing.T) {
	repo, mock, cleanup := setupNoticeMock(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "date_issued", "content", "month", "year"}).
		AddRow("bad-uuid", "bad-date", 123, "wrong", "wrong")

	mock.ExpectQuery("SELECT id, date_issued, content, month, year").
		WillReturnRows(rows)

	_, err := repo.GetAllNotices()
	if err == nil {
		t.Errorf("expected scan error, got nil")
	}
}

func TestGetNoticesByMonthYear_Success(t *testing.T) {
	repo, mock, cleanup := setupNoticeMock(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "date_issued", "content", "month", "year"}).
		AddRow(uuid.New(), time.Now(), "NoticeMonthYear", 8, 2025)

	mock.ExpectQuery("SELECT id, date_issued, content, month, year").
		WithArgs(8, 2025).
		WillReturnRows(rows)

	results, err := repo.GetNoticesByMonthYear(8, 2025)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 notice, got %d", len(results))
	}
}

func TestGetNoticesByMonthYear_NoRows(t *testing.T) {
	repo, mock, cleanup := setupNoticeMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, date_issued, content, month, year").
		WithArgs(9, 2025).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_issued", "content", "month", "year"}))

	results, err := repo.GetNoticesByMonthYear(9, 2025)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 notices, got %d", len(results))
	}
}

func TestGetNoticesByMonthYear_QueryError(t *testing.T) {
	repo, mock, cleanup := setupNoticeMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, date_issued, content, month, year").
		WithArgs(10, 2025).
		WillReturnError(errors.New("db error"))

	_, err := repo.GetNoticesByMonthYear(10, 2025)
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}

func TestGetNoticesByMonthYear_ScanError(t *testing.T) {
	repo, mock, cleanup := setupNoticeMock(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "date_issued", "content", "month", "year"}).
		AddRow("bad-uuid", "bad-date", 123, "wrong", "wrong")

	mock.ExpectQuery("SELECT id, date_issued, content, month, year").
		WithArgs(11, 2025).
		WillReturnRows(rows)

	_, err := repo.GetNoticesByMonthYear(11, 2025)
	if err == nil {
		t.Errorf("expected scan error, got nil")
	}
}

func TestGetNoticesByYear_Success(t *testing.T) {
	repo, mock, cleanup := setupNoticeMock(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "date_issued", "content", "month", "year"}).
		AddRow(uuid.New(), time.Now(), "NoticeYear", 8, 2025)

	mock.ExpectQuery("SELECT id, date_issued, content, month, year").
		WithArgs(2025).
		WillReturnRows(rows)

	results, err := repo.GetNoticesByYear(2025)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 notice, got %d", len(results))
	}
}

func TestGetNoticesByYear_NoRows(t *testing.T) {
	repo, mock, cleanup := setupNoticeMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, date_issued, content, month, year").
		WithArgs(2026).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_issued", "content", "month", "year"}))

	results, err := repo.GetNoticesByYear(2026)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 notices, got %d", len(results))
	}
}

func TestGetNoticesByYear_QueryError(t *testing.T) {
	repo, mock, cleanup := setupNoticeMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, date_issued, content, month, year").
		WithArgs(2027).
		WillReturnError(errors.New("db error"))

	_, err := repo.GetNoticesByYear(2027)
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}

func TestGetNoticesByYear_ScanError(t *testing.T) {
	repo, mock, cleanup := setupNoticeMock(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "date_issued", "content", "month", "year"}).
		AddRow("bad-uuid", "bad-date", 123, "wrong", "wrong")

	mock.ExpectQuery("SELECT id, date_issued, content, month, year").
		WithArgs(2028).
		WillReturnRows(rows)

	_, err := repo.GetNoticesByYear(2028)
	if err == nil {
		t.Errorf("expected scan error, got nil")
	}
}