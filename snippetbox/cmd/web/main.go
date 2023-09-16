package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	infoLog *log.Logger
	errLog  *log.Logger
}

func main() {
	//cmd example: go run ./ -port=":5000"
	port := flag.String("port", ":4000", "your custom port that belongs between (3000-9000)")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		infoLog: infoLog,
		errLog:  errorLog,
	}

	srv := &http.Server{
		Addr:     *port,
		ErrorLog: app.errLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("starting server on port %v", *port)
	if err := srv.ListenAndServe(); err != nil {
		errorLog.Fatal(err)
	}

}
