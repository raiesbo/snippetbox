package models

import (
	"database/sql"
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
func (m *SnippetModel) Insert(title string, content string) (int, error) {
	// Write the SQL statement we want to execute. I've split it over two lines
	// for readability (which is whay it's surrounded with backquotes instad
	// of norma double quotes).
	query := `INSERT INTO snippets (title, content, created, expires) 
	VALUES($1, $2, timezone('utc', now()), timezone('utc', now()))
	RETURNING id`

	// Use the Exec() method on the embeded connection pool to execute the
	// statement. The first parameter is the SAL statement, followed by the
	// values for the placeholder parameters: title, content and expiry in
	// that order. This method returns a sqlResult type, which contains some
	// basic informatin about what happened when the statement was executed.
	id := 0
	m.DB.QueryRow(query, title, content).Scan(&id)

	// Use the LastInsertid() method on the result to get the ID of our
	// newly inserted record in the snippets table.
	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return 0, err
	// }

	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Lastest() ([]Snippet, error) {
	return nil, nil
}
