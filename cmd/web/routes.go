package main

import (
	"net/http"

	"github.com/justinas/alice"
)

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

	// Create a new middleware chain containingh the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Update these routess to use the new dynamic middleware cahain followed by
	// the appropriate handler funciton. Note that because the alice ThenFunc()
	// method returns a http.Handler (rather than a http.HandlerFunc) we also
	// need to switc to registering the route using the mux.Handle() method.
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snipppetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snipppetCreate))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// Pas the servemux as the 'next' parameter to the commonHeaders middleware.routes
	// Because commonHeader is just a function, and the function returns a
	// http.Handler we don't need to do anything else.'
	// return app.recoverPanic(app.logRequest(commonHeaders(mux)))

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
