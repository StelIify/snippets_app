package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"snippetapp.olex/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	latestSnippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = latestSnippets
	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
			return
		}
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.html", data)
}
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "create.html", data)
	//code to show html form to user
}
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
