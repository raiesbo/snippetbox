package main

import (
	"log"
	"net/http"
)

func main() {
	// use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()

	// Create a file server shich serves files out of the "./ui/static" directory.
	// Notes that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Register the other routes as normal...
	// The "{$}" prevents trailing slash URLs from becoming "catch it all"
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snipppetView) // Add the {id} wildcard segment
	mux.HandleFunc("GET /snippet/create", snipppetCreate)
	// Create the new route, which is restricted to POST requests only.
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// Print a log message to say that the server is starting.
	log.Print("staging server on :4000")

	// use the http.ListerAndServe() function to start a new web server.
	// we pass in two parameters:
	// - the TCP network address to listen on (in this case ":4000")
	// - the servemux we just created
	// If http.ListenAndServe() returns an error we use the log.Fatal()
	// function to log the error message and exit. Note that any error returned by
	// http.ListenAndServe() is always non-nil.
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
