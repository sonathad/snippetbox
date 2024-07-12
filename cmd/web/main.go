package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// #TODO move the config out if it grows too big
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

	// Scalable way to add dependencies to the handlers
	app := &application{
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	}

	mux := http.NewServeMux()

	// Fileserver used to serve static files
	fileServer := http.FileServer(customFS{http.Dir("./assets/static/")})

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Declaring a new http.Server struct, to leverage errorLog
	// for logging server problems
	// #OLD: err := http.ListenAndServe(cfg.addr, mux)
	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: app.errorLog,
		Handler:  mux,
	}

	app.infoLog.Printf("Starting server on %s", cfg.addr)
	err := srv.ListenAndServe()
	app.errorLog.Fatal(err)
}

type customFS struct {
	fs http.FileSystem
}

// I define a filesystem with the intention to hide the static file tree navigation.
func (nfs customFS) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	// Returns 404 Not Found for all directories that don't have an index.html
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
