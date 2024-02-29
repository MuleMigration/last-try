package Service

import (
	"errors"
	repository "feedback/Repository"
	"feedback/dto"
	MMLogger "feedback/logger"
)

var log MMLogger.Logger

type GetFeedbackService struct {
	repo repository.GetFeedBackRepoI
}

type GetFeedBackServiceI interface {
	GetAllFeedbacks(request dto.Feedback) (*dto.GetFeedBackResponse, error)
}

func NewGetFeedbackService(repo repository.GetFeedBackRepoI) GetFeedbackService {
	return GetFeedbackService{repo: repo}
}

func (s GetFeedbackService) GetAllFeedbacks(request dto.Feedback) (*dto.GetFeedBackResponse, error) {

	FeedbackData, err := s.repo.FetchFeedbacks(request)

	if err != nil {

		return nil, errors.New(err.Message)
	}

	if len(FeedbackData.FeedbackDetails) == 0 {
		return nil, errors.New("No Content Found")
	}

	return FeedbackData, nil
}
