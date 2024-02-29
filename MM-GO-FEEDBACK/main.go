package main

import (
	"net/http"
	controller "postFeedback/Controller"
	repository "postFeedback/Repository"
	service "postFeedback/Service"
	MongoConnect "postFeedback/mongoconnect"

	"github.com/gorilla/mux"
)

func main() {
	Database := MongoConnect.MongoDB{}
	Repository := repository.NewFeedbackRepository(&Database)
	Service := service.NewFeedbackService(&Repository)
	Controller := controller.NewFeedbackController(&Service)

	router := mux.NewRouter()
	router.HandleFunc("/postfeedback", Controller.PostFeedback).Methods("POST")
	http.ListenAndServe("localhost:5001", router)
}
