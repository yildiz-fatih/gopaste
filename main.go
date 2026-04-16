package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/yildiz-fatih/gopaste/internal/models"
)

type application struct {
	logger     *slog.Logger
	pasteModel *models.PasteModel
	templates  map[string]*template.Template
}

func main() {
	host := "0.0.0.0"
	var port int

	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Error("DATABASE_URL environment variable is not set")
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

	parsedTemplates, err := parseTemplates()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:     logger,
		pasteModel: &models.PasteModel{DB: db},
		templates:  parsedTemplates,
	}

	server := &http.Server{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Handler:  app.newRouter(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("Starting server", "host", host, "port", port)
	err = server.ListenAndServe() // err is always non-nil
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
