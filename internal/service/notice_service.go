//go:generate mockgen -source=notice_service.go -destination=../mocks/notice_mock_service.go -package=mocks
package service

import (
	"time"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/repository"
)

type NoticeServiceInterface interface {
	IssueNotice(content string, month time.Month, year int) error
	GetNotices() ([]model.Notice, error)
	GetNoticesByMonthYear(month time.Month, year int) ([]model.Notice, error)
	GetNoticesByYear(year int) ([]model.Notice, error)
}

type NoticeService struct {
	NoticeRepo repository.NoticeRepositoryInterface
}

func NewNoticeService(repo repository.NoticeRepositoryInterface) *NoticeService {
	return &NoticeService{NoticeRepo: repo}
}

func (s *NoticeService) IssueNotice(content string, month time.Month, year int) error {
	notice := model.Notice{
		DateIssued: time.Now(),
		Content:    content,
		Month:      month,
		Year:       year,
	}
	return s.NoticeRepo.SaveNotice(notice)
}

func (s *NoticeService) GetNotices() ([]model.Notice, error) {
	return s.NoticeRepo.GetAllNotices()
}

func (s *NoticeService) GetNoticesByMonthYear(month time.Month, year int) ([]model.Notice, error) {
	return s.NoticeRepo.GetNoticesByMonthYear(month, year)
}

func (s *NoticeService) GetNoticesByYear(year int) ([]model.Notice, error) {
	return s.NoticeRepo.GetNoticesByYear(year)
}