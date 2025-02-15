package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"
	"time"

	"github.com/suhel-kap/broker/event"
	"github.com/suhel-kap/broker/logs"
	"github.com/suhel-kap/toolbox"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		app.logEventWithRPC(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		tools.ErrorJSON(w, errors.New("unknown action"), http.StatusBadRequest)
	}
}

////////////////////////GRPC WAY TO COMMUNICATE////////////////////////

func (app *Config) LogEventWithGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	conn, err := grpc.NewClient("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var payload = toolbox.JsonResponse{
		Error:   false,
		Message: "Log has been written with gRPC",
	}

	tools.WriteJSON(w, http.StatusOK, payload)
}

////////////////////////RPC WAY TO COMMUNICATE////////////////////////

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) logEventWithRPC(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var payload = toolbox.JsonResponse{
		Error:   false,
		Message: result,
	}

	tools.WriteJSON(w, http.StatusOK, payload)
}

////////////////////////RABBITMQ WAY TO COMMUNICATE////////////////////////

func (app *Config) logEventWithRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		tools.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var payload = toolbox.JsonResponse{
		Error:   false,
		Message: "Log event sent to RabbitMQ",
	}

	tools.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}

////////////////////////JSON WAY TO COMMUNICATE////////////////////////

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
