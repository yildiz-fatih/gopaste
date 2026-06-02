package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/yildiz-fatih/gopaste/internal/models"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 20 * time.Second
	idleTimeout       = 120 * time.Second
)

type application struct {
	logger      *slog.Logger
	pasteModel  *models.PasteModel
	templates   map[string]*template.Template
	baseURL     string
	redisClient *redis.Client
}

func main() {
	host := "0.0.0.0"

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Error("DATABASE_URL is not set")
		os.Exit(1)
	}

	portStr := os.Getenv("PORT")
	if portStr == "" {
		logger.Error("PORT is not set")
		os.Exit(1)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		logger.Error("Invalid PORT value")
		os.Exit(1)
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		logger.Error("BASE_URL is not set")
		os.Exit(1)
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		logger.Error("REDIS_URL is not set")
		os.Exit(1)
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("Connected to the database")

	redisOptions, err := redis.ParseURL(redisURL)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	redisClient := redis.NewClient(redisOptions)
	defer redisClient.Close()

	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("Connected to Redis")

	// upsert the help paste (insert if not exists, update if it does)
	helpPaste, err := os.ReadFile("help.md")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	query := `INSERT INTO pastes (slug, content, created, expires) 
	VALUES ('help', $1, NOW(), NOW() + INTERVAL '100 years')
	ON CONFLICT (slug) DO UPDATE SET content = EXCLUDED.content`
	_, err = db.Exec(query, string(helpPaste))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("Help paste upserted")

	parsedTemplates, err := parseTemplates()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:      logger,
		pasteModel:  &models.PasteModel{DB: db},
		templates:   parsedTemplates,
		baseURL:     baseURL,
		redisClient: redisClient,
	}

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", host, port),
		Handler:           app.newRouter(),
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	logger.Info("Starting server", "host", host, "port", port)
	err = server.ListenAndServe() // err is always non-nil
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
