package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// TODO Replace this with .env vars
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Could also store multiple settings in a single struct
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	mux := http.NewServeMux()

	staticPath, err := app.getUIPath("static")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.Dir(staticPath))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// Can use http.ServerFile() to serve individual file from handler but unlike
	// http.FileServer, it does not sanitize input with filepath.Clean()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog, // ensure the server is using the new errorLog
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}
