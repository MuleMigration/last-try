package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	service "postFeedback/Service"
	"postFeedback/dto"
)

type PostFeedbackController struct {
	service service.PostFeedBackServiceI
}

type PostFeedBackControllerI interface {
	PostFeedback(w http.ResponseWriter, r *http.Request)
	WriteResponse(w http.ResponseWriter, statusCode int, data interface{})
}

func NewFeedbackController(service service.PostFeedBackServiceI) PostFeedbackController {
	return PostFeedbackController{service: service}
}

func (c *PostFeedbackController) PostFeedback(w http.ResponseWriter, r *http.Request) {
	var feedback dto.Request

	json.NewDecoder(r.Body).Decode(&feedback)

	err := dto.Validate(feedback)
	if err != nil {
		c.WriteResponse(w, 400, &dto.Response{StatusCode: "400", StatusMessage: "Bad Request"})
		return
	}

	result, err := c.service.PostFeedback(feedback)
	if err != nil {
		c.WriteResponse(w, 400, &dto.Response{StatusCode: "500", StatusMessage: "Failed to Insert the Feedback."})
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

	json.NewEncoder(w).Encode(data)
}
