package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspondot the fields in our DB snippets table
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a Snippet type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a ewn snippet into the data base.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// Write the SQL statement we want to execute. I've split it over two lines
	// for readability (which is whay it's surrounded with backquotes instad
	// of norma double quotes).
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 DAY' * $3)
	RETURNING id;`

	// Use the Exec() method on the embeded connection pool to execute the
	// statement. The first parameter is the SAL statement, followed by the
	// values for the placeholder parameters: title, content and expiry in
	// that order. This method returns a sqlResult type, which contains some
	// basic informatin about what happened when the statement was executed.
	id := 0
	err := m.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}
	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return id, nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	// Initialize a new zeroed Snippet struct.
	var s Snippet

	// Write the SQL statement we want to execute.
	stmt := `SELECT id, title, content, created, expires
	FROM snippets
	WHERE expires > CURRENT_TIMESTAMP AND id = $1;`

	// Use the QueryRow() method on the connectin pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder paramter. This returns a pointer to a sql.Row object which
	// holds the result from the database.
	// use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet strucut.
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// if the query retusn no ros, then row.Scan() will return a
		// sql.ErrNoRows erro. We use the errors.Is() function check for that
		// error specifically, and return our own ErrNorRecords error instead
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	// If everything went OK, then return the filled Snippet struct.
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Lastest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires
	FROM snippets
	WHERE expires > CURRENT_TIMESTAMP
	ORDER BY id DESC
	LIMIT 10;`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close) to ensure the sql.Rows resultset is
	// always property closed before the Lastest() method returns. This defe
	// statememetn should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil result set.
	defer rows.Close()

	var snippets []Snippet

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		var s Snippet
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created. Again, the argument to row.Scan()
		// must be pointers to the place youwant to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets
		snippets = append(snippets, s)
	}

	// When the rows.Next loop has finished we call rows.Err() to retrieve any
	// error that was entoutered during the iteration. It's important to
	// call this - don't assument that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Snippets slice.
	return snippets, nil
}
