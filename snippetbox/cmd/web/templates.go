package main

import (
	"html/template"
	"path/filepath"

	"github.com/johnmerga/Mastering_Go/snippetbox/internal/models"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)
	pages, err := filepath.Glob(".tmpl.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		files := []string{
			"../../ui/html/base.tmpl.html",
			"../../ui/html/partials/nav.tmpl.html",
			page,
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
