package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// Define a home handler funciton which writes a byte slice containing "hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	// User the Header().Add() method to add a "Server: Go" header to the
	// response header map. Thr first paramter is the header name, and
	// the second parameter is the header value.
	w.Header().Add("Server", "Go")

	w.Write([]byte("hello from Snippetbox"))
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
