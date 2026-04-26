package models

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"math/big"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type Paste struct {
	ID      int
	Slug    string
	Content string
	Created time.Time
	Expires time.Time
}

type PasteModel struct {
	DB *sql.DB
}

func randomSlug(length int) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

	result := make([]byte, length)
	for i := range length {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		result[i] = alphabet[n.Int64()]
	}

	return string(result), nil
}

func (m *PasteModel) Get(slug string) (Paste, error) {
	query := `SELECT id, slug, content, created, expires 
	FROM pastes 
	WHERE expires > NOW() AND slug = $1`

	var p Paste

	err := m.DB.QueryRow(query, slug).Scan(&p.ID, &p.Slug, &p.Content, &p.Created, &p.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Paste{}, ErrNotFound
		} else {
			return Paste{}, err
		}
	}

	return p, nil
}

func (m *PasteModel) Insert(content string, expires int) (string, error) {
	query := `INSERT INTO pastes (slug, content, created, expires) 
	VALUES ($1, $2, NOW(), NOW() + $3 * INTERVAL '1 hour')
	RETURNING slug`

	const slugLength = 6
	const maxRetries = 3

	// try to insert
	for range maxRetries {
		slug, err := randomSlug(slugLength)
		if err != nil {
			return "", err
		}

		var insertedSlug string
		err = m.DB.QueryRow(query, slug, content, expires).Scan(&insertedSlug)
		if err != nil {
			// (expected) unique violation error, run the loop again
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				continue
			}

			// (unexpected) database error
			return "", err
		}

		return insertedSlug, nil

	}

	// exhausted retries, return error
	return "", errors.New("failed to generate unique slug")
}
