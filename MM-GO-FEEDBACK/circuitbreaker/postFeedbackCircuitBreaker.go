// Sequence to Pseudocode	01/31/2024 02:20 PM
package circuitbreaker

import (
	"context"
	"time"

	MMLogger "postFeedback/logger"
	MMErr "postFeedback/mmerror"
	"postFeedback/mongoConnect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CircuitBreakerDetail struct {
	ID             string    `bson:"_id,omitempty"`
	APIName        string    `bson:"api_name"`
	CurrentState   string    `bson:"current_state"`
	FailureCount   int       `bson:"failure_count"`
	SuccessCount   int       `bson:"success_count"`
	NextAttempt    time.Time `bson:"next_attempt"`
	StateTimeStamp time.Time `bson:"state_time_stamp"`
	CreatedBy      string    `bson:"created_by"`
	CreatedOn      time.Time `bson:"created_on"`
	LastModifiedBy string    `bson:"last_modified_by"`
	LastModifiedOn time.Time `bson:"last_modified_on"`
	IsActive       int       `bson:"is_active"`
}

var logger *MMLogger.Logger
var detail CircuitBreakerDetail

type CircuitBreaker struct {
	client MongoConnect.MongoDBInterface
}

const (
	Closed           = "Closed"
	HalfOpen         = "HalfOpen"
	Open             = "Open"
	SuccessThreshold = 3
	FailureThreshold = 3
)

type Request struct{}
type Response struct{}

type CircuitBreakerInterface interface {
	Execute(serviceName string, request *Request) (*Response, *MMErr.AppError)
}

func NewCircuitBreaker(client MongoConnect.MongoDBInterface) *CircuitBreaker {
	logger = MMLogger.NewLogger()
	return &CircuitBreaker{client: client}
}

var collection *mongo.Collection

func (cb *CircuitBreaker) Execute(service func(interface{}) (interface{}, interface{}), serviceName string, request interface{}) (interface{}, *MMErr.AppError) {
	db, err := cb.client.GetAppClient()
	if err != nil {
		logger.Info("Unexpected Error", err)
		return nil, err
	}

	collection := db.Collection("CircuitBreaker")

	defer db.Client().Disconnect(context.TODO())

	filter := bson.M{"api_name": serviceName}

	findErr := collection.FindOne(context.Background(), filter).Decode(&detail)
	if err != nil {
		logger.Info("Unexpected Error", findErr.Error())
		return nil, MMErr.NewUnexpectedError(findErr.Error())
	}
	ShouldAttempt := time.Now().After(detail.NextAttempt)
	var response interface{}
	if detail.CurrentState == Closed || (detail.CurrentState == HalfOpen && ShouldAttempt) {
		response, _ := service(request)
		response.statusCode
		if response["StatusCode"] != 500 {
			cb.UpdateCircuitBreakerState(serviceName, detail)
			return response, nil
		} else {
			detail.FailureCount++
			if detail.FailureCount > 3 {
				detail.CurrentState = Open
				detail.NextAttempt = time.Now()
			}
			cb.UpdateCircuitBreakerState(serviceName, detail)
			return nil, cb.HandleFailure(serviceName, "Unexpected Error")
		}
	} else if detail.CurrentState == Open && !ShouldAttempt {
		return nil, cb.HandleFailure(serviceName, "Circuit State is Open")
	}

	return response, nil
}

func (cb *CircuitBreaker) UpdateCircuitBreakerState(serviceName string, detail CircuitBreakerDetail) {

	filter := bson.M{"api_name": serviceName}
	update := bson.M{"$set": bson.M{"current_state": detail.CurrentState}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		logger.Info("Unexpected Error", err.Error())
	}
}

func (cb *CircuitBreaker) GetCircuitBreakerState(serviceName string) string {
	filter := bson.M{"api_name": serviceName}
	var detail dto.CircuitBreakerDetail
	err := collection.FindOne(context.Background(), filter).Decode(&detail)
	if err != nil {
		logger.Info("Unexpected Error", err.Error())
	}

	return detail.CurrentState
}

func (cb *CircuitBreaker) HandleFailure(serviceName string, message string) *MMErr.AppError {
	return MMErr.NewUnexpectedError(message)
}
