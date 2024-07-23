package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/suhel-kap/toolbox"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

var tools = toolbox.Tools{}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := toolbox.JsonResponse{
		Error:   false,
		Message: "Broker service is running",
	}

	_ = tools.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItem(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		tools.ErrorJSON(w, errors.New("unknown action"), http.StatusBadRequest)
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		tools.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	} else if response.StatusCode != http.StatusAccepted {
		tools.ErrorJSON(w, errors.New("unknown error"), http.StatusBadRequest)
		return
	}

	// create a var we will read the response into
	var jsonFromService toolbox.JsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if jsonFromService.Error {
		tools.ErrorJSON(w, errors.New(jsonFromService.Message), http.StatusUnauthorized)
		return
	}

	var payload toolbox.JsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	tools.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	logServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		tools.ErrorJSON(w, errors.New("unknown error"), http.StatusBadRequest)
		return
	}

	var jsonFromService toolbox.JsonResponse
	jsonFromService.Error = false
	jsonFromService.Message = "Logged!"

	tools.WriteJSON(w, http.StatusAccepted, jsonFromService)
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// call the mail service
	mailServiceUrl := "http://mail-service/send"

	request, err := http.NewRequest("POST", mailServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		tools.ErrorJSON(w, errors.New("unknown error"), http.StatusBadRequest)
		return
	}

	var jsonFromService toolbox.JsonResponse
	jsonFromService.Error = false
	jsonFromService.Message = "Mail sent!"

	tools.WriteJSON(w, http.StatusAccepted, jsonFromService)
}
