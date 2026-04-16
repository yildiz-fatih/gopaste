package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yildiz-fatih/gopaste/internal/models"
)

func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	app.writeTemplate(w, "home.tmpl", nil)
}

func (app *application) handlePasteView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	p, err := app.pasteModel.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.NotFound(w, r)
			return
		} else {
			app.writeServerError(w, err)
			return
		}
	}

	data := templateData{
		Paste: p,
	}

	app.writeTemplate(w, "paste_view.tmpl", data)
}

func (app *application) handlePasteCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.writeClientError(w, http.StatusBadRequest)
		return
	}

	content := r.PostForm.Get("content")
	language := r.PostForm.Get("language")

	expires, err := strconv.Atoi(r.PostForm.Get("expires")) // hours
	if err != nil {
		app.writeClientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.pasteModel.Insert(content, language, expires)
	if err != nil {
		app.writeServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/paste/%d", id), http.StatusSeeOther)
}

func (app *application) handleHelp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("help"))
}
