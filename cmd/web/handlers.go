package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/yildiz-fatih/gopaste/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	app.writeTemplate(w, "home.tmpl", nil)
}

func (app *application) handlePasteView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	p, err := app.getPaste(r.Context(), slug)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.NotFound(w, r)
			return
		} else {
			app.writeServerError(w, err)
			return
		}
	}

	var data templateData

	if p.PasswordHash != nil {
		data = templateData{}
	} else {
		data = templateData{
			Paste: p,
		}
	}

	data.FullURL = fmt.Sprintf("%s/paste/%s", app.baseURL, p.Slug)

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

	var hashedPassword *string

	password := r.PostForm.Get("password")
	if strings.TrimSpace(password) != "" {
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			app.writeServerError(w, err)
			return
		}

		hashedString := string(hashedBytes)
		hashedPassword = &hashedString
	}

	// write to database
	paste, err := app.pasteModel.Insert(content, expires, hashedPassword)
	if err != nil {
		app.writeServerError(w, err)
		return
	}

	// write to redis
	pasteJson, err := json.Marshal(paste)
	if err != nil {
		app.logger.Error("redis marshal error", "error", err)
	} else {
		err = app.redisClient.Set(r.Context(), fmt.Sprintf("paste:%s", paste.Slug), pasteJson, time.Until(paste.Expires))
		if err != nil {
			app.logger.Error("redis set error", "error", err)
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/paste/%s", paste.Slug), http.StatusSeeOther)
}

func (app *application) handlePasteUnlock(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	err := r.ParseForm()
	if err != nil {
		app.writeClientError(w, http.StatusBadRequest)
		return
	}

	p, err := app.getPaste(r.Context(), slug)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.NotFound(w, r)
			return
		} else {
			app.writeServerError(w, err)
			return
		}
	}

	if p.PasswordHash == nil {
		data := templateData{
			Paste:   p,
			FullURL: fmt.Sprintf("%s/paste/%s", app.baseURL, p.Slug),
		}
		app.writeTemplate(w, "paste_view.tmpl", data)
		return
	}

	enteredPassword := r.PostForm.Get("password")

	err = bcrypt.CompareHashAndPassword([]byte(*p.PasswordHash), []byte(enteredPassword))
	if err != nil {
		data := templateData{
			Error:   "incorrect password",
			FullURL: fmt.Sprintf("%s/paste/%s", app.baseURL, p.Slug),
		}
		app.writeTemplate(w, "paste_view.tmpl", data)
		return
	}

	data := templateData{
		Paste:   p,
		FullURL: fmt.Sprintf("%s/paste/%s", app.baseURL, p.Slug),
	}

	app.writeTemplate(w, "paste_view.tmpl", data)
}
