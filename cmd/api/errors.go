package main

import (
	"fmt"
	"net/http"
)

func (app *application) internalServerError(respone http.ResponseWriter, request *http.Request, err error) {
	fmt.Printf("INTERNAL_SERVER_ERROR\n")
	fmt.Printf("METHOD: %s\n", request.Method)
	fmt.Printf("PATH: %s\n", request.URL.Path)
	fmt.Printf("ERRORS: %s\n", err.Error())

	errResponse := map[string]string{
		"error": "internal server error",
	}

	writeJson(respone, http.StatusInternalServerError, errResponse)
}

func (app *application) badRequestError(respone http.ResponseWriter, request *http.Request, err error) {
	fmt.Printf("BAD_REQUEST_ERROR\n")
	fmt.Printf("METHOD: %s\n", request.Method)
	fmt.Printf("PATH: %s\n", request.URL.Path)
	fmt.Printf("ERRORS: %s\n", err.Error())

	errResponse := map[string]string{
		"error": err.Error(),
	}

	writeJson(respone, http.StatusBadRequest, errResponse)
}

func (app *application) notFoundError(respone http.ResponseWriter, request *http.Request, err error) {
	fmt.Printf("NOT_FOUND_ERROR\n")
	fmt.Printf("METHOD: %s\n", request.Method)
	fmt.Printf("PATH: %s\n", request.URL.Path)
	fmt.Printf("ERRORS: %s\n", err.Error())

	errResponse := map[string]string{
		"error": "not found",
	}

	writeJson(respone, http.StatusNotFound, errResponse)
}
