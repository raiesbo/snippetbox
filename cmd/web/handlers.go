package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Define a home handler funciton which writes a byte slice containing "hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	// User the Header().Add() method to add a "Server: Go" header to the
	// response header map. Thr first paramter is the header name, and
	// the second parameter is the header value.
	w.Header().Add("Server", "Go")

	// Initialize a slice contaiing the paths to the two files. It's important
	// to note taht the file containing our base template must be the FIRST
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl", // Include the navigation partial in the template files.
	}

	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message, use
	// The http.Error() fucntio to send an  Internal Server Error response to the
	// user, and then return from the handler so no subsequent code is executed.

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a tesmplate set. Notice that we use ... to pass the contents
	// of the files slice as variadic arguments.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Then we use the Execute() method on the template set to write the
	// template content as the response body. The last parameter to Execute()
	// represents any dynamic data that we want to pass in, which for now we'll
	// leave as nil.

	// Uset the ExecuteTemplate() method to write the content of the "base"
	// template as the reponse body.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	// w.Write([]byte("hello from Snippetbox"))
}

// Add a snippetView handler function.
func snipppetView(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id wildcard from the request using r.PathValue()
	// it can't be converted to an integer, or the value is less thatn 1, we
	// return a 404 page not found reponse
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Use the fmt.Sprintf() function to interpolate the id value with a message,
	// then write it as the HTTP response.
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id) // Fprintf write the formatted string to "w" ResponseWritter.
}

// Add a snippetCreate handler function.
func snipppetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

// Add a snippetCreatePost handler function.
func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Use the w.WriteHeader() method to send a 201 status code.
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("Save a new snippet..."))
}
