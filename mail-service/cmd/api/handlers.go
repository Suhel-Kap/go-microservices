package main

import (
	"net/http"

	"github.com/suhel-kap/toolbox"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := app.tool.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.tool.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		app.tool.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := toolbox.JsonResponse{
		Error:   false,
		Message: "Mail sent successfully" + requestPayload.To,
	}

	app.tool.WriteJSON(w, http.StatusAccepted, payload)
}
