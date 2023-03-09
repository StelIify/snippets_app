package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetapp.olex/internal/models"
)

func (app *application) home(w http.ResponseWriter, request *http.Request){
	if request.URL.Path != "/"{
		app.notFound(w)
		return
	}

	latestSnippets, err := app.snippets.Latest()

	if err != nil{
		app.serverError(w, err)
		return
	}
	app.render(w, http.StatusOK, "home.html", &templateData{Snippets: latestSnippets})
}

func (app *application) snippetView(w http.ResponseWriter, request *http.Request){
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1{
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil{
		if errors.Is(err, models.ErrNoRecord){
			app.notFound(w)
		}else{
			app.serverError(w, err)
			return
		}
	}

	app.render(w, http.StatusOK, "view.html", &templateData{Snippet: snippet})
}

func (app *application) snippetCreate(w http.ResponseWriter, request *http.Request){
	if request.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	id, err := app.snippets.Insert("Go database with dynamic expires 2", "some content 2", 10)
	
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, request, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	
}