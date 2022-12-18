package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	basePath, err := getHTMLPath()
	if err != nil {
		log.Print(err.Error())
		errCode := http.StatusInternalServerError
		http.Error(w, http.StatusText(errCode), errCode)
	}

	files := []string{
		filepath.Join(basePath, "base.html"),
		filepath.Join(basePath, "pages", "home.html"),
		filepath.Join(basePath, "partials", "nav.html"),
	}

	templateSet, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		errCode := http.StatusInternalServerError
		http.Error(w, http.StatusText(errCode), errCode)
		return
	}

	err = templateSet.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		errCode := http.StatusInternalServerError
		http.Error(w, http.StatusText(errCode), errCode)
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID: %d", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		errCode := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(errCode), errCode)
		return
	}

	w.Write([]byte("Create a new snippet..."))
}

func getHTMLPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, "ui", "html"), nil
}
