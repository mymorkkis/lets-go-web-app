package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// TODO Replace this with .env vars
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Could also store multiple settings in a single struct
	flag.Parse()

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

	log.Printf("Starting server on %s", *addr)
	err = http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

func getUIPath(subfolder string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, "ui", subfolder), nil
}
