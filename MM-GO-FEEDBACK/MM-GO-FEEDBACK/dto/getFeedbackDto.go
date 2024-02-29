package dto

import (
	MMErr "feedback/mmerror"

	"github.com/go-playground/validator"
)

type FeedbackSort struct {
	Column string `json:"column,omitempty"`
	Order  string `json:"order,omitempty"`
}

type FeedbackFilterObj struct {
	ProjectName      string `json:"projectName,omitempty"`
	OrganizationName string `json:"organizationName,omitempty"`
	Rating           int    `json:"rating,omitempty"`
	FromDate         string `json:"fromDate,omitempty"`
	EndDate          string `json:"endDate,omitempty"`
}

type Feedback struct {
	Limit               int               `json:"limit,omitempty"`
	Offset              int               `json:"offset,omitempty"`
	FeedbackSearchField string            `json:"feedbackSearchField,omitempty"`
	FeedbackSort        FeedbackSort      `json:"feedbackSort,omitempty"`
	FeedbackFilterObj   FeedbackFilterObj `json:"feedbackFilterObj,omitempty"`
}

type ProjectFeedback struct {
	Project_Name     string `bson:"project_name" json:"projectName"`
	UserName         string `bson:"user_name" json:"userName"`
	Organization     string `bson:"organization_name" json:"organization"`
	Submitted_On     string `bson:"created_on" json:"submittedOn"`
	Feedback_rating  int    `bson:"feedback_rating" json:"feedbackRating"`
	Feedback_comment string `bson:"feedback_comment" json:"feedbackComment"`
}

type DistinctValues struct {
	Feedback_rating   []int    `bson:"distinct_feedback_ratings" json:"feedbackRating"`
	ProjectName       []string `bson:"distinct_project_names" json:"projectNames"`
	Organization_name []string `bson:"distinct_organization_names" json:"organization_name"`
}

type ProjectNames struct {
	ProjectName []string `bson:"distinctValues" json:"projectNames"`
}

type OrganizationDetails struct {
	Organization_name string `bson:"organization_name" json:"organization_name"`
}

type FeedbackRating struct {
	Feedback_rating int `bson:"feedback_rating" json:"feedbackRating"`
}
type Ratings struct {
	Rating []int `bson:"distinctValues"`
}

type TotalCount struct {
	TotalCount int `bson:"total_count"`
}

type GetFeedBackResponse struct {
	FeedbackDetails []ProjectFeedback `json:"feedbackDetails"`
	ProjectNames    []string          `bson:"user_name" json:"projectName"`
	Orgnames        []string          `json:"organizationName"`
	TotalCount      int               `json:"totalCount"`
	Ratings         []int             `json:"ratings"`
}

type ErrResponse struct {
	StatusCode    int    `json:"StatusCode"`
	StatusMessage string `json:"StatusMessage"`
}

func Validate(p Feedback) *MMErr.AppError {
	validate := validator.New()

	if err := validate.Struct(p); err != nil {
		return MMErr.NewBadRequestError("Bad request")
	}
	return nil
}
