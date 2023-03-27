package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"snippetapp.olex/internal/models"
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

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(dbpool)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	templateCache, err := newTemplateCache()

	if err != nil {
		erorrLog.Fatal(err)
	}

	app := &application{
		infoLog:        infoLog,
		errorLog:       erorrLog,
		snippets:       &models.SnippetModel{DB: dbpool},
		users:          &models.UserModel{DB: dbpool},
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}

	server := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: erorrLog,
	}

	infoLog.Printf("Starting server on %s", *addr)
	// err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	err = server.ListenAndServe()

	erorrLog.Fatal(err)
}

type application struct {
	infoLog        *log.Logger
	errorLog       *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}
