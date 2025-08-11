package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/MananLed/majorProjectSMS/constants"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
)

type NoticeRepositoryInterface interface {
	SaveNotice(notice model.Notice) error
	GetAllNotices() ([]model.Notice, error)
	GetNoticesByMonthYear(month string, year string) ([]model.Notice, error)
	GetNoticesByYear(year string) ([]model.Notice, error)
}

type NoticeRepository struct {
	mu sync.Mutex
}

func (r *NoticeRepository) loadNotices() ([]model.Notice, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	fileData, err := os.ReadFile(string(constants.NoticeDataPath))
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Notice{}, nil
		}
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}

	var notices []model.Notice
	if err := json.Unmarshal(fileData, &notices); err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}

	return notices, nil
}

func (r *NoticeRepository) SaveNotice(notice model.Notice) error {
	notices, err := r.loadNotices()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	notice.ID = utils.GenerateUUID()
	notices = append(notices, notice)

	r.mu.Lock()
	defer r.mu.Unlock()
	data, err := json.MarshalIndent(notices, "", "  ")
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	return os.WriteFile(string(constants.NoticeDataPath), data, 0644)
}

func (r *NoticeRepository) GetAllNotices() ([]model.Notice, error) {
	return r.loadNotices()
}

func (r *NoticeRepository) GetNoticesByMonthYear(month string, year string) ([]model.Notice, error) {
	notices, err := r.loadNotices()
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}

	var noticesOfMonth []model.Notice

	for _, n := range notices {
		if n.Month == month && n.Year == year {
			noticesOfMonth = append(noticesOfMonth, n)
		}
	}

	return noticesOfMonth, nil
}

func (r *NoticeRepository) GetNoticesByYear(year string) ([]model.Notice, error) {
	notices, err := r.loadNotices()

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}

	var noticesYear []model.Notice

	for _, notice := range notices {
		if notice.Year == year {
			noticesYear = append(noticesYear, notice)
		}
	}
	return noticesYear, nil
}
