package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	fileserver := http.FileServer(http.Dir("./ui/static/"))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	r.Method(http.MethodGet, "/static/", http.StripPrefix("static", fileserver))

	r.Get("/", app.home)
	r.Get("/snippet/view/{id}", app.snippetView)
	r.Get("/snippet/create", app.snippetCreate)
	r.Post("/snippet/create", app.snippetCreatePost)

	standard := alice.New(app.recoverPanic, app.logRequest, securityHeaders)
	return standard.Then(r)
}
