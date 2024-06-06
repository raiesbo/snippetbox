package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
	"github.com/raiesbo/snippetbox/internal/models"
)

// Define an applicaton structu to hold the application-wide dependencies for the
// web application. For now we'll only include the structured looger , but we'll
// add more to this as the build progresses.
// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers.
type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a decoder instance
	formDecoder := form.NewDecoder()

	// Use the scs.New() function to initialize a new session manager. Then we
	// configure it to use our Postgres database as the session store, and set a
	// lifetime of 12 hours (so that sessions automattically expire 12 hours
	// affeter first bein created).
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	// Make sure that the Secure attribute is set on our session cookies.
	// Settings this means that the cookie will only be setn by a user's web
	// browser when the HTTPS connection is being used (and won't be sent over an
	// unsecure HTTP connection).
	sessionManager.Cookie.Secure = true

	// And add the session manage to our application dependencies.

	// Initialize a new instance of our applicaton struct, containing the
	// dependencies (for now, just the structured logger).
	// Initialize a models.SnippetModel instance cotaining the connection pool
	// and add it to the application dependencies.
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Initialize a tls.Config struct to hold the non-default TLS setting swe
	// want the server to use. In this case the only thing that we're changing
	// is the curve preferences value, so that on ly elliptic curves with
	// assembly implementations are used.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Initialize a new http.Server struc. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before
	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		// Create a *log.Logger from our structured logger handler, which writes
		// log entries at Error level, and assign it to the ErrorLog field.
		ErrorLog:  slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
		TLSConfig: tlsConfig,
	}

	// Print a log message to say that the server is starting.
	// Uset the Infor() method to log the starting server mesaage at Info severity
	// (along with the listen address as an atribute).
	logger.Info("starting server", slog.Any("add", *addr))

	// Call the ListenAndServe() method on our new http.Server strcut to start
	// the server
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
