package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

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

func handlePasteView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	res := fmt.Sprintf("display a paste with ID %s", id)
	w.Write([]byte(res))
}

func handlePasteCreate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("creating paste..."))
}

func handleHelp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("help"))
}
