package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	tmplFiles := []string{
		"./views/base.tmpl", // base must be parsed first
		"./views/home.tmpl",
	}
	tmpl, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		app.writeServerError(w, err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.writeServerError(w, err)
		return
	}
}

func (app *application) handlePasteView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	res := fmt.Sprintf("display a paste with ID %s", id)
	w.Write([]byte(res))
}

func (app *application) handlePasteCreate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("creating paste..."))
}

func (app *application) handleHelp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("help"))
}
