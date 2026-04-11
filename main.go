package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	host := "0.0.0.0"
	var port int

	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.Parse()

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("GET /{$}", handleHome)
	mux.HandleFunc("GET /paste/{id}", handlePasteView)
	mux.HandleFunc("POST /paste", handlePasteCreate)
	mux.HandleFunc("GET /help", handleHelp)

	log.Printf("Server running on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), mux)
	log.Fatal(err)
}
