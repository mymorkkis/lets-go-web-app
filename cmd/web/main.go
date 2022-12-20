package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	mux := http.NewServeMux()

	staticPath, err := getUIPath("static")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.Dir(staticPath))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// Can use http.ServerFile() to serve individual file from handler but unlike
	// http.FileServer, it does not sanitize input with filepath.Clean()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on :4000")
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

func getUIPath(subfolder string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, "ui", subfolder), nil
}
