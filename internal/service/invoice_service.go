package service

import (
	"time"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/repository"
)

type InvoiceServiceInterface interface {
	GenerateInvoice(amount float64, month time.Month, year int) error
	GetInvoiceByMonthAndYear(month time.Month, year int) (*model.Invoice, error)
	GetInvoicesByYear(year int) ([]model.Invoice, error)
}

type InvoiceService struct {
	InvoiceRepo repository.InvoiceRepositoryInterface
}

func NewInvoiceService(repo repository.InvoiceRepositoryInterface) *InvoiceService {
	return &InvoiceService{InvoiceRepo: repo}
}

func (s *InvoiceService) GenerateInvoice(amount float64, month time.Month, year int) error {
	invoice := model.Invoice{
		Amount: amount,
		Month:  month,
		Year:   year,
	}
	return s.InvoiceRepo.SaveInvoice(invoice)
}

func (s *InvoiceService) GetInvoiceByMonthAndYear(month time.Month, year int) (*model.Invoice, error) {
	return s.InvoiceRepo.GetInvoiceByMonthAndYear(month, year)
}

func (s *InvoiceService) GetInvoicesByYear(year int) ([]model.Invoice, error) {
	return s.InvoiceRepo.GetInvoicesByYear(year)
}