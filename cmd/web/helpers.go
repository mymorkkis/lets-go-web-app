package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) getUIPath(subfolder string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, "ui", subfolder), nil
}

func (app *application) getTemplateSetForHTMLPage(pageName string) (*template.Template, error) {
	htmlPath, err := app.getUIPath("html")
	if err != nil {
		return nil, err
	}

	page := fmt.Sprintf("%s.html", pageName)

	files := []string{
		filepath.Join(htmlPath, "base.html"),
		filepath.Join(htmlPath, "partials", "nav.html"),
		filepath.Join(htmlPath, "pages", page),
	}

	templateSet, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return templateSet, nil
}
