package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
	"snippetapp.olex/internal/models"
	"html/template"
)


func main() {
	addr := flag.String("addr", ":8000", "HTTP network address")
	flag.Parse()
	
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	erorrLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	
	dbpool, db_error := pgxpool.New(context.Background(), os.Getenv("SNIPPET_DB_URL"))
	if db_error != nil {
		erorrLog.Fatal(db_error)
	}
	defer dbpool.Close()

	templateCache, err := newTemplateCache()
	if err != nil{
		erorrLog.Fatal(err)
	}

	app := &application{
		infoLog:  infoLog,
		errorLog: erorrLog,
		snippets: &models.SnippetModel{DB: dbpool},
		templateCache: templateCache,
	}

	server := &http.Server{
		Addr: *addr,
		Handler: app.routes(),
		ErrorLog: erorrLog,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServe()

	erorrLog.Fatal(err)
}

type application struct {
	infoLog *log.Logger
	errorLog *log.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template 
}