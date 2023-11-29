package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/johnmerga/Mastering_Go/snippetbox/internal/models"
	"github.com/johnmerga/Mastering_Go/snippetbox/internal/valiator"
	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title              string `form:"title"`
	Content            string `form:"content"`
	Expires            int    `form:"expires"`
	valiator.Validator `form:"-"`
}

type signupForm struct {
	Name               string `form:"name"`
	Email              string `form:"email"`
	Password           string `form:"password"`
	valiator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "snippet.form.tmpl.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(valiator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(valiator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(valiator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(valiator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "snippet.form.tmpl.html", data)
		return
	}
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "New Snippet created successfully")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// user related
func (app *application) userSignupForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = signupForm{}
	app.render(w, http.StatusOK, "signUp.tmpl.html", data)
	return
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form signupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}
	form.CheckField(valiator.NotBlank(form.Name), "name", "name should not be blank")
	form.CheckField(valiator.MaxChars(form.Name, 50), "name", "name should not be more than 50 characters")
	form.CheckField(valiator.NotBlank(form.Email), "email", "email should not be blank")
	form.CheckField(valiator.MaxChars(form.Email, 80), "email", "email should not be more than 50 characters")
	form.CheckField(valiator.Matches(form.Email, valiator.EmailRX), "email", "email should be valid")
	form.CheckField(valiator.NotBlank(form.Password), "password", "password should not be blank")
	form.CheckField(valiator.MinChars(form.Password, 10), "password", "password should be at least 10 characters")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusBadRequest, "signUp.tmpl.html", data)
		return
	}
	err = app.users.Create(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusBadRequest, "signUp.tmpl.html", data)
			return
		}
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "user successfully created. please login")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLoginForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for logging in a user...")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
