package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/mymorkkis/lets-go-web-app/internal/models"
	"github.com/mymorkkis/lets-go-web-app/ui"
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

	pages, err := fs.Glob(
		ui.Files,
		filepath.Join("html", "pages", "*.html"),
	)
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		templateSet := template.New(name).Funcs(functions)

		patterns := []string{
			filepath.Join("html", "base.html"),
			filepath.Join("html", "partials", "*.html"),
			page,
		}

		templateSet, err = templateSet.ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = templateSet
	}

	return cache, nil
}
