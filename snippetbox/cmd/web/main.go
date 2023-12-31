package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/fatih/color"
	"github.com/go-playground/form/v4"

	_ "github.com/go-sql-driver/mysql"
	"github.com/johnmerga/Mastering_Go/snippetbox/internal/models"
)

type application struct {
	infoLog        *log.Logger
	errLog         *log.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	debug          bool
}

func main() {
	//cmd example: go run ./ -port=":5000"
	port := flag.String("port", ":4000", "your custom port that belongs between (3000-9000)")
	dsn := flag.String("dsn", "web:password@tcp(127.0.0.1:3306)/snippetbox?parseTime=true", "MySQL data source name")
	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	infoColor := color.New(color.FgCyan).Add(color.Bold)
	errorColor := color.New(color.FgRed).Add(color.Bold)
	infoLog := log.New(os.Stdout, infoColor.Sprint("INFO:\t"), log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, errorColor.Sprint("ERROR\t"), log.Ldate|log.Ltime|log.Lshortfile)

	dbPool, err := openDb(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Println(color.GreenString("Database connection established successfully"))
	defer dbPool.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	formDecoder := form.NewDecoder()
	// session config
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(dbPool)
	sessionManager.Lifetime = time.Hour * 12
	app := &application{
		infoLog:        infoLog,
		errLog:         errorLog,
		snippets:       &models.SnippetModel{DB: dbPool},
		users:          &models.UserModel{DB: dbPool},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		debug:          *debug,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256}, // add curve preferences as you see fit
		MinVersion:       tls.VersionTLS12,                         // TLS 1.2 is the minimum version we should support
		// Restricting cipher suites
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
	srv := &http.Server{
		Addr:      *port,
		ErrorLog:  app.errLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,      // used to close the connection if the client has been idle for a while
		ReadTimeout:  5 * time.Second,  // used to limit the time the server will wait to read the request body
		WriteTimeout: 10 * time.Second, // used to limit the time the server will wait before writing a response
	}
	pColor := color.New(color.FgGreen).Add(color.ResetItalic)
	infoLog.Printf("- 🚀starting server on PORT%s", pColor.Sprint(*port))
	if err := srv.ListenAndServeTLS("../../tls/cert.pem", "../../tls/key.pem"); err != nil {
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
