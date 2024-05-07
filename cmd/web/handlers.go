package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/raiesbo/snippetbox/internal/models"
)

// Define a home handler funciton which writes a byte slice containing "hello from Snippetbox" as the response body.
// Change the signature of the home handler so it is defined as method against *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// User the Header().Add() method to add a "Server: Go" header to the
	// response header map. Thr first paramter is the header name, and
	// the second parameter is the header value.
	w.Header().Add("Server", "Go")

	snippets, err := app.snippets.Lastest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// Initialize a slice contaiing the paths to the two files. It's important
	// to note taht the file containing our base template must be the FIRST
	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/pages/home.tmpl",
	// 	"./ui/html/partials/nav.tmpl", // Include the navigation partial in the template files.
	// }

	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message, use
	// The http.Error() fucntio to send an  Internal Server Error response to the
	// user, and then return from the handler so no subsequent code is executed.

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a tesmplate set. Notice that we use ... to pass the contents
	// of the files slice as variadic arguments.
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	// Because the home handler isnow a method agaist the application
	// 	// struct it can access its fields, including the structured logger. We'll
	// 	// use this to create a log entry at Error level cotaining the error
	// 	// message, also including the request method and URI as attributes to
	// 	// assist with debugging.
	// 	// // app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
	// 	// // http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	app.serverError(w, r, err) // Refactor previous code with serverError helper method.
	// 	return
	// }

	// Then we use the Execute() method on the template set to write the
	// template content as the response body. The last parameter to Execute()
	// represents any dynamic data that we want to pass in, which for now we'll
	// leave as nil.

	// Uset the ExecuteTemplate() method to write the content of the "base"
	// template as the reponse body.
	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	// And we also need to update the code here to use the structured logger too.
	// 	// // app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
	// 	// // http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	app.serverError(w, r, err) // Refactor previous code with serverError helper method.
	// }

	// w.Write([]byte("hello from Snippetbox"))
}

// Add a snippetView handler function.
// Change the signature of the snippetView handler so it is defined as a method
// against *application
func (app *application) snipppetView(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id wildcard from the request using r.PathValue()
	// it can't be converted to an integer, or the value is less thatn 1, we
	// return a 404 page not found reponse
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	//Use the SnippetModel's Get() method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Write the snippet data as a plain text HTTP response body.
	fmt.Fprintf(w, "%+v", snippet)

	// Use the fmt.Sprintf() function to interpolate the id value with a message,
	// then write it as the HTTP response.
	// fmt.Fprintf(w, "Display a specific snippet with ID %d, %s...", id, snippet.Title) // Fprintf write the formatted string to "w" ResponseWritter.
}

// Add a snippetCreate handler function.
func (app *application) snipppetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

// Add a snippetCreatePost handler function.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	// Pass the data to the SnippetModdel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Use the w.WriteHeader() method to send a 201 status code.
	// w.WriteHeader(http.StatusCreated)

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	// w.Write([]byte("Save a new snippet..."))
}
