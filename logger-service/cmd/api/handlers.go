package main

import (
	"log"
	"log-service/data"
	"net/http"

	"github.com/suhel-kap/toolbox"
)

type JsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	log.Println("WriteLog requested")
	tools := toolbox.Tools{}

	// read json into var
	var requestPayload JsonPayload
	_ = tools.ReadJSON(w, r, &requestPayload)

	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	log.Printf("Inserting log entry: %v", event)

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		log.Println("Error inserting log entry")
		log.Println(err)
		tools.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	resp := toolbox.JsonResponse{
		Error:   false,
		Message: "logged",
	}

	log.Println("WriteLog completed")

	tools.WriteJSON(w, http.StatusAccepted, resp)
	log.Println("WriteLog response sent")
}
