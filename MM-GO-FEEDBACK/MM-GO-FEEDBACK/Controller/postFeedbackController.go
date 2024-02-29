package Controller

import (
	"encoding/json"
	service "feedback/Service"
	"feedback/dto"
	MMLogger "feedback/logger"
	"fmt"
	"net/http"

	"github.com/sony/gobreaker"
)

var log *MMLogger.Logger

type PostFeedbackController struct {
	service service.PostFeedBackServiceI
	cb      *gobreaker.CircuitBreaker
}

type PostFeedBackControllerI interface {
	PostFeedback(w http.ResponseWriter, r *http.Request)
	WriteResponse(w http.ResponseWriter, statusCode int, data interface{})
}

func NewPostFeedbackController(service service.PostFeedBackServiceI, cb *gobreaker.CircuitBreaker) PostFeedbackController {
	return PostFeedbackController{service: service, cb: cb}
}

// PC_NO_1.14
func (c *PostFeedbackController) PostFeedback(w http.ResponseWriter, r *http.Request) {
	var feedback dto.Request

	json.NewDecoder(r.Body).Decode(&feedback)
	//PC_NO_1.17
	err := dto.PostValidate(feedback)
	if err != nil {
		c.WriteResponse(w, 400, &dto.PostResponse{StatusCode: 400, StatusMessage: "Bad Request"})
		return
	}

	result, serviceErr := c.cb.Execute(func() (interface{}, error) {
		return c.service.PostFeedback(feedback)

	})
	//PC_NO_1.48 - PC_NO_1.50
	if serviceErr != nil {
		c.WriteResponse(w, 500, serviceErr.Error())
		return
	}

	fmt.Println(result)
	c.WriteResponse(w, 200, result)
}

func (c PostFeedbackController) WriteResponse(w http.ResponseWriter, statusCode int, data interface{}) {

	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	out := data
	if statusCode == 500 {
		out = map[string]interface{}{"StatusCode": statusCode, "StatusMessage": data}
	}

	json.NewEncoder(w).Encode(out)
}
