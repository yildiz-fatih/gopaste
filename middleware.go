package main

import "net/http"

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
