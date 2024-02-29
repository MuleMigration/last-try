package OpenAI

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	MMLogger "postFeedback/logger"
	MMErr "postFeedback/mmerror"
)

var log *MMLogger.Logger

type ChatRequest struct {
	Request string
}

type ChatResponse struct {
	Response string
}

type OpenAI struct{}

type OpenAiInterface interface {
	GetChatResponse(chatRequest *ChatRequest) (*ChatResponse, *MMErr.AppError)
}

func NewOpenAI() OpenAI {
	return OpenAI{}
}

const API_KEY = "your_api_key"
const API_URL = "https://api.openai.com/v1/engines/davinci-codex/completions"

func (*OpenAI) GetChatResponse(chatRequest *ChatRequest) (*ChatResponse, *MMErr.AppError) {
	log = MMLogger.NewLogger()
	requestBody, err := json.Marshal(chatRequest.Request)
	if err != nil {
		log.Info("Internal Server Error", err.Error())
		return nil, MMErr.NewUnexpectedError(err.Error())
	}

	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Info("Internal Server Error", err.Error())
		return nil, MMErr.NewUnexpectedError(err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Info("Internal Server Error", err.Error())
		return nil, MMErr.NewUnexpectedError(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Info("Internal Server Error", err.Error())
		return nil, MMErr.NewUnexpectedError(errors.New("received non-200 response from OpenAI Chat API").Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Info("Internal Server Error", err.Error())
		return nil, MMErr.NewUnexpectedError(err.Error())
	}

	var chatResponse ChatResponse
	err = json.Unmarshal(body, &chatResponse.Response)
	if err != nil {
		log.Info("Internal Server Error", err.Error())
		return nil, MMErr.NewUnexpectedError(errors.New("failed to unmarshal response body").Error())
	}

	return &chatResponse, nil
}
