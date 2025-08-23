package service

import (
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/repository"
)

type FeedbackServiceInterface interface{
	IssueFeedback(content string, residentID string,flat string, rating int32) error
	GetFeedbacks() ([]model.Feedback, error)
	GetFeedbackByID(id string) ([]model.Feedback, error)
}

type FeedbackService struct {
	FeedbackRepo repository.FeedbackRepositoryInterface
}

func NewFeedbackService(repo repository.FeedbackRepositoryInterface) *FeedbackService {
	return &FeedbackService{FeedbackRepo: repo}
}

func (s *FeedbackService) IssueFeedback(content string, residentID string,flat string, rating int32) error {
	feedback := model.Feedback{
		ResidentID: residentID,
		Rating:     rating,
		Content:    content,
		Flat: flat,
	}

	return s.FeedbackRepo.SaveFeedback(feedback)
}

func (s *FeedbackService) GetFeedbacks() ([]model.Feedback, error) {
	return s.FeedbackRepo.GetAllFeedbacks()
}

func (s *FeedbackService) GetFeedbackByID(id string) ([]model.Feedback, error){
	return s.FeedbackRepo.GetFeedbacksByID(id)
}