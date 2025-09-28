package main

import (
	"fmt"
	"net/http"
)

func (app *application) internalServerError(respone http.ResponseWriter, request *http.Request, err error) {
	fmt.Printf("INTERNAL_SERVER_ERROR")
	fmt.Printf("METHOD: %s", request.Method)
	fmt.Printf("PATH: %s", request.URL.Path)
	fmt.Printf("ERRORS: %s", err.Error())

	writeJson(respone, http.StatusInternalServerError, "internal server error")
}

func (app *application) badRequestError(respone http.ResponseWriter, request *http.Request, err error) {
	fmt.Printf("BAD_REQUEST_ERROR")
	fmt.Printf("METHOD: %s", request.Method)
	fmt.Printf("PATH: %s", request.URL.Path)
	fmt.Printf("ERRORS: %s", err.Error())

	writeJson(respone, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(respone http.ResponseWriter, request *http.Request, err error) {
	fmt.Printf("NOT_FOUND_ERROR")
	fmt.Printf("METHOD: %s", request.Method)
	fmt.Printf("PATH: %s", request.URL.Path)
	fmt.Printf("ERRORS: %s", err.Error())

	writeJson(respone, http.StatusNotFound, "not found")
}
