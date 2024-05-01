package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// Define an applicaton structu to hold the application-wide dependencies for the
// web application. For now we'll only include the structured looger , but we'll
// add more to this as the build progresses.
type application struct {
	logger *slog.Logger
}

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP networkd address")
	//Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This read id the command-line flag value and assigns it to the addr
	// ortherwise it will always contain the default value of ":4000". if any errors are will be terminated.
	flag.Parse()

	// Use the slog.New() funciton to initialize a new structured logger, which
	// writes to the standard out stream and usese the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}))

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

	// Initialize a new instance of our applicaton struct, containing the
	// dependencies (for now, just the structured logger).
	app := &application{logger: logger}

	// Register the other routes as normal...
	// The "{$}" prevents trailing slash URLs from becoming "catch it all"
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snipppetView) // Add the {id} wildcard segment
	mux.HandleFunc("GET /snippet/create", app.snipppetCreate)
	// Create the new route, which is restricted to POST requests only.
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Print a log message to say that the server is starting.
	// Uset the Infor() method to log the starting server mesaage at Info severity
	// (along with the listen address as an atribute).
	logger.Info("starting server", slog.Any("add", *addr))

	// use the http.ListerAndServe() function to start a new web server.
	// we pass in two parameters:
	// - the TCP network address to listen on (in this case ":4000")
	// - the servemux we just created
	// If http.ListenAndServe() returns an error we use the log.Fatal()
	// function to log the error message and exit. Note that any error returned by
	// http.ListenAndServe() is always non-nil.
	err := http.ListenAndServe(*addr, mux) // We pass the dereferenced addr pointer to the ListenAndServer too.
	// And we also use the Error() method to log any error message rturnd by
	// http.ListenAndServe() at Error severity (with no additional attributes),
	// and then call os.Exit(1) to terminate the application with exit code 1.
	logger.Error(err.Error())
	os.Exit(1)
}
