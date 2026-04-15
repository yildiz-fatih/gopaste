package main

import "net/http"

func (app *application) writeServerError(w http.ResponseWriter, err error) {
	app.logger.Error(err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // TODO: make a custom error page
}

func (app *application) writeClientError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}
