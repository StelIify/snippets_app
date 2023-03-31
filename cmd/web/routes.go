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

	r.Method(http.MethodGet, "/static/*", http.StripPrefix("/static/", fileserver))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	r.Method(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	r.Method(http.MethodGet, "/snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	r.Method(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignUp))
	r.Method(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignUpPost))
	r.Method(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	r.Method(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	r.Method(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	r.Method(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	r.Method(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, securityHeaders)

	return standard.Then(r)
}
