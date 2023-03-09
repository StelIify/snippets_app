package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetapp.olex/internal/models"
)

func (app *application) home(writer http.ResponseWriter, request *http.Request){
	if request.URL.Path != "/"{
		app.notFound(writer)
		return
	}

	latestSnippets, err := app.snippets.Latest()

	if err != nil{
		app.serverError(writer, err)
		return
	}
	app.render(writer, http.StatusOK, "home.html", &templateData{Snippets: latestSnippets})
}

func (app *application) snippetView(writer http.ResponseWriter, request *http.Request){
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1{
		app.notFound(writer)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil{
		if errors.Is(err, models.ErrNoRecord){
			app.notFound(writer)
		}else{
			app.serverError(writer, err)
			return
		}
	}

	app.render(writer, http.StatusOK, "view.html", &templateData{Snippet: snippet})
}

func (app *application) snippetCreate(writer http.ResponseWriter, request *http.Request){
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", "POST")
		app.clientError(writer, http.StatusMethodNotAllowed)
		return
	}
	id, err := app.snippets.Insert("Go database with dynamic expires 2", "some content 2", 10)
	
	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	
}