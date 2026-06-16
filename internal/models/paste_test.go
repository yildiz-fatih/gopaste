package models

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestRandomSlug(t *testing.T) {
	t.Run("correct length", func(t *testing.T) {
		lengths := []int{1, 8, 16}

		for _, want := range lengths {
			slug, err := randomSlug(want)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got := len(slug)

			if want != got {
				t.Errorf("want %d, got %d", want, got)
			}
		}
	})
	t.Run("error for non-positive length", func(t *testing.T) {
		lengths := []int{0, -1, -100}

		for _, length := range lengths {
			_, err := randomSlug(length)
			if !errors.Is(err, errInvalidLength) {
				t.Errorf("expected an error for length %d", length)
			}
		}
	})
	t.Run("only uses allowed characters", func(t *testing.T) {
		// intentionally duplicated so the test is independent of the implementation
		const allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

		slug, err := randomSlug(50)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for _, char := range slug {
			if !strings.ContainsRune(allowedChars, char) {
				t.Errorf("contains disallowed character %q", char)
			}
		}
	})
}

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// set up
	ctx := context.Background()

	migrationFiles, err := filepath.Glob(filepath.Join("..", "..", "migrations", "*.up.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if len(migrationFiles) == 0 {
		t.Fatal("no migration files found")
	}

	pgContainer, err := postgres.Run(ctx, "postgres:18.3-alpine",
		postgres.WithInitScripts(migrationFiles...),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// connect
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestPasteModel_Get(t *testing.T) {
	t.Run("paste exists", func(t *testing.T) {
		// set up a test database
		db := newTestDB(t)

		// insert a paste
		query := `INSERT INTO pastes (slug, content, created, expires)
	VALUES ($1, $2, NOW(), NOW() + $3 * INTERVAL '1 hour')
	RETURNING *`
		var paste Paste
		slug := "test-slug"
		err := db.QueryRow(query, slug, "test content", 1).Scan(&paste.ID, &paste.Slug, &paste.Content, &paste.Created, &paste.Expires)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// call Get and check the return
		pasteModel := PasteModel{DB: db}
		got, err := pasteModel.Get(slug)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(paste, got) {
			t.Errorf("want %+v, got %+v", paste, got)
		}
	})
	t.Run("paste does not exist", func(t *testing.T) {
		// set up a test database
		db := newTestDB(t)

		// call Get and check the return
		pasteModel := PasteModel{DB: db}
		_, got := pasteModel.Get("doesnotmatterwhatiputhere")
		if !errors.Is(got, ErrNotFound) {
			t.Errorf("want %v, got %v", ErrNotFound, got)
		}
	})
	t.Run("paste expired", func(t *testing.T) {
		// set up a test database
		db := newTestDB(t)

		// insert a paste
		query := `INSERT INTO pastes (slug, content, created, expires)
		VALUES ($1, $2, NOW(), NOW() + $3 * INTERVAL '1 hour')`
		slug := "test-slug"
		_, err := db.Exec(query, slug, "test content", -1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// call Get and check the return
		pasteModel := PasteModel{DB: db}
		_, got := pasteModel.Get(slug)
		if !errors.Is(got, ErrNotFound) {
			t.Errorf("want %v, got %v", ErrNotFound, got)
		}
	})
}

func TestPasteModel_Insert(t *testing.T) {
	t.Run("no collision", func(t *testing.T) {
		// set up a test database
		db := newTestDB(t)

		pasteModel := PasteModel{DB: db}
		wantContent := "test-content"
		wantExpires := 1
		gotPaste, err := pasteModel.Insert(wantContent, wantExpires)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if wantContent != gotPaste.Content {
			t.Errorf("want %q, got %q", wantContent, gotPaste.Content)
		}

		gotExpires := int(gotPaste.Expires.Sub(gotPaste.Created).Hours())
		if wantExpires != gotExpires {
			t.Errorf("want %d, got %d", wantExpires, gotExpires)
		}
	})
}
