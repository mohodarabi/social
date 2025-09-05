package main

import "net/http"

func (app *application) healthCheckHandler(response http.ResponseWriter, request *http.Request) {

	response.Write([]byte("im healthy, dont worry :)"))

}
