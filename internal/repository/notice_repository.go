package repository

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
)

type NoticeRepositoryInterface interface {
	SaveNotice(notice model.Notice) error
	GetAllNotices() ([]model.Notice, error)
	GetNoticesByMonthYear(month time.Month, year int) ([]model.Notice, error)
	GetNoticesByYear(year int) ([]model.Notice, error)
}

type NoticeRepository struct {
	DB *sql.DB
	mu sync.Mutex
}

func NewNoticeRepository(db *sql.DB) *NoticeRepository {
	return &NoticeRepository{DB: db}
}

func (r *NoticeRepository) SaveNotice(notice model.Notice) error {

	notice.ID = utils.GenerateUUID()

	query := `
		INSERT INTO notices (id, date_issued, content, month, year)
		VALUES ($1, $2, $3, $4, $5)
	`
	r.mu.Lock()
	_, err := r.DB.Exec(query, notice.ID, notice.DateIssued, notice.Content, notice.Month, notice.Year)
	r.mu.Unlock()

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}
	return nil
}

func (r *NoticeRepository) GetAllNotices() ([]model.Notice, error) {
	query := `
	SELECT id, date_issued, content, month, year 
	FROM notices 
	ORDER BY date_issued DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}
	defer rows.Close()

	var notices []model.Notice
	for rows.Next() {
		var n model.Notice
		if err := rows.Scan(&n.ID, &n.DateIssued, &n.Content, &n.Month, &n.Year); err != nil {
			logger.LogToFile(fmt.Sprintf("error: %v", err))
			return nil, err
		}
		notices = append(notices, n)
	}

	return notices, nil
}

func (r *NoticeRepository) GetNoticesByMonthYear(month time.Month, year int) ([]model.Notice, error) {
	query := `
	SELECT id, date_issued, content, month, year 
	FROM notices 
	WHERE month = $1 AND year = $2 
	ORDER BY date_issued DESC
	`
	r.mu.Lock()
	rows, err := r.DB.Query(query, month, year)
	r.mu.Unlock()

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}
	defer rows.Close()

	var notices []model.Notice
	for rows.Next() {
		var n model.Notice
		if err := rows.Scan(&n.ID, &n.DateIssued, &n.Content, &n.Month, &n.Year); err != nil {
			logger.LogToFile(fmt.Sprintf("error: %v", err))
			return nil, err
		}
		notices = append(notices, n)
	}

	return notices, nil
}

func (r *NoticeRepository) GetNoticesByYear(year int) ([]model.Notice, error) {
	query := `
	SELECT id, date_issued, content, month, year 
	FROM notices WHERE year = $1 
	ORDER BY date_issued DESC
	`
	r.mu.Lock()
	rows, err := r.DB.Query(query, year)
	r.mu.Unlock()

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}
	defer rows.Close()

	var notices []model.Notice
	for rows.Next() {
		var n model.Notice
		if err := rows.Scan(&n.ID, &n.DateIssued, &n.Content, &n.Month, &n.Year); err != nil {
			logger.LogToFile(fmt.Sprintf("error: %v", err))
			return nil, err
		}
		notices = append(notices, n)
	}

	return notices, nil
}
