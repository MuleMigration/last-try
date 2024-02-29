package Controller

import (
	"encoding/json"
	service "feedback/Service"
	"feedback/dto"
	"net/http"

	"github.com/sony/gobreaker"
)

type GetFeedbackController struct {
	service service.GetFeedBackServiceI
	cb      *gobreaker.CircuitBreaker
}

type GetFeedBackControllerI interface {
	GetAllFeedbacks(response http.ResponseWriter, request *http.Request)
	WriteResponse(w http.ResponseWriter, statusCode int, data interface{})
}

func NewGetFeedbackController(service service.GetFeedBackServiceI, cb *gobreaker.CircuitBreaker) GetFeedbackController {
	return GetFeedbackController{service: service, cb: cb}
}

func (c GetFeedbackController) GetAllFeedbacks(response http.ResponseWriter, request *http.Request) {
	var feedbackRequest dto.Feedback

	json.NewDecoder(request.Body).Decode(&feedbackRequest)

	err := dto.Validate(feedbackRequest)
	if err != nil {
		c.WriteResponse(response, 400, &dto.ErrResponse{StatusCode: 400, StatusMessage: "Bad Request"})
		return
	}

	Feedbackdata, serviceErr := c.cb.Execute(func() (interface{}, error) {
		return c.service.GetAllFeedbacks(feedbackRequest)
	})

	if serviceErr != nil {
		c.WriteResponse(response, 500, &dto.ErrResponse{StatusCode: 500, StatusMessage: serviceErr.Error()})
		return
	}

	c.WriteResponse(response, 200, Feedbackdata)
}

func (c GetFeedbackController) WriteResponse(w http.ResponseWriter, statusCode int, data interface{}) {

	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(data)
}
