package Repository

import (
	"context"
	"feedback/dto"
	MMErr "feedback/mmerror"
	MongoConnect "feedback/mongoconnect"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostFeedbackRepository struct {
	client MongoConnect.MongoDBInterface
}

type PostFeedbackRepoI interface {
	InsertFeedback(request dto.Request) (*dto.PostResponse, *MMErr.AppError)
}

func NewPostFeedbackRepository(db MongoConnect.MongoDBInterface) PostFeedbackRepository {
	log.SetLevel(2)
	return PostFeedbackRepository{client: db}
}

func (r *PostFeedbackRepository) InsertFeedback(request dto.Request) (*dto.PostResponse, *MMErr.AppError) {

	AppDb, AppDbErr := r.client.GetAppClient()
	if AppDbErr != nil {
		log.Error("Error getting Application Client", AppDbErr)
		return nil, AppDbErr
	}

	feedbackCollection := AppDb.Collection("ProjectFeedback")
	userDetailCollection := AppDb.Collection("UserDetail")
	organizationDetail := AppDb.Collection("OrganizationDetail")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user_id, err := primitive.ObjectIDFromHex(request.UserId)
	if err != nil {
		log.Error("Error", err)
		return nil, MMErr.NewUnexpectedError("failed to decode the user id from hex")
	}

	//PC_NO_1.25 - //PC_NO_1.30

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "$eq", Value: user_id}}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "organization_id", Value: 1},
		}}},
	}

	cursor, err := userDetailCollection.Aggregate(context.Background(), pipeline)

	if err != nil {
		log.Error("Error", err)
		return nil, MMErr.NewUnexpectedError("Failed in fetching the organization id from the database")
	}

	defer cursor.Close(context.Background())

	var user dto.PostUserDetails
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		cursor.Decode(&user)
	}

	if err := cursor.Err(); err != nil {
		log.Error("Error", err)
		return nil, MMErr.NewUnexpectedError("Failed to decode the user details")
	}

	//PC_NO_1.31 - //PC_NO_1.34
	pipeline = mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "$eq", Value: user.Organization_ID}}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "organization_id", Value: 1},
			{Key: "organization_name", Value: 1},
		}}},
	}

	cursor, err = organizationDetail.Aggregate(context.Background(), pipeline)

	if err != nil {
		log.Error("Error", err)
		return nil, MMErr.NewUnexpectedError("Failed to fetch the organization details")
	}

	var org dto.OrgDetails
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		cursor.Decode(&org)
	}

	if err := cursor.Err(); err != nil {
		return nil, MMErr.NewUnexpectedError("Failed to decode the organization details")
	}

	project_id, err2 := primitive.ObjectIDFromHex(request.ProjectId)
	if err2 != nil {
		return nil, MMErr.NewUnexpectedError("Failed to decode the project id from hex")
	}

	//PC_NO_1.35- PC_NO_1.38

	respone := dto.PostProjectFeedback{
		ProjectID:        project_id,
		ProjectName:      request.ProjectName,
		UserID:           user_id,
		OrganizationName: org.OrganizationName,
		FeedbackRating:   request.FeedbackRating,
		FeedbackComment:  request.FeedbackComment,
		CreatedBy:        user_id,
		CreatedOn:        time.Now(),
		ModifiedBy:       user_id,
		ModifiedOn:       time.Now(),
		IsActive:         1,
	}

	_, Updateerr := feedbackCollection.InsertOne(context.Background(), respone)

	if Updateerr != nil {
		return nil, MMErr.NewUnexpectedError("Failed to insert the post feedback into the database")
	}

	return &dto.PostResponse{StatusCode: http.StatusCreated, StatusMessage: "Sucessfully Inserted the Feedback."}, nil
}
