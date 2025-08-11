package service

import (
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/repository"
)

type InvoiceService struct {
	InvoiceRepo repository.InvoiceRepositoryInterface
}

func NewInvoiceService(repo repository.InvoiceRepositoryInterface) *InvoiceService{
	return &InvoiceService{InvoiceRepo: repo} 
}

func (s *InvoiceService) GenerateInvoice(amount float64, month string, year string) error{
	invoice := model.Invoice{
		Amount: amount,
		Month: month,
		Year: year,
	}

	return s.InvoiceRepo.SaveInvoice(invoice)
}

func (s *InvoiceService) GetInvoiceByMonthAndYear(month string, year string) (*model.Invoice, error){
	return s.InvoiceRepo.GetInvoiceByMonthAndYear(month, year)
}

func (s *InvoiceService) GetInvoicesByYear(year string) ([]model.Invoice, error){
	return s.InvoiceRepo.GetInvoicesByYear(year)
}