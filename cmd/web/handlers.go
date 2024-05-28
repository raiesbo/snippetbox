package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/raiesbo/snippetbox/internal/models"
	"github.com/raiesbo/snippetbox/internal/validatior"
)

// Update our snippetCreateFomr struct to include struct tags which tell the
// decoder how to map HTML form values into the differen tstruct fields. So, for
// example, here we'are telling the decoder to store the value from th HTML form
// input with the name "title" in th eTitle field. The struct tag `form:"-"`
// tells the decoder to completely ingnore a field during decoding.
type snippetCreateForm struct {
	Title                string `form:"title"`
	Content              string `form:"content"`
	Expires              int    `form:"expires"`
	validatior.Validator `form:"_"`
}

// Define a snippetCreatForm struct to represent the form adata and validation
// errors for the form fields. Note that all the struct fields are deliberately
// exported (i.e. start with a capital letter). This is because structu fields
// must be exported in order to be read by the html/template package when
// rendering the template.
// type snippetCreateForm struct {
// 	Title   string
// 	Content string
// 	Expires int
// 	validatior.Validator
// }

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
}

// Add a snippetCreate handler function.
func (app *application) snipppetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateCache(r)

	// Initialized a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a greate opportunity to set any default or
	// 'initial' values fo rthe form --- here we set the initial value for the
	// snippet expiry to 365 days.
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, r, http.StatusOK, "create.tmpl", data)
	// w.Write([]byte("Display a form for creating a new snippet..."))
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

	// Declare a new emtpy instance of the snippet CreateForm struct.
	var form snippetCreateForm

	// Call the Decode() method of the form decoder, passing in the current
	// request and *a pointer* to our snippetCreateForm struct. This will
	// exxentially fill our struct with the relevant values from the HTML form.
	// If there is a problem, we return a 400 Bad Request response to the client.
	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// Then validate an duse the data as normal...

	// Use the r.PostForm.Get() method to retrieve the title and content
	// from the r.PostForm map.
	// title := r.PostForm.Get("title")
	// content := r.PostForm.Get("content")

	// The r.PostForm.Get() method always returns the form data as a *string*.
	// However, we're expecting our expires value to be a number, and wan to
	// represent it in our Go code as an integer. So we need to manualy convert
	// the form data to an integer using strconv.Atoi(), and we send 400 Bad
	// Request respones if the conversion fails.
	// expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	// if err != nil {
	// 	app.clientError(w, r, http.StatusBadRequest)
	// 	return
	// }

	// Create a instance of the snippetCreateForm struct containing the values
	// from the form and an empty map for any validation errors.
	// form := snippetCreateForm{
	// 	Title:   r.PostForm.Get("title"),
	// 	Content: r.PostForm.Get("content"),
	// 	Expires: expires,
	// }

	// Because the Validator struct is embeded by the snippetCreateForm struct,
	// we ca nall CehckField() directly on it to execute our validation checks.
	form.CheckField(validatior.NotBlank(form.Title), "title", "This field can not be blank")
	form.CheckField(validatior.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validatior.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validatior.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// Initialize a map to hold any validation  errors for the form fields.
	// fieldErrors := make(map[string]string)

	// Check taht the title value is not black and is not more than 100
	// charatcters long. If it fails either of those check, ass a message to the
	// errors map using the fiel name as the key.
	// if strings.TrimSpace(form.Title) == "" {
	// 	form.FieldErrors["title"] = "This fields cannot be black"
	// } else if utf8.RuneCountInString(form.Title) > 100 {
	// 	form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	// }

	// Check that the Content value isn't black.
	// if strings.TrimSpace(form.Content) == "" {
	// 	form.FieldErrors["content"] = "This field cannot be blank"
	// }

	// Check the expires value matches one of the permitted values (1, 7, 365).
	// if expires != 1 && expires != 7 && expires != 365 {
	// 	form.FieldErrors["expires"] = "This field must equal 1, 7, 365"
	// }

	// If there are any validation errors, then redisplay the craete.tmpl template,
	// passing in the snippetCreateForm instance as dynamic data in the Form
	// field. Note that we use the HTTP status code 422 Unprocessable Entity
	// when sendin the response to indicate that tehre was a validation error.
	if !form.Valid() {
		data := app.newTemplateCache(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	// Pass the data to the SnippetModdel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
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
