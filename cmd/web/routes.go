package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/raiesbo/snippetbox/ui"
)

// The routes9) method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	// use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()

	// Use the http.FileServerFS() function to create a HTTP handler which
	// serves the embedded files in ui.Files. It's important to note that our
	// static files are contained in the "static" folder of the ui.Files
	// embedded filesystem. So, for example, our CSS stylesheet is located at
	// "static/css/main.css". This means that we no longer need to strip the
	// prefix from the request URL -- any requests that start with /static/ can
	// just be passed directly to the file server and the corresponding static
	// file will be served (so long as it exists).
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// Add a new GET /ping route
	mux.HandleFunc("GET /ping", ping)

	// Create a file server shich serves files out of the "./ui/static" directory.
	// Notes that the path given to the http.Dir function is relative to the project
	// directory root.
	// fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Create a new middleware chain containingh the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// Update these routess to use the new dynamic middleware cahain followed by
	// the appropriate handler funciton. Note that because the alice ThenFunc()
	// method returns a http.Handler (rather than a http.HandlerFunc) we also
	// need to switc to registering the route using the mux.Handle() method.
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snipppetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected (authenticated-only) application routes, using a new "protected"
	// middleware chain which includes the requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snipppetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// Pas the servemux as the 'next' parameter to the commonHeaders middleware.routes
	// Because commonHeader is just a function, and the function returns a
	// http.Handler we don't need to do anything else.'
	// return app.recoverPanic(app.logRequest(commonHeaders(mux)))

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
