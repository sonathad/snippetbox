package main

import (
	"flag"
	"log"
	"net/http"
	"os"
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

	app.infoLog.Printf("Starting server on %s", cfg.addr)
	err := srv.ListenAndServe()
	app.errorLog.Fatal(err)
}
