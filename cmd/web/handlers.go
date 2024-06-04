package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	// err := r.ParseForm()
	// if err != nil {
	// 	app.clientError(w, r, http.StatusBadRequest)
	// 	return
	// }

	// Declare a new emtpy instance of the snippet CreateForm struct.
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	// Because the Validator struct is embeded by the snippetCreateForm struct,
	// we ca nall CehckField() directly on it to execute our validation checks.
	form.CheckField(validatior.NotBlank(form.Title), "title", "This field can not be blank")
	form.CheckField(validatior.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validatior.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validatior.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

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

	// Use the Put() method to add as string value ("Snippet succesfully
	// created!") and the corresponding key ("flash") to the session data.
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	// Use the w.WriteHeader() method to send a 201 status code.
	// w.WriteHeader(http.StatusCreated)

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
