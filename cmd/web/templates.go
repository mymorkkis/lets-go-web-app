package main

import (
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/mymorkkis/lets-go-web-app/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any // TODO Fix this
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	htmlPath := filepath.Join(wd, "ui", "html")

	pages, err := filepath.Glob(
		filepath.Join(htmlPath, "pages", "*.html"),
	)
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		templateSet := template.New(name).Funcs(functions)

		files := []string{
			filepath.Join(htmlPath, "base.html"),
			page,
		}

		templateSet, err := templateSet.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		templateSet, err = templateSet.ParseGlob(
			filepath.Join(htmlPath, "partials", "*.html"),
		)
		if err != nil {
			return nil, err
		}

		cache[name] = templateSet
	}

	return cache, nil
}
