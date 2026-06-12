package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yildiz-fatih/gopaste/internal/models"
)

// implements models.PasteRepository
type mockPasteModel struct {
	paste models.Paste
	err   error
}

func (m *mockPasteModel) Get(slug string) (models.Paste, error) {
	return m.paste, m.err // return the preset paste and error
}

func (m *mockPasteModel) Insert(content string, expires int) (models.Paste, error) {
	return m.paste, m.err // return the preset paste and error
}

// implements models.RedisRepository
type mockRedisClient struct {
	value string
	err   error
}

func (m *mockRedisClient) Get(ctx context.Context, key string) (string, error) {
	return m.value, m.err // return the preset value and error
}

func (m *mockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return nil // do nothing, just return nil
}

func TestHandlePasteView(t *testing.T) {
	slug := "abc123"
	pasteContent := "hello world"

	// The paste our fake DB will "return"
	dbPaste := models.Paste{
		ID:      1,
		Slug:    slug,
		Content: pasteContent,
		Created: time.Now(),
		Expires: time.Now().Add(time.Hour),
	}

	pasteJson, err := json.Marshal(dbPaste)
	if err != nil {
		t.Fatal(err)
	}

	parsedTemplates, err := parseTemplates()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		dbPaste    models.Paste
		dbErr      error
		redisValue string
		redisErr   error
		wantCode   int
		wantBody   string
	}{
		{
			name:       "cache miss, found in db, 200",
			dbPaste:    dbPaste,
			dbErr:      nil,
			redisValue: "",
			redisErr:   redis.Nil, // redis.Nil = "key not found" = cache miss
			wantCode:   http.StatusOK,
			wantBody:   pasteContent,
		},
		{
			name:       "cache miss, not found in db, 404",
			dbErr:      models.ErrNotFound,
			redisValue: "",
			redisErr:   redis.Nil,
			wantCode:   http.StatusNotFound,
		},
		{
			name:       "cache miss, db error, 500",
			dbErr:      errors.New("database error"),
			redisValue: "",
			redisErr:   redis.Nil,
			wantCode:   http.StatusInternalServerError,
		},
		{
			name:       "cache hit, 200",
			dbErr:      errors.New("db should not be called on cache hit"),
			redisValue: string(pasteJson),
			wantCode:   http.StatusOK,
			wantBody:   pasteContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{
				logger: slog.New(slog.DiscardHandler), // discard logs during tests
				pasteModel: &mockPasteModel{
					paste: tt.dbPaste,
					err:   tt.dbErr,
				},
				templates: parsedTemplates,
				baseURL:   "http://localhost:1234",
				redisClient: &mockRedisClient{
					value: tt.redisValue,
					err:   tt.redisErr,
				},
			}

			// make request and record response
			resRecorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/paste/"+slug, nil)
			req.SetPathValue("slug", slug)

			app.handlePasteView(resRecorder, req)

			gotCode := resRecorder.Code
			gotBody := resRecorder.Body.String()

			// assert
			if tt.wantCode != gotCode {
				t.Errorf("status: want %d, got %d", tt.wantCode, gotCode)
			}

			if tt.wantBody != "" && !strings.Contains(gotBody, tt.wantBody) {
				t.Errorf("body: want it to contain %q, got\n%q", tt.wantBody, gotBody)
			}
		})
	}
}

func TestHandlePasteCreate(t *testing.T) {
	/*
		contract:
			in: 	http request (POST, form in body)
			out:	http response (303 redirect, "/paste/createdslug")
	*/
	/*
		happy:
			- "valid form, 303 redirect"
		sad:
			- "invalid form, 400" -> hard to test, ignored
			- "invalid expires field, 400"
			- "db write fails, 500"
	*/
	slug := "abc123"
	pasteContent := "hello world"

	// The paste our fake DB will "return"
	dbPaste := models.Paste{
		ID:      1,
		Slug:    slug,
		Content: pasteContent,
		Created: time.Now(),
		Expires: time.Now().Add(time.Hour),
	}

	tests := []struct {
		name         string       // test name
		formValues   url.Values   // form values for request body
		dbPaste      models.Paste // mock return value
		dbErr        error        // mock return value
		wantCode     int          // expected status code
		wantLocation string       // expected Location header
	}{
		{
			name: "valid form, 303 redirect",
			formValues: url.Values{
				"content": {pasteContent},
				"expires": {"60"},
			},
			dbPaste:      dbPaste,
			dbErr:        nil,
			wantCode:     http.StatusSeeOther,
			wantLocation: "/paste/" + slug,
		},
		{
			name: "invalid expires field, 400",
			formValues: url.Values{
				"content": {pasteContent},
				"expires": {"notanumber"},
			},
			dbErr:    nil,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "db write fails, 500",
			formValues: url.Values{
				"content": {pasteContent},
				"expires": {"60"},
			},
			dbErr:    errors.New("database error"),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{
				logger: slog.New(slog.DiscardHandler), // discard logs during tests
				pasteModel: &mockPasteModel{
					paste: tt.dbPaste,
					err:   tt.dbErr,
				},
				templates:   nil,
				baseURL:     "http://localhost:1234",
				redisClient: &mockRedisClient{},
			}

			// make request and record response
			resRecorder := httptest.NewRecorder()

			body := strings.NewReader(tt.formValues.Encode())
			req := httptest.NewRequest(http.MethodPost, "/paste", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			app.handlePasteCreate(resRecorder, req)

			gotCode := resRecorder.Code
			gotLocation := resRecorder.Header().Get("Location")

			// assert
			if tt.wantCode != gotCode {
				t.Errorf("status: want %d, got %d", tt.wantCode, gotCode)
			}

			if tt.wantLocation != "" && tt.wantLocation != gotLocation {
				t.Errorf("location: want %q, got %q", tt.wantLocation, gotLocation)
			}
		})
	}
}
