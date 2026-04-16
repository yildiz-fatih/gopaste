package main

import (
	"errors"
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
	if err != nil {
		return nil, err
	}
	if len(tmplFiles) == 0 {
		return nil, errors.New("No template files found")
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
