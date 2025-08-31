package repository 

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/google/uuid"
)

func setupInvoiceMock(t *testing.T) (*InvoiceRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	repo := NewInvoiceRepository(db)
	cleanup := func() { db.Close() }
	return repo, mock, cleanup
}

func TestSaveInvoice_Success(t *testing.T) {
	repo, mock, cleanup := setupInvoiceMock(t)
	defer cleanup()

	invoice := model.Invoice{Month: 1, Year: 2025, Amount: 1000}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO invoices (id, month, year, amount)
		VALUES ($1, $2, $3, $4)`)).
		WithArgs(sqlmock.AnyArg(), invoice.Month, invoice.Year, invoice.Amount).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.SaveInvoice(invoice)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSaveInvoice_Error(t *testing.T) {
	repo, mock, cleanup := setupInvoiceMock(t)
	defer cleanup()

	invoice := model.Invoice{Month: 2, Year: 2025, Amount: 2000}

	mock.ExpectExec("INSERT INTO invoices").
		WillReturnError(errors.New("insert failed"))

	err := repo.SaveInvoice(invoice)
	if err == nil || err.Error() != "insert failed" {
		t.Errorf("expected insert failed error, got %v", err)
	}
}

func TestGetInvoiceByMonthAndYear_Success(t *testing.T) {
	repo, mock, cleanup := setupInvoiceMock(t)
	defer cleanup()

	id := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "month", "year", "amount"}).
		AddRow(id, 3, 2025, 3000)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, month, year, amount
		FROM invoices
		WHERE month = $1 AND year = $2`)).
		WithArgs(3, 2025).
		WillReturnRows(rows)

	result, err := repo.GetInvoiceByMonthAndYear(time.March, 2025)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result == nil || result.Amount != 3000 {
		t.Errorf("expected invoice with amount 3000, got %v", result)
	}
}

func TestGetInvoiceByMonthAndYear_NotFound(t *testing.T) {
	repo, mock, cleanup := setupInvoiceMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, month, year, amount").
		WithArgs(4, 2025).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetInvoiceByMonthAndYear(time.April, 2025)
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
	if err == nil || err.Error() != "invoice not found" {
		t.Errorf("expected invoice not found error, got %v", err)
	}
}

func TestGetInvoiceByMonthAndYear_ScanError(t *testing.T) {
	repo, mock, cleanup := setupInvoiceMock(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "month", "year", "amount"}).
		AddRow("bad-uuid", "wrong", "wrong", "wrong")

	mock.ExpectQuery("SELECT id, month, year, amount").
		WithArgs(5, 2025).
		WillReturnRows(rows)

	_, err := repo.GetInvoiceByMonthAndYear(time.May, 2025)
	if err == nil {
		t.Errorf("expected scan error, got nil")
	}
}

func TestGetInvoicesByYear_Success(t *testing.T) {
	repo, mock, cleanup := setupInvoiceMock(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "month", "year", "amount"}).
		AddRow(uuid.New(), 6, 2025, 6000).
		AddRow(uuid.New(), 7, 2025, 7000)

	mock.ExpectQuery("SELECT id, month, year, amount").
		WithArgs(2025).
		WillReturnRows(rows)

	results, err := repo.GetInvoicesByYear(2025)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 invoices, got %d", len(results))
	}
}

func TestGetInvoicesByYear_NoRows(t *testing.T) {
	repo, mock, cleanup := setupInvoiceMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, month, year, amount").
		WithArgs(2025).
		WillReturnRows(sqlmock.NewRows([]string{"id", "month", "year", "amount"}))

	results, err := repo.GetInvoicesByYear(2025)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 invoices, got %d", len(results))
	}
}

func TestGetInvoicesByYear_QueryError(t *testing.T) {
	repo, mock, cleanup := setupInvoiceMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, month, year, amount").
		WithArgs(2026).
		WillReturnError(errors.New("db error"))

	_, err := repo.GetInvoicesByYear(2026)
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}

func TestGetInvoicesByYear_ScanError(t *testing.T) {
	repo, mock, cleanup := setupInvoiceMock(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "month", "year", "amount"}).
		AddRow("bad-uuid", "wrong", "wrong", "wrong")

	mock.ExpectQuery("SELECT id, month, year, amount").
		WithArgs(2025).
		WillReturnRows(rows)

	_, err := repo.GetInvoicesByYear(2025)
	if err == nil {
		t.Errorf("expected scan error, got nil")
	}
}