package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Define a home handler function, which writes a byte slice
// containing "Hello from Snippetbox" as the response body
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	templatePaths := []string{
		"./assets/html/pages/index.html",
		"./assets/html/base.html",
		"./assets/html/partials/nav.html",
	}

	ts, err := template.ParseFiles(templatePaths...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// writes the content of the "base" template as the response
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		// #OLD: app.errorLog.Print(err.Error())
		// http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		app.serverError(w, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// adds an 'Allow: POST' header to show which method is allowed in the endpoint
		w.Header().Set("Allow", http.MethodPost)

		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed"))
		// calls WriteHeader, is equivalent to the above
		// http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet"))
}
