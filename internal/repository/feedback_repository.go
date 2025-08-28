package repository

import (

	"fmt"
	"sync"
	"database/sql"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
)

type FeedbackRepository struct {
	mu sync.Mutex
	DB *sql.DB
}

type FeedbackRepositoryInterface interface {
	SaveFeedback(model.Feedback) error
	GetFeedbacksByID(string) ([]model.Feedback, error)
	GetAllFeedbacks() ([]model.Feedback, error)
}

func NewFeedbackRepository(db *sql.DB) *FeedbackRepository {
	return &FeedbackRepository{DB: db}
}

func (r *FeedbackRepository) SaveFeedback(feedback model.Feedback) error {
	feedback.ID = utils.GenerateUUID()

	query := `
		INSERT INTO feedbacks (id, resident_id, rating, content, flat_no)
		VALUES ($1, $2, $3, $4, $5)
	`
	r.mu.Lock()
	_, err := r.DB.Exec(query, feedback.ID, feedback.ResidentID, feedback.Rating, feedback.Content, feedback.Flat)
	r.mu.Unlock()

	if err != nil {logger.LogToFile(fmt.Sprintf("error: %v", err))}
	return err
}

func (r *FeedbackRepository) GetFeedbacksByID(residentID string) ([]model.Feedback, error) {
	query := `
		SELECT id, resident_id, rating, content, flat_no
		FROM feedbacks
		WHERE resident_id = $1
	`

	r.mu.Lock()
	rows, err := r.DB.Query(query, residentID)
	r.mu.Unlock()

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}

	defer rows.Close()

	var feedbacks []model.Feedback
	for rows.Next() {
		var f model.Feedback
		if err := rows.Scan(&f.ID, &f.ResidentID, &f.Rating, &f.Content, &f.Flat); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}
	return feedbacks, nil
}

func (r *FeedbackRepository) GetAllFeedbacks() ([]model.Feedback, error) {
	query := `
		SELECT id, resident_id, rating, content, flat_no
		FROM feedbacks
	`
	r.mu.Lock()
	rows, err := r.DB.Query(query)
	r.mu.Unlock()

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}
	defer rows.Close()

	var feedbacks []model.Feedback
	for rows.Next() {
		var f model.Feedback
		if err := rows.Scan(&f.ID, &f.ResidentID, &f.Rating, &f.Content, &f.Flat); err != nil {
			logger.LogToFile(fmt.Sprintf("error: %v", err))
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}
	return feedbacks, nil
}