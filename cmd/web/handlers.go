package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/raiesbo/snippetbox/internal/models"
)

// Define a home handler funciton which writes a byte slice containing "hello from Snippetbox" as the response body.
// Change the signature of the home handler so it is defined as method against *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Lastest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateCache(r)
	data.Snippets = snippets

	// use the new render helper.
	app.render(w, r, http.StatusOK, "home.tmpl", data)

	// for _, snippet := range snippets {
	// 	fmt.Fprintf(w, "%+v\n", snippet)
	// }

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

	// data := templateData{
	// 	Snippets: snippets,
	// }

	// Then we use the Execute() method on the template set to write the
	// template content as the response body. The last parameter to Execute()
	// represents any dynamic data that we want to pass in, which for now we'll
	// leave as nil.

	// Uset the ExecuteTemplate() method to write the content of the "base"
	// template as the reponse body.
	// err = ts.ExecuteTemplate(w, "base", data)
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

	data := app.newTemplateCache(r)
	data.Snippet = snippet

	// use the new render helper.
	app.render(w, r, http.StatusOK, "view.tmpl", data)

	// Initialize a slice containing the paths to the view.tmpl file,
	// plus the base layout and navigation partial that we made earlier
	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/view.tmpl",
	// }

	// Parse the templates files...
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// Create an instance of a templateData struct holding the snippet data.
	// data := templateData{
	// 	Snippet: snippet,
	// }

	// And then execute them.
	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// }

	// Use the fmt.Sprintf() function to interpolate the id value with a message,
	// then write it as the HTTP response.
	// fmt.Fprintf(w, "Display a specific snippet with ID %d, %s...", id, snippet.Title) // Fprintf write the formatted string to "w" ResponseWritter.
}

// Add a snippetCreate handler function.
func (app *application) snipppetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateCache(r)
	app.render(w, r, http.StatusOK, "create.tmpl", data)
	w.Write([]byte("Display a form for creating a new snippet..."))
}

// Add a snippetCreatePost handler function.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// First we call the r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same wa for PUT and PATCH
	// requests. If there are any errors,  we use our app.ClientError() helper to
	// send a 400 Bad request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the title and content
	// from the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// The r.PostForm.Get() method always returns the form data as a *string*.
	// However, we're expecting our expires value to be a number, and wan to
	// represent it in our Go code as an integer. So we need to manualy convert
	// the form data to an integer using strconv.Atoi(), and we send 400 Bad
	// Request respones if the conversion fails.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// Initialize a map to hold any validation  errors for the form fields.
	fieldErrors := make(map[string]string)

	// Check taht the title value is not black and is not more than 100
	// charatcters long. If it fails either of those check, ass a message to the
	// errors map using the fiel name as the key.
	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This fields cannot be black"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	// Check that the Content value isn't black.
	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "This field cannot be blank"
	}

	// Check the expires value matches one of the permitted values (1, 7, 365).
	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors["expires"] = "This field must equal 1, 7, 365"
	}

	fmt.Println(fieldErrors)

	// If there are any errors, dump them in a plain HTTP respose and
	// return from the handler.
	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}

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

func (a *application) newTemplateCache(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
	}
}
