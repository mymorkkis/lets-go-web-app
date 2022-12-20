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

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)

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

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}

func getUIPath(subfolder string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, "ui", subfolder), nil
}
