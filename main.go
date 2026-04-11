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

	app.logger.Info("Starting server", "host", host, "port", port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), app.newRouter()) // err is always non-nil
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}
