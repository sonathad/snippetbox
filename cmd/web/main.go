package main

import (
	"log"
	"net/http"
	"path/filepath"
)

const serverPort = ":4000"

func main() {
	mux := http.NewServeMux()

	// Fileserver used to serve static files
	fileServer := http.FileServer(customFileSystem{http.Dir("./assets/static/")})

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Printf("Starting server on %s", serverPort)
	err := http.ListenAndServe(serverPort, mux)
	log.Fatal(err)
}

type customFileSystem struct {
	fs http.FileSystem
}

// I define a filesystem with the intention to hide the static file tree navigation.
func (nfs customFileSystem) Open(path string) (http.File, error) {
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
