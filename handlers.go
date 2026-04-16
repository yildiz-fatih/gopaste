package main

import (
	"database/sql"
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

	query := `SELECT id, content, language, created, expires 
	FROM pastes 
	WHERE expires > NOW() AND id = $1`

	var p models.Paste
	err = app.db.QueryRow(query, id).Scan(&p.ID, &p.Content, &p.Language, &p.Created, &p.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			app.writeServerError(w, err)
		}
		return
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

	query := `INSERT INTO pastes (content, language, created, expires) 
	VALUES ($1, $2, NOW(), NOW() + $3 * INTERVAL '1 hour')
	RETURNING id`

	var id int
	err = app.db.QueryRow(query, content, language, expires).Scan(&id)
	if err != nil {
		app.writeServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/paste/%d", id), http.StatusSeeOther)
}

func (app *application) handleHelp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("help"))
}
