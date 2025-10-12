package repository

import (
	"fmt"
	"sync"

	"database/sql"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
	"github.com/google/uuid"
)

type FeedbackRepository struct {
	mu sync.Mutex
	DB *sql.DB
}

type FeedbackRepositoryInterface interface {
	SaveFeedback(model.Feedback) error
	GetFeedbacksByID(string) ([]model.Feedback, error)
	GetAllFeedbacks() ([]model.Feedback, error)
	IsFeedbackPresent(requestID uuid.UUID) (bool, error)
}

func NewFeedbackRepository(db *sql.DB) *FeedbackRepository {
	return &FeedbackRepository{DB: db}
}

func (r *FeedbackRepository) SaveFeedback(feedback model.Feedback) error {
	if feedback.ID == uuid.Nil {
		feedback.ID = utils.GenerateUUID()
	}

	var user model.User

	queryForUserDetails := `
		SELECT u.first_name, u.middle_name, u.last_name FROM
		users u WHERE u.id = $1
	`

	row := r.DB.QueryRow(queryForUserDetails, feedback.ResidentID)

	err := row.Scan(&user.FirstName, &user.MiddleName, &user.LastName)

	if err != nil{
		logger.LogToFile(fmt.Sprintf("error: %v", err))
	}

	var name string

	if user.MiddleName != "" {
		name = fmt.Sprintf("%s %s %s", user.FirstName, user.MiddleName, user.LastName)
	} else {
		name = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}

	feedback.ResidentName = name

	queryForUpdatingFeedbackStatus := `
		UPDATE service_requests set feedback_given = $1 WHERE request_id = $2
	`

	_ , err = r.DB.Exec(queryForUpdatingFeedbackStatus, true, feedback.RequestID)


	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return fmt.Errorf("failed to update feedback given status: %v", err)
	}

	queryForServiceDetails := `
		SELECT assigned_to, service_type, date, time_slot from service_requests
		WHERE request_id = $1
	`

	row = r.DB.QueryRow(queryForServiceDetails, feedback.RequestID)

	err = row.Scan(&feedback.AssignedTo, &feedback.ServiceType, &feedback.Date, &feedback.TimeSlot)

	if err != nil{
		logger.LogToFile(fmt.Sprintf("error: %v", err))
	}

	query := `
		INSERT INTO feedbacks (id, resident_id, rating, content, flat_no, username, request_id, assigned_to, service_type, date, time_slot)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = r.DB.Exec(query, feedback.ID, feedback.ResidentID, feedback.Rating, feedback.Content, feedback.Flat, feedback.ResidentName, feedback.RequestID, feedback.AssignedTo, feedback.ServiceType, feedback.Date, feedback.TimeSlot)


	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
	}
	return err
}

func (r *FeedbackRepository) GetFeedbacksByID(residentID string) ([]model.Feedback, error) {
	query := `
		SELECT id, resident_id, rating, content, flat_no, username, request_id, assigned_to, service_type, date, time_slot
		FROM feedbacks
		WHERE resident_id = $1
	`


	rows, err := r.DB.Query(query, residentID)


	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}

	defer rows.Close()

	var feedbacks []model.Feedback
	for rows.Next() {
		var f model.Feedback
		if err := rows.Scan(&f.ID, &f.ResidentID, &f.Rating, &f.Content, &f.Flat, &f.ResidentName, &f.AssignedTo, &f.ServiceType, &f.Date, &f.TimeSlot); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}
	return feedbacks, nil
}

func (r *FeedbackRepository) GetAllFeedbacks() ([]model.Feedback, error) {
	query := `
		SELECT id, resident_id, rating, content, flat_no, username, request_id, assigned_to, service_type, date, time_slot
		FROM feedbacks
	`

	rows, err := r.DB.Query(query)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}
	defer rows.Close()

	var feedbacks []model.Feedback
	for rows.Next() {
		var f model.Feedback
		if err := rows.Scan(&f.ID, &f.ResidentID, &f.Rating, &f.Content, &f.Flat, &f.ResidentName, &f.RequestID, &f.AssignedTo, &f.ServiceType, &f.Date, &f.TimeSlot); err != nil {
			logger.LogToFile(fmt.Sprintf("error: %v", err))
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}
	return feedbacks, nil
}

func (r *FeedbackRepository) IsFeedbackPresent(requestID uuid.UUID) (bool, error){

	var exists bool 

	query := `SELECT EXISTS (SELECT 1 FROM feedbacks WHERE request_id = $1);`

	err := r.DB.QueryRow(query, requestID).Scan(&exists)

	if(err != nil){
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return false, err 
	}

	return exists, nil
}