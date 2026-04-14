package main

import "net/http"

func (app *application) newRouter() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("GET /{$}", app.handleHome)
	mux.HandleFunc("GET /paste/{id}", app.handlePasteView)
	mux.HandleFunc("POST /paste", app.handlePasteCreate)
	mux.HandleFunc("GET /help", app.handleHelp)

	// return mux
	return app.logRequest(mux)
}
