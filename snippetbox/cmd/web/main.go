package main

import (
	"database/sql"
	"flag"
	"github.com/fatih/color"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/johnmerga/Mastering_Go/snippetbox/internal/models"
)

type application struct {
	infoLog       *log.Logger
	errLog        *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	//cmd example: go run ./ -port=":5000"
	port := flag.String("port", ":4000", "your custom port that belongs between (3000-9000)")
	dsn := flag.String("dsn", "web:password@tcp(127.0.0.1:3306)/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoColor := color.New(color.FgCyan).Add(color.Bold)
	errorColor := color.New(color.FgRed).Add(color.Bold)
	infoLog := log.New(os.Stdout, infoColor.Sprint("INFO:\t"), log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, errorColor.Sprint("ERROR\t"), log.Ldate|log.Ltime|log.Lshortfile)

	dbPool, err := openDb(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Println("Database connection successfully opened")
	defer dbPool.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	app := &application{
		infoLog: infoLog,
		errLog:  errorLog,
		snippets: &models.SnippetModel{
			DB: dbPool,
		},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *port,
		ErrorLog: app.errLog,
		Handler:  app.routes(),
	}
	pColor := color.New(color.FgGreen).Add(color.ResetItalic)
	infoLog.Printf("- 🚀starting server on PORT%s", pColor.Sprint(*port))
	if err := srv.ListenAndServe(); err != nil {
		errorLog.Fatal(err)
	}

}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// create a connection and check for any errors
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
