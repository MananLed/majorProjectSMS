package service

import (
	"errors"
	"testing"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/google/uuid"
)

type MockInvoiceRepo struct {
	invoices map[uuid.UUID]model.Invoice
	saveErr  error
	getErr   error
}

func (m *MockInvoiceRepo) SaveInvoice(inv model.Invoice) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	if m.invoices == nil {
		m.invoices = make(map[uuid.UUID]model.Invoice)
	}

	if inv.ID == uuid.Nil {
		inv.ID = uuid.New()
	}
	m.invoices[inv.ID] = inv
	return nil
}

func (m *MockInvoiceRepo) GetInvoiceByMonthAndYear(month time.Month, year int) (*model.Invoice, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	for _, inv := range m.invoices {
		if inv.Month == month && inv.Year == year {
			return &inv, nil
		}
	}
	return nil, errors.New("invoice not found")
}

func (m *MockInvoiceRepo) GetInvoicesByYear(year int) ([]model.Invoice, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var res []model.Invoice
	for _, inv := range m.invoices {
		if inv.Year == year {
			res = append(res, inv)
		}
	}
	return res, nil
}


func TestGenerateInvoice(t *testing.T) {
	mockRepo := &MockInvoiceRepo{invoices: make(map[uuid.UUID]model.Invoice)}
	service := NewInvoiceService(mockRepo)

	err := service.GenerateInvoice(1500.75, time.August, 2025)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if len(mockRepo.invoices) != 1 {
		t.Errorf("expected 1 invoice saved, got %d", len(mockRepo.invoices))
	}
}

func TestGetInvoiceByMonthAndYear(t *testing.T) {
	mockRepo := &MockInvoiceRepo{invoices: make(map[uuid.UUID]model.Invoice)}
	service := NewInvoiceService(mockRepo)

	id := uuid.New()
	inv := model.Invoice{ID: uuid.New(), Amount: 2000, Month: time.September, Year: 2025}
	mockRepo.invoices[id] = inv

	got, err := service.GetInvoiceByMonthAndYear(time.September, 2025)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if got.Amount != inv.Amount {
		t.Errorf("expected amount %v, got %v", inv.Amount, got.Amount)
	}
}

func TestGetInvoiceByMonthAndYear_NotFound(t *testing.T) {
	mockRepo := &MockInvoiceRepo{invoices: make(map[uuid.UUID]model.Invoice)}
	service := NewInvoiceService(mockRepo)

	_, err := service.GetInvoiceByMonthAndYear(time.October, 2025)
	if err == nil {
		t.Error("expected error for not found, got nil")
	}
}

func TestGetInvoicesByYear(t *testing.T) {
	mockRepo := &MockInvoiceRepo{invoices: make(map[uuid.UUID]model.Invoice)}
	service := NewInvoiceService(mockRepo)

	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	mockRepo.invoices[id1] = model.Invoice{ID: id1, Amount: 1000, Month: time.January, Year: 2025}
	mockRepo.invoices[id2] = model.Invoice{ID: id2, Amount: 2000, Month: time.February, Year: 2025}
	mockRepo.invoices[id3] = model.Invoice{ID: id3, Amount: 3000, Month: time.March, Year: 2024}

	invoices, err := service.GetInvoicesByYear(2025)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if len(invoices) != 2 {
		t.Errorf("expected 2 invoices, got %d", len(invoices))
	}
}

func TestSaveInvoiceError(t *testing.T) {
	mockRepo := &MockInvoiceRepo{saveErr: errors.New("db error")}
	service := NewInvoiceService(mockRepo)

	err := service.GenerateInvoice(999, time.August, 2025)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetInvoiceError(t *testing.T) {
	mockRepo := &MockInvoiceRepo{getErr: errors.New("db error")}
	service := NewInvoiceService(mockRepo)

	_, err := service.GetInvoiceByMonthAndYear(time.August, 2025)
	if err == nil {
		t.Error("expected error, got nil")
	}

	_, err = service.GetInvoicesByYear(2025)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
