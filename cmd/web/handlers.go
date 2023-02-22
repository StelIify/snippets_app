package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(writer http.ResponseWriter, request *http.Request){
	if request.URL.Path != "/"{
		http.NotFound(writer, request)
		return
	}
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}
	ts, err := template.ParseFiles(files...)

	if err != nil{
		log.Println(err.Error())
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(writer, "base", nil)
	if err != nil{
		log.Println(err.Error())
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func snippetView(writer http.ResponseWriter, request *http.Request){
	id, error := strconv.Atoi(request.URL.Query().Get("id"))
	if error != nil || id < 1{
		http.NotFound(writer, request)
		return
	}
	fmt.Fprintf(writer, "Displaying a specific snipper with id %d..", id)
}

func snippetCreate(writer http.ResponseWriter, request *http.Request){
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", "POST")
		http.Error(writer, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}
	writer.Write([]byte("Hello from snippet create"))
}