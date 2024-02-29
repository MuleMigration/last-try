package Repository

import (
	"context"
	"feedback/dto"
	MMLogger "feedback/logger"
	MongoConnect "feedback/mongoconnect"
	"fmt"
	"time"

	MMErr "feedback/mmerror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var log MMLogger.Logger

type GetFeedbackRepository struct {
	client MongoConnect.MongoDBInterface
}

type GetFeedBackRepoI interface {
	FetchFeedbacks(request dto.Feedback) (*dto.GetFeedBackResponse, *MMErr.AppError)
}

func NewGetFeedbackRepository(db MongoConnect.MongoDBInterface) GetFeedbackRepository {
	return GetFeedbackRepository{client: db}
}

func (r GetFeedbackRepository) FetchFeedbacks(request dto.Feedback) (*dto.GetFeedBackResponse, *MMErr.AppError) {

	feedbacks := make([]dto.ProjectFeedback, 0)
	AppDb, AppDbErr := r.client.GetAppClient()
	if AppDbErr != nil {
		log.Error("Error getting Application Client", AppDbErr)
		return nil, AppDbErr
	}
	feedbackCollection := AppDb.Collection("ProjectFeedback")
	// ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	var pipeline bson.A

	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "UserDetail"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "User_Detail"},
		}},
	}
	pipeline = append(pipeline, lookupStage)

	unwindStage := bson.D{{Key: "$unwind", Value: "$User_Detail"}}
	pipeline = append(pipeline, unwindStage)

	if request.FeedbackFilterObj.ProjectName != "" {
		projectMatchStage := bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "project_name", Value: bson.D{{Key: "$eq", Value: request.FeedbackFilterObj.ProjectName}}},
			}},
		}
		pipeline = append(pipeline, projectMatchStage)
	}

	if request.FeedbackFilterObj.OrganizationName != "" {
		organizationMatchStage := bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "organization_name", Value: bson.D{{Key: "$eq", Value: request.FeedbackFilterObj.OrganizationName}}},
			}},
		}
		pipeline = append(pipeline, organizationMatchStage)
	}

	if request.FeedbackFilterObj.Rating != 0 {
		ratingMatchStage := bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "feedback_rating", Value: bson.D{{Key: "$eq", Value: request.FeedbackFilterObj.Rating}}},
			}},
		}
		pipeline = append(pipeline, ratingMatchStage)
	}

	if request.FeedbackFilterObj.FromDate != "" {
		dateStr := request.FeedbackFilterObj.FromDate
		layout := "02/01/2006"
		fromDate, err := time.Parse(layout, dateStr)
		if err != nil {
			log.Error("error in date conversion:", err)
			return nil, MMErr.NewUnexpectedError("The Date conversion error")
		}
		bsonDateTime := primitive.NewDateTimeFromTime(fromDate)
		projectMatchStage := bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "created_on", Value: bson.D{{Key: "$gte", Value: bsonDateTime}}},
			}},
		}
		pipeline = append(pipeline, projectMatchStage)
	}

	if request.FeedbackFilterObj.EndDate != "" {
		dateStr := request.FeedbackFilterObj.EndDate
		layout := "02/01/2006"
		fromDate, err := time.Parse(layout, dateStr)
		if err != nil {
			log.Error("error in date conversion:", err)
			return nil, MMErr.NewUnexpectedError("The Date conversion error")
		}
		bsonDateTime := primitive.NewDateTimeFromTime(fromDate)
		projectMatchStage := bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "created_on", Value: bson.D{{Key: "$lte", Value: bsonDateTime}}},
			}},
		}
		pipeline = append(pipeline, projectMatchStage)
	}

	sortValue := -1 //default descending order
	if request.FeedbackSort.Order == "asc" {
		sortValue = 1
	}

	if request.FeedbackSort.Column != "" {
		sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: request.FeedbackSort.Column, Value: sortValue}}}}
		pipeline = append(pipeline, sortStage)
	}

	offset := 0
	if request.Offset != 0 {
		offset = request.Offset
	}
	skipStage := bson.D{{Key: "$skip", Value: offset}}
	pipeline = append(pipeline, skipStage)

	limit := 10

	if request.Limit != 0 {
		limit = request.Limit
	}

	limitStage := bson.D{{Key: "$limit", Value: limit}}
	pipeline = append(pipeline, limitStage)

	//to search the user name
	addFieldStage := bson.D{
		{Key: "$addFields", Value: bson.D{
			{Key: "fullName", Value: bson.D{
				{Key: "$concat", Value: bson.A{
					"$User_Detail.first_name",
					" ",
					"$User_Detail.last_name",
				}},
			}},
		}},
	}

	pipeline = append(pipeline, addFieldStage)

	if request.FeedbackSearchField != "" {
		searchStage := bson.D{

			{Key: "$match", Value: bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "organization_name", Value: bson.D{{Key: "$regex", Value: request.FeedbackSearchField}, {Key: "$options", Value: "i"}}}},
					bson.D{{Key: "feedback_comment", Value: bson.D{{Key: "$regex", Value: request.FeedbackSearchField}, {Key: "$options", Value: "i"}}}},
					bson.D{{Key: "last_modified_by", Value: bson.D{{Key: "$regex", Value: request.FeedbackSearchField}, {Key: "$options", Value: "i"}}}},
					bson.D{{Key: "project_name", Value: bson.D{{Key: "$regex", Value: request.FeedbackSearchField}, {Key: "$options", Value: "i"}}}},
					bson.D{{Key: "User_Detail.first_name", Value: bson.D{{Key: "$regex", Value: request.FeedbackSearchField}, {Key: "$options", Value: "i"}}}},
					bson.D{{Key: "User_Detail.last_name", Value: bson.D{{Key: "$regex", Value: request.FeedbackSearchField}, {Key: "$options", Value: "i"}}}},
					bson.D{{Key: "fullName", Value: bson.D{{Key: "$regex", Value: request.FeedbackSearchField}, {Key: "$options", Value: "i"}}}}, // Include user name in the search
					// Include user name in the search

					// bson.D{{Key: "$unwind", Value: "$User_Detail"}},
					// bson.D{
					// 	{Key: "$expr", Value: bson.D{
					// 		{Key: "$regexMatch", Value: bson.D{
					// 			{Key: "input", Value: bson.D{
					// 				{Key: "$concat", Value: bson.A{"$User_Detail.first_name", " ", "$User_Details.last_name"}},
					// 			}},
					// 			{Key: "regex", Value: request.FeedbackSearchField},
					// 			{Key: "options", Value: "i"},
					// 		}},
					// 	}},
					// // },
					// bson.D{
					// 	{Key: "$project",
					// 		Value: bson.D{
					// 			{Key: "name",
					// 				Value: bson.D{
					// 					{Key: "$concat",
					// 						Value: bson.A{
					// 							"$User_Detail.first_name",
					// 							" ",
					// 							"$User_Detail.last_name",
					// 						},
					// 					},
					// 				},
					// 			},
					// 		},
					// 	},
					// },
					// bson.D{{Key: "$match", Value: bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: request.FeedbackSearchField}}}}}},
				}},
			}},
		}
		pipeline = append(pipeline, searchStage)
	}

	// 			Value: bson.D{
	// 				{Key: "$dateToString",
	// 					Value: bson.D{
	// 						{Key: "format", Value: "%d/%m/%Y %H:%M"},
	// 						{Key: "date", Value: "$created_on"},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	}},
	// }

	projectStage := bson.D{
		{
			Key: "$project", Value: bson.D{
				{Key: "total_count", Value: bson.D{
					{Key: "$sum", Value: 1},
				}},
				{Key: "user_name", Value: bson.D{
					{Key: "$concat", Value: bson.A{"$User_Detail.first_name", " ", "$User_Detail.last_name"}},
				}},
				{Key: "feedback_rating", Value: 1},
				{Key: "feedback_comment", Value: 1},
				{Key: "organization_name", Value: 1},
				{Key: "project_name", Value: 1},
				{Key: "_id", Value: 0},
				{Key: "created_on", Value: bson.D{
					{Key: "$dateToString", Value: bson.D{
						{Key: "format", Value: "%d/%m/%Y %H:%M"},
						{Key: "date", Value: "$created_on"},
					}},
				}},
			},
		},
	}

	pipeline = append(pipeline, projectStage)

	cursor, err := feedbackCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Info("The Fetching Feedback error:", err)
		return nil, MMErr.NewUnexpectedError("Error while executing the query")
	}
	log.Info("Cursor for the feeback: ", cursor)
	defer cursor.Close(context.Background())

	for cursor.Next(context.TODO()) {
		var feedback dto.ProjectFeedback
		cursor.Decode(&feedback)
		feedbacks = append(feedbacks, feedback)
		log.Info("decoded feedback:", feedbacks)
	}

	log.Info("decode feedback: ", feedbacks)
	if err := cursor.Err(); err != nil {
		return nil, MMErr.NewUnexpectedError("Decode feedback error")
	}

	//

	if len(feedbacks) == 0 {
		return nil, MMErr.NewNoContentError("No feedback found")
	}

	distinctValues := bson.A{
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: primitive.Null{}},
					{Key: "distinct_feedback_ratings", Value: bson.D{{Key: "$addToSet", Value: "$feedback_rating"}}},
					{Key: "distinct_project_names", Value: bson.D{{Key: "$addToSet", Value: "$project_name"}}},
					{Key: "distinct_organization_names", Value: bson.D{{Key: "$addToSet", Value: "$organization_name"}}},
				},
			},
		},
	}

	fbCollection := AppDb.Collection("ProjectFeedback")

	DistinctValuesCursor, err := fbCollection.Aggregate(context.Background(), distinctValues)
	if err != nil {
		log.Error("Getting the distinctValues error")
		return nil, MMErr.NewUnexpectedError("Getting distinctValues error")
	}

	defer DistinctValuesCursor.Close(context.Background())

	//
	var DistinctValues dto.DistinctValues

	if DistinctValuesCursor.Next(context.Background()) {

		if err := DistinctValuesCursor.Decode(&DistinctValues); err != nil {
			// Handle error
			fmt.Println("Error decoding distinct values:", err.Error())
		}
	}

	//geting total count

	totalFB := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "is_active", Value: 1}}}},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: primitive.Null{}},
					{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
				},
			},
		},
	}
	totalFBcursor, err := fbCollection.Aggregate(context.Background(), totalFB)

	if err != nil {
		return nil, MMErr.NewUnexpectedError("Error while Getting the total count")
	}
	defer totalFBcursor.Close(context.Background())

	var count dto.TotalCount

	if totalFBcursor.Next(context.Background()) {
		if err := totalFBcursor.Decode(&count); err != nil {
			log.Error("Failed to decode the total count", err)
			return nil, MMErr.NewUnexpectedError("Error while decoding  the total count")
		}
	}

	// var orgResult []string
	// for _, organization := range organizations {
	// 	orgResult = append(orgResult, organization.Organization_name)
	// }

	// var ratingResult []int
	// for _, r := range ratings {
	// 	ratingResult = append(ratingResult, r.Rating)
	// }

	return &dto.GetFeedBackResponse{TotalCount: count.TotalCount, FeedbackDetails: feedbacks, ProjectNames: DistinctValues.ProjectName, Orgnames: DistinctValues.Organization_name, Ratings: DistinctValues.Feedback_rating}, nil
}
