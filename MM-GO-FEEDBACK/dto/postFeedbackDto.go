package dto

import (
	MMErr "postFeedback/mmerror"
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectFeedback struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID        primitive.ObjectID `bson:"project_id" json:"project_id" `
	ProjectName      string             `bson:"project_name" json:"project_name"`
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`
	OrganizationName string             `bson:"organization_name" json:"organization_name"`
	FeedbackRating   int                `bson:"feedback_rating" json:"feedback_rating"`
	FeedbackComment  string             `bson:"feedback_comment" json:"feedback_comment"`
	CreatedBy        primitive.ObjectID `bson:"created_by" json:"created_by"`
	CreatedOn        time.Time          `bson:"created_on" json:"created_on"`
	ModifiedBy       primitive.ObjectID `bson:"last_modified_by" json:"modified_by"`
	ModifiedOn       time.Time          `bson:"last_modified_on" json:"modified_on"`
	IsActive         int                `bson:"is_active" json:"is_active"`
}

type Request struct {
	ProjectId       string `bson:"project_id" json:"project_id" validate:"required"`
	ProjectName     string `bson:"project_name" json:"project_name" validate:"required"`
	UserId          string `bson:"user_id" json:"user_id" validate:"required"`
	FeedbackRating  int    `bson:"feedback_rating" json:"feedback Rating" validate:"oneof=1 2 3 4 5"`
	FeedbackComment string `bson:"feedback_comment" json:"feedback Comment" validate:"required"`
}

type Response struct {
	StatusCode    string `json:"StatusCode"`
	StatusMessage string `json:"StatusMessage"`
}

type UserDetails struct {
	Organization_ID primitive.ObjectID `bson:"organization_id" json:"organization_id"`
}

type OrgDetails struct {
	OrganizationName string `bson:"organization_name" json:"organization_name"`
}

func Validate(p Request) *MMErr.AppError {
	validate := validator.New()

	if err := validate.Struct(p); err != nil {
		return MMErr.NewBadRequestError("Bad request")
	}
	return nil
}
