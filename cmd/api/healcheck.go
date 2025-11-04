package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(response http.ResponseWriter, request *http.Request) {
	data := map[string]string{
		"status": "ok",
	}

	if err := writeJsonResponse(response, http.StatusOK, data); err != nil {
		app.internalServerError(response, request, err)
	}
}
