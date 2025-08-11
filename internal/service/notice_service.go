package service

import (
	"time"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/repository"
)

type NoticeService struct {
	NoticeRepo repository.NoticeRepositoryInterface
}

func NewNoticeService(repo repository.NoticeRepositoryInterface) *NoticeService {
	return &NoticeService{NoticeRepo: repo}
}

func (s *NoticeService) IssueNotice(content string, month string, year string) error {
	notice := model.Notice{
		DateIssued: time.Now().Format("2006-01-02 15:04:05"),
		Content:    content,
		Month: month,
		Year: year,
	}
	return s.NoticeRepo.SaveNotice(notice)
}

func (s *NoticeService) GetNotices() ([]model.Notice, error) {
	return s.NoticeRepo.GetAllNotices()
}

func (s *NoticeService) GetNoticesByMonthYear(month string, year string) ([]model.Notice, error) {	
	return s.NoticeRepo.GetNoticesByMonthYear(month, year)
}

func (s *NoticeService) GetNoticesByYear(year string) ([]model.Notice, error){
	return s.NoticeRepo.GetNoticesByYear(year)
}