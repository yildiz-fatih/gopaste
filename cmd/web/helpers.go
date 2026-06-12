package main

import (
	"errors"
	"net/http"
)

func (app *application) writeServerError(w http.ResponseWriter, err error) {
	app.logger.Error(err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // TODO: make a custom error page
}

func (app *application) writeClientError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (app *application) writeTemplate(w http.ResponseWriter, name string, data any) {
	tmpl, ok := app.templates[name]
	if !ok {
		app.writeServerError(w, errors.New("Template not found in cache"))
		return
	}

	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.writeServerError(w, err)
		return
	}
}
