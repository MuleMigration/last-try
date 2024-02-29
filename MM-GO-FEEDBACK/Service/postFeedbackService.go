package service

import (
	repository "postFeedback/Repository"
	"postFeedback/dto"
	MMErr "postFeedback/mmerror"
)

type PostFeedbackService struct {
	repo repository.PostFeedbackRepoI
}

type PostFeedBackServiceI interface {
	PostFeedback(request dto.Request) (*dto.Response, *MMErr.AppError)
}

func NewFeedbackService(repo repository.PostFeedbackRepoI) PostFeedbackService {
	return PostFeedbackService{repo: repo}
}

func (s *PostFeedbackService) PostFeedback(request dto.Request) (*dto.Response, *MMErr.AppError) {	

	feedback, err := s.repo.InsertFeedback(request)
	if err != nil {
		return nil, MMErr.NewUnexpectedError("Cursor Error")
	}

	return feedback, nil
}
