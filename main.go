package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	host := "0.0.0.0"
	var port int

	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.Parse()

	app := &application{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("GET /{$}", app.handleHome)
	mux.HandleFunc("GET /paste/{id}", app.handlePasteView)
	mux.HandleFunc("POST /paste", app.handlePasteCreate)
	mux.HandleFunc("GET /help", app.handleHelp)

	app.logger.Info("Starting server", "host", host, "port", port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), mux) // err is always non-nil
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}
