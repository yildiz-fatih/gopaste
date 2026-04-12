package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/yildiz-fatih/gopaste/internal/models"
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

// placeholder code for now
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

	fmt.Fprintf(w, "ID: %d,\nContent: %s,\nLanguage: %s,\nCreated: %s,\nExpires: %s\n", p.ID, p.Content, p.Language, p.Created, p.Expires)
}

// placeholder code for now
func (app *application) handlePasteCreate(w http.ResponseWriter, r *http.Request) {
	content := "example content"
	language := "plaintext"
	expires := 7 // days

	query := `INSERT INTO pastes (content, language, created, expires) 
	VALUES ($1, $2, NOW(), NOW() + $3 * INTERVAL '1 day')
	RETURNING id`

	var id int
	err := app.db.QueryRow(query, content, language, expires).Scan(&id)
	if err != nil {
		app.writeServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/paste/%d", id), http.StatusSeeOther)
}

func (app *application) handleHelp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("help"))
}
