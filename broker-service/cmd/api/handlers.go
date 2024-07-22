package main

import (
	"net/http"

	"github.com/suhel-kap/toolbox"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	tools := toolbox.Tools{}

	payload := toolbox.JsonResponse{
		Error:   false,
		Message: "Broker service is running",
	}

	_ = tools.WriteJSON(w, http.StatusOK, payload)
}
