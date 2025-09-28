package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJson(response http.ResponseWriter, status int, data any) error {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	return json.NewEncoder(response).Encode(data)
}

func readJson(response http.ResponseWriter, request *http.Request, data any) error {
	maxBytes := 1_048_578

	request.Body = http.MaxBytesReader(response, request.Body, int64(maxBytes))

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func writeJsonError(response http.ResponseWriter, status int, message string) error {

	type envelope struct {
		Error string `json:"error"`
	}

	return writeJson(response, status, &envelope{Error: message})

}
