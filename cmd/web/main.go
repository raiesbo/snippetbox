package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	// Import the models package that we just creates.
	"github.com/raiesbo/snippetbox/internals/models"
)

// Define an applicaton structu to hold the application-wide dependencies for the
// web application. For now we'll only include the structured looger , but we'll
// add more to this as the build progresses.
// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers.
type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP networkd address")

	dsn := flag.String("dsn", "postgres://root:root@127.0.0.1:5432/snippetbox?sslmode=disable", "Postges data source name")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This read id the command-line flag value and assigns it to the addr
	// ortherwise it will always contain the default value of ":4000". if any errors are will be terminated.
	flag.Parse()

	// Use the slog.New() funciton to initialize a new structured logger, which
	// writes to the standard out stream and usese the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}))

	// To keep the main() function tidy I'ave put the code for creating a connection
	// pool into the separate openDB() function below. We pass OpenDB() the DSN
	// from the command-lint flag.
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialize a new instance of our applicaton struct, containing the
	// dependencies (for now, just the structured logger).
	// Initialize a models.SnippetModel instance cotaining the connection pool
	// and add it to the application dependencies.
	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

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
	// Call the enw app.routes() method to get the servemux containing our routes,
	// and pass that to http.ListenAndServe().
	err = http.ListenAndServe(*addr, app.routes()) // We pass the dereferenced addr pointer to the ListenAndServer too.
	// And we also use the Error() method to log any error message rturnd by
	// http.ListenAndServe() at Error severity (with no additional attributes),
	// and then call os.Exit(1) to terminate the application with exit code 1.
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
