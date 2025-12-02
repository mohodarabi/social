package main

import (
	"net/http"
)

func (app *application) internalServerError(respone http.ResponseWriter, request *http.Request, err error) {
	app.logger.Infow("INTERNAL_SERVER_ERROR\n", "METHOD: ", request.Method, "PATH: ", request.URL.Path, "ERRORS: ", err.Error())
	writeJson(respone, http.StatusInternalServerError, "internal server error")
}

func (app *application) badRequestError(respone http.ResponseWriter, request *http.Request, err error) {
	app.logger.Infow("BAD_REQUEST_ERROR\n", "METHOD: ", request.Method, "PATH: ", request.URL.Path, "ERRORS: ", err.Error())
	writeJsonError(respone, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(respone http.ResponseWriter, request *http.Request, err error) {
	app.logger.Infow("NOT_FOUND_ERROR\n", "METHOD: ", request.Method, "PATH: ", request.URL.Path, "ERRORS: ", err.Error())
	writeJsonError(respone, http.StatusNotFound, "not found")
}
