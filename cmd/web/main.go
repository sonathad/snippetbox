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

func main() {
	// setting default address if no argument is provided
	// #OLD: addr := flag.String("addr", ":4000", "HTTP network address")
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.Parse()

	// Decouple logging early
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	// Fileserver used to serve static files
	fileServer := http.FileServer(customFS{http.Dir("./assets/static/")})

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Declaring a new http.Server struct, to leverage errorLog
	// for logging server problems
	// #OLD: err := http.ListenAndServe(cfg.addr, mux)
	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", cfg.addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
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
