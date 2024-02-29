package main

import (
	controller "feedback/Controller"
	repository "feedback/Repository"
	service "feedback/Service"
	config "feedback/circuitbreaker"
	MongoConnect "feedback/mongoconnect"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	Database := MongoConnect.MongoDB{}
	//post
	PostRepository := repository.NewPostFeedbackRepository(&Database)
	PostService := service.NewPostFeedbackService(&PostRepository)
	circuitBreaker := config.CircuitBreakerConfig()
	PostController := controller.NewPostFeedbackController(&PostService, circuitBreaker)

	//get
	Repository := repository.NewGetFeedbackRepository(&Database)
	GetService := service.NewGetFeedbackService(Repository)
	cb := config.CircuitBreakerConfig()

	GetController := controller.NewGetFeedbackController(GetService, cb)

	router := mux.NewRouter()
	//
	router.HandleFunc("/postfeedback", PostController.PostFeedback).Methods("POST")

	router.HandleFunc("/getfeedbacks", GetController.GetAllFeedbacks).Methods("GET")
	fmt.Println("Service Running in port 5001")

	if err := http.ListenAndServe("localhost:5001", router); err != nil {
		fmt.Println("error from the server: ", err)
	}
}
