package main

import (
	"html/template"
	"path/filepath"

	"github.com/raiesbo/snippetbox/internal/models"
)

// Define a templateData type to act as the holding sttructure for
// any dynamic data that we want to pass to our HTML templates.
// At th emoment it only contains one field, but we'll add more
// to it as the build progresses.
type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepaht.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/html/page/*.tmpl"
	// us a slice of all the filepaths for our application 'pate' templates
	// like [ui/html/page/home.tmpl]
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Loop through the page filepaths one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'home.tmpl') from the full filepath
		// and assign it to the name variable
		name := filepath.Base(page)

		// Parse the base template file into a template set,

		ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Create a slice containing the filepaths for our base template, any
		// parts and the page
		// files := []string{
		// 	"./ui/html/base.tmpl",
		// 	"./ui/html/partials/nav.tmpl",
		// 	page,
		// }

		// // Parse the files inot a template set.
		// ts, err := template.ParseFiles(files...)
		// if err != nil {
		// 	return nil, err
		// }

		// Add the template set to the map, using the name of the page
		// (like 'home.tmpl') as the key.
		cache[name] = ts
	}

	return cache, nil
}
