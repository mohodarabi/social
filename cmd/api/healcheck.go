package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(response http.ResponseWriter, request *http.Request) {
	data := map[string]string{
		"status": "ok",
	}

	if err := writeJson(response, http.StatusOK, data); err != nil {
		writeJsonError(response, http.StatusInternalServerError, err.Error())
	}
}
