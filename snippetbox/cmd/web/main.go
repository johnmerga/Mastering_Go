package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	//cmd example: go run ./ -port=":5000"
	port := flag.String("port", ":4000", "your custom port that belongs between (3000-9000)")
	flag.Parse()
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("../../ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux.HandleFunc("/home", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	infoLog.Printf("starting server on port %v", *port)
	if err := http.ListenAndServe(*port, mux); err != nil {
		errorLog.Fatal(err)
	}

}
