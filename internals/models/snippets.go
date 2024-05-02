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
func (m *SnippetModel) Insert(title string, content, string, expires int) (int, error) {
	return 0, nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Lastest() ([]Snippet, error) {
	return nil, nil
}