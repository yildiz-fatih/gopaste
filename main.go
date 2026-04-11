package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	host := "0.0.0.0"
	port := 8080

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handleHome)
	mux.HandleFunc("GET /paste/{id}", handlePasteView)
	mux.HandleFunc("POST /paste", handlePasteCreate)
	mux.HandleFunc("GET /help", handleHelp)

	log.Printf("Server running on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), mux)
	log.Fatal(err)
}
