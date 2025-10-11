package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
)

type InvoiceRepositoryInterface interface {
	SaveInvoice(model.Invoice) error
	GetInvoiceByMonthAndYear(time.Month, int) (*model.Invoice, error)
	GetInvoicesByYear(int) ([]model.Invoice, error)
}

type InvoiceRepository struct {
	mu sync.Mutex
	DB *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{DB: db}
}

func (r *InvoiceRepository) SaveInvoice(invoice model.Invoice) error {
	invoice.ID = utils.GenerateUUID()

	query := `
		INSERT INTO invoices (id, month, year, amount)
		VALUES ($1, $2, $3, $4)
	`


	_, err := r.DB.Exec(query, invoice.ID, invoice.Month, invoice.Year, invoice.Amount)

	
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	return nil
}

func (r *InvoiceRepository) GetInvoiceByMonthAndYear(month time.Month, year int) (*model.Invoice, error) {
	var invoice model.Invoice

	query := `
		SELECT id, month, year, amount
		FROM invoices
		WHERE month = $1 AND year = $2
	`


	err := r.DB.QueryRow(query, int(month), year).Scan(&invoice.ID, &invoice.Month, &invoice.Year, &invoice.Amount)


	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invoice not found")
		}
		logger.LogToFile(fmt.Sprintf("error fetching invoice: %v", err))
		return nil, err
	}

	return &invoice, nil
}

func (r *InvoiceRepository) GetInvoicesByYear(year int) ([]model.Invoice, error) {
	var query string
	
	if year == 0{
		query = `
		SELECT id, month, year, amount
		FROM invoices
	`
	}else{
	query = `
		SELECT id, month, year, amount
		FROM invoices
		WHERE year = $1
	`}

	var rows *sql.Rows 
	var err error
	if year == 0{
		rows, err = r.DB.Query(query)
	}else{
		rows, err = r.DB.Query(query, year)
	}


	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}
	defer rows.Close()

	var invoices []model.Invoice
	for rows.Next() {
		var inv model.Invoice
		if err := rows.Scan(&inv.ID, &inv.Month, &inv.Year, &inv.Amount); err != nil {
			logger.LogToFile(fmt.Sprintf("error scanning invoice row: %v", err))
			return nil, err
		}
		invoices = append(invoices, inv)
	}

	return invoices, nil
}