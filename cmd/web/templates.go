package main

import "github.com/raiesbo/snippetbox/internal/models"

// Define a templateData type to act as the holding sttructure for
// any dynamic data that we want to pass to our HTML templates.
// At th emoment it only contains one field, but we'll add more
// to it as the build progresses.
type templateData struct {
	Snippet models.Snippet
}
