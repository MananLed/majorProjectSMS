package service

import (
	"errors"
	"testing"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/google/uuid"
)


type MockNoticeRepo struct {
	notices map[uuid.UUID]model.Notice
	saveErr error
	getErr  error
}

func (m *MockNoticeRepo) SaveNotice(notice model.Notice) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	if m.notices == nil {
		m.notices = make(map[uuid.UUID]model.Notice)
	}
	if notice.ID == uuid.Nil {
		notice.ID = uuid.New()
	}
	m.notices[notice.ID] = notice
	return nil
}

func (m *MockNoticeRepo) GetAllNotices() ([]model.Notice, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var res []model.Notice
	for _, n := range m.notices {
		res = append(res, n)
	}
	return res, nil
}

func (m *MockNoticeRepo) GetNoticesByMonthYear(month time.Month, year int) ([]model.Notice, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var res []model.Notice
	for _, n := range m.notices {
		if n.Month == month && n.Year == year {
			res = append(res, n)
		}
	}
	return res, nil
}

func (m *MockNoticeRepo) GetNoticesByYear(year int) ([]model.Notice, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var res []model.Notice
	for _, n := range m.notices {
		if n.Year == year {
			res = append(res, n)
		}
	}
	return res, nil
}


func TestIssueNotice(t *testing.T) {
	mockRepo := &MockNoticeRepo{notices: make(map[uuid.UUID]model.Notice)}
	service := NewNoticeService(mockRepo)

	err := service.IssueNotice("Water supply will be off", time.August, 2025)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if len(mockRepo.notices) != 1 {
		t.Errorf("expected 1 notice saved, got %d", len(mockRepo.notices))
	}
}

func TestGetAllNotices(t *testing.T) {
	mockRepo := &MockNoticeRepo{notices: make(map[uuid.UUID]model.Notice)}
	service := NewNoticeService(mockRepo)

	id := uuid.New()
	mockRepo.notices[id] = model.Notice{
		ID:         id,
		Content:    "Maintenance work",
		DateIssued: time.Now(),
		Month:      time.September,
		Year:       2025,
	}

	notices, err := service.GetNotices()
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if len(notices) != 1 {
		t.Errorf("expected 1 notice, got %d", len(notices))
	}
}

func TestGetNoticesByMonthYear(t *testing.T) {
	mockRepo := &MockNoticeRepo{notices: make(map[uuid.UUID]model.Notice)}
	service := NewNoticeService(mockRepo)

	id := uuid.New()
	mockRepo.notices[id] = model.Notice{
		ID:         id,
		Content:    "Festival celebration",
		DateIssued: time.Now(),
		Month:      time.October,
		Year:       2025,
	}

	notices, err := service.GetNoticesByMonthYear(time.October, 2025)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if len(notices) != 1 {
		t.Errorf("expected 1 notice, got %d", len(notices))
	}
}

func TestGetNoticesByYear(t *testing.T) {
	mockRepo := &MockNoticeRepo{notices: make(map[uuid.UUID]model.Notice)}
	service := NewNoticeService(mockRepo)

	id1 := uuid.New()
	id2 := uuid.New()
	mockRepo.notices[id1] = model.Notice{ID: id1, Content: "Meeting", Month: time.January, Year: 2025}
	mockRepo.notices[id2] = model.Notice{ID: id2, Content: "Holiday", Month: time.February, Year: 2024}

	notices, err := service.GetNoticesByYear(2025)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if len(notices) != 1 {
		t.Errorf("expected 1 notice for 2025, got %d", len(notices))
	}
}

func TestSaveNoticeError(t *testing.T) {
	mockRepo := &MockNoticeRepo{saveErr: errors.New("db error")}
	service := NewNoticeService(mockRepo)

	err := service.IssueNotice("Power cut notice", time.November, 2025)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetNoticeError(t *testing.T) {
	mockRepo := &MockNoticeRepo{getErr: errors.New("db error")}
	service := NewNoticeService(mockRepo)

	_, err := service.GetNotices()
	if err == nil {
		t.Error("expected error, got nil")
	}
	_, err = service.GetNoticesByMonthYear(time.December, 2025)
	if err == nil {
		t.Error("expected error, got nil")
	}
	_, err = service.GetNoticesByYear(2025)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
