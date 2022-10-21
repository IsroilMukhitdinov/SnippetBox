package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	fs := http.FileServer(http.Dir(app.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	return mux
}
