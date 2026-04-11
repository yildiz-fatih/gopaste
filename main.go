package main

import (
	"fmt"
	"html/template"
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
	tmplFiles := []string{
		"./views/base.tmpl", // base must be parsed first
		"./views/home.tmpl",
	}
	tmpl, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError) // TODO: make a custom error page
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError) // TODO: make a custom error page
		return
	}
}
