package main

import "net/http"

func (app *application) routes() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	staticPath, err := app.getUIPath("static")
	if err != nil {
		return nil, err
	}

	// Can use http.ServerFile() to serve individual file from handler but unlike
	// http.FileServer, it does not sanitize input with filepath.Clean()
	fileServer := http.FileServer(http.Dir(staticPath))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux, nil
}
