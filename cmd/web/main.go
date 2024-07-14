package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	// can I remove this?
	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	addr string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// setting default address if no argument is provided
	// #OLD: addr := flag.String("addr", ":4000", "HTTP network address")
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")

	dsn := flag.String("dsn", "web:kekw123@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Scalable way to add dependencies
	app := &application{
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	}

	// Declaring a new http.Server struct, to leverage errorLog
	// for logging server problems
	// #OLD: err := http.ListenAndServe(cfg.addr, mux)
	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	db, err := openDB(*dsn)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	defer db.Close()

	app.infoLog.Printf("Starting server on %s", cfg.addr)
	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
