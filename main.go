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

	log.Printf("Server running on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), mux)
	log.Fatal(err)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello GoPaste"))
}
