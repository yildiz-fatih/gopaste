package main

import (
	"fmt"
	"net/http"
)

func (app *application) logRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("Request received",
			"method", r.Method,
			"uri", r.URL.RequestURI(),
			"remote_addr", r.RemoteAddr,
		)

		next.ServeHTTP(w, r)
	}
	handler := http.HandlerFunc(fn)
	return handler
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				w.Header().Set("Connection", "close")
				app.writeServerError(w, fmt.Errorf("%v", err))
			}
		}()

		next.ServeHTTP(w, r)
	}
	handler := http.HandlerFunc(fn)
	return handler
}
