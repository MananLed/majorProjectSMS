package service

import (
	"errors"
	"testing"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/google/uuid"
)

type MockFeedbackRepo struct {
	feedbacks map[uuid.UUID]model.Feedback
	saveErr   error
	getErr    error
}

func (m *MockFeedbackRepo) SaveFeedback(f model.Feedback) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	if m.feedbacks == nil {
		m.feedbacks = make(map[uuid.UUID]model.Feedback)
	}
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	m.feedbacks[f.ID] = f
	return nil
}

func (m *MockFeedbackRepo) GetFeedbacksByID(residentID string) ([]model.Feedback, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var results []model.Feedback
	for _, f := range m.feedbacks {
		if f.ResidentID == residentID {
			results = append(results, f)
		}
	}
	return results, nil
}

func (m *MockFeedbackRepo) GetAllFeedbacks() ([]model.Feedback, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var all []model.Feedback
	for _, f := range m.feedbacks {
		all = append(all, f)
	}
	return all, nil
}

func TestIssueFeedback_Success(t *testing.T) {
	mockRepo := &MockFeedbackRepo{feedbacks: make(map[uuid.UUID]model.Feedback)}
	service := NewFeedbackService(mockRepo)

	err := service.IssueFeedback("Great service", "resident1", "101", 5)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(mockRepo.feedbacks) != 1 {
		t.Errorf("expected 1 feedback saved, got %d", len(mockRepo.feedbacks))
	}
}

func TestIssueFeedback_Failure(t *testing.T) {
	mockRepo := &MockFeedbackRepo{saveErr: errors.New("db error")}
	service := NewFeedbackService(mockRepo)

	err := service.IssueFeedback("Bad service", "resident2", "202", 2)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetFeedbackByID(t *testing.T) {
	mockRepo := &MockFeedbackRepo{feedbacks: make(map[uuid.UUID]model.Feedback)}
	service := NewFeedbackService(mockRepo)

	id1 := uuid.New()
	id2 := uuid.New()

	mockRepo.feedbacks[id1] = model.Feedback{ID: id1, ResidentID: "resident1", Rating: 4, Content: "Good", Flat: "101"}
	mockRepo.feedbacks[id2] = model.Feedback{ID: id2, ResidentID: "resident2", Rating: 2, Content: "Bad", Flat: "202"}

	results, err := service.GetFeedbackByID("resident1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(results) != 1 || results[0].ResidentID != "resident1" {
		t.Errorf("expected 1 feedback for resident1, got %+v", results)
	}
}

func TestGetFeedbacks(t *testing.T) {
	mockRepo := &MockFeedbackRepo{feedbacks: make(map[uuid.UUID]model.Feedback)}
	service := NewFeedbackService(mockRepo)

	id := uuid.New()
	mockRepo.feedbacks[id] = model.Feedback{ID: id, ResidentID: "resident1", Rating: 5, Content: "Awesome", Flat: "A-101"}

	results, err := service.GetFeedbacks()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 feedback, got %d", len(results))
	}
}

func TestGetFeedbacks_Error(t *testing.T) {
	mockRepo := &MockFeedbackRepo{getErr: errors.New("db error")}
	service := NewFeedbackService(mockRepo)

	_, err := service.GetFeedbacks()
	if err == nil {
		t.Error("expected error, got nil")
	}
}
