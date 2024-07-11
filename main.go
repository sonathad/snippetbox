package main

import (
	"log"
	"net/http"
)

const PORT = ":4000"

// Define a home handler function, which writes a byte slice
// containing "Hello from Snippetbox" as the response body
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// adds an 'Allow: POST' header to show which method is allowed in the endpoint
		w.Header().Set("Allow", http.MethodPost)

		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed"))
		// calls WriteHeader, is equivalent to the above
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet"))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Start a new web server, passing the servemux I made.
	log.Printf("Starting server on %s", PORT)
	err := http.ListenAndServe(PORT, mux)

	// log.Fatal runs os.Exit(1) after logging the error (if any)
	log.Fatal(err)
}
