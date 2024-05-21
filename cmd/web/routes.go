package main

import "net/http"

// The routes9) method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
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
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snipppetView) // Add the {id} wildcard segment
	mux.HandleFunc("GET /snippet/create", app.snipppetCreate)
	// Create the new route, which is restricted to POST requests only.
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Pas the servemux as the 'next' parameter to the commonHeaders middleware.routes
	// Because commonHeader is just a function, and the function returns a
	// http.Handler we don't need to do anything else.'
	return app.logRequest(commonHeaders(mux))
}
