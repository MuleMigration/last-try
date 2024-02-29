package Service

import (
	"errors"
	repository "feedback/Repository"
	"feedback/dto"
)

type PostFeedbackService struct {
	repo repository.PostFeedbackRepoI
}

type PostFeedBackServiceI interface {
	PostFeedback(request dto.Request) (*dto.PostResponse, error)
}

func NewPostFeedbackService(repo repository.PostFeedbackRepoI) PostFeedbackService {

	log.SetLevel(2)
	return PostFeedbackService{repo: repo}
}

// PC_NO_1.19 - PC_NO_1.20
func (s *PostFeedbackService) PostFeedback(request dto.Request) (*dto.PostResponse, error) {

	postResponse, err := s.repo.InsertFeedback(request)
	//PC_NO_1.45 - PC_NO_1.47
	if err != nil {
		log.Error("Error getting repo response " + err.Message)
		return nil, errors.New(err.Message)
	}

	return postResponse, nil
}
