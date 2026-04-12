package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

type application struct {
	logger *slog.Logger
	db     *sql.DB
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

	app := &application{
		logger: logger,
		db:     db,
	}

	logger.Info("Starting server", "host", host, "port", port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), app.newRouter()) // err is always non-nil
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
