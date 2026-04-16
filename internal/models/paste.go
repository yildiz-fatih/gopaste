package models

import (
	"database/sql"
	"errors"
	"time"
)

type Paste struct {
	ID       int
	Content  string
	Language string
	Created  time.Time
	Expires  time.Time
}

type PasteModel struct {
	DB *sql.DB
}

func (m *PasteModel) Get(id int) (Paste, error) {
	query := `SELECT id, content, language, created, expires 
	FROM pastes 
	WHERE expires > NOW() AND id = $1`

	var p Paste

	err := m.DB.QueryRow(query, id).Scan(&p.ID, &p.Content, &p.Language, &p.Created, &p.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// http.NotFound(w, r)
			return Paste{}, err
		} else {
			// app.writeServerError(w, err)
			return Paste{}, err
		}
	}

	return p, nil
}

func (m *PasteModel) Insert(content string, language string, expires int) (int, error) {
	query := `INSERT INTO pastes (content, language, created, expires) 
	VALUES ($1, $2, NOW(), NOW() + $3 * INTERVAL '1 hour')
	RETURNING id`

	var id int
	err := m.DB.QueryRow(query, content, language, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
