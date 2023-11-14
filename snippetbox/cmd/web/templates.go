package main

import "github.com/johnmerga/Mastering_Go/snippetbox/internal/models"

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
