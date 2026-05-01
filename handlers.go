package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/yildiz-fatih/gopaste/internal/models"
)

func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	app.writeTemplate(w, "home.tmpl", nil)
}

func (app *application) handlePasteView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	p, err := app.pasteModel.Get(slug)
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
		Paste:   p,
		FullURL: fmt.Sprintf("%s/paste/%s", app.baseURL, p.Slug),
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

	expires, err := strconv.Atoi(r.PostForm.Get("expires")) // hours
	if err != nil {
		app.writeClientError(w, http.StatusBadRequest)
		return
	}

	// write to database
	paste, err := app.pasteModel.Insert(content, expires)
	if err != nil {
		app.writeServerError(w, err)
		return
	}

	// write to redis
	pasteJson, err := json.Marshal(paste)
	if err != nil {
		app.writeServerError(w, err)
		return
	}

	err = app.redisClient.Set(r.Context(), fmt.Sprintf("paste:%s", paste.Slug), pasteJson, time.Until(paste.Expires)).Err()
	if err != nil {
		app.logger.Error("Failed to cache paste in redis", "error", err)
	}

	http.Redirect(w, r, fmt.Sprintf("/paste/%s", paste.Slug), http.StatusSeeOther)
}
