package main

import (
	"html/template"
	"path/filepath"

	"github.com/yildiz-fatih/gopaste/internal/models"
)

type templateData struct {
	Paste models.Paste
}

func parseTemplates() (map[string]*template.Template, error) {
	parsed := make(map[string]*template.Template)

	tmplFiles, err := filepath.Glob("./views/pages/*.tmpl")
	if err != nil || len(tmplFiles) == 0 {
		return nil, err
	}

	for _, file := range tmplFiles {
		tmpl, err := template.ParseFiles("./views/base.tmpl", file)
		if err != nil {
			return nil, err
		}

		filename := filepath.Base(file)
		parsed[filename] = tmpl
	}

	return parsed, nil
}
