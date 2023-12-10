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
type loginForm struct {
	Email              string `form:"email"`
	Password           string `form:"password"`
	valiator.Validator `form:"-"`
}

type passwordForm struct {
	OldPassword        string `form:"oldPassword"`
	NewPassword        string `form:"newPassword"`
	ConfirmPassword    string `form:"confirmPassword"`
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
	form.CheckField(valiator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

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

// post signUp
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
		app.render(w, http.StatusUnprocessableEntity, "signUp.tmpl.html", data)
		return
	}
	err = app.users.Create(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signUp.tmpl.html", data)
			return
		}
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "user successfully created. please login")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = loginForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", data)
	return
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form loginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(valiator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(valiator.Matches(form.Email, valiator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(valiator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusBadRequest, "login.tmpl.html", data)
		return
	}
	id, err := app.users.AuthenticateUser(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("incorect email or password")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
			return
		}
		app.serverError(w, err)
		return
	}
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Login successful")
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	redirectPath := app.sessionManager.PopString(r.Context(), "redirectPathAfterLogin")
	if redirectPath != "" {
		app.sessionManager.Remove(r.Context(), "redirectPathAfterLogin")
		http.Redirect(w, r, redirectPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "about.tmpl.html", data)
}

// account

func (app *application) account(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	user, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.sessionManager.Put(r.Context(), "flash", "You must be logged in to access this page")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.User = user
	app.render(w, http.StatusOK, "account.tmpl.html", data)
}

func (app *application) accountPasswordUpdateForm(w http.ResponseWriter, r *http.Request) {
	form := passwordForm{}

	data := app.newTemplateData(r)
	data.Form = form
	app.render(w, http.StatusOK, "password.tmpl.html", data)
}

func (app *application) accountPasswordUpdatePost(w http.ResponseWriter, r *http.Request) {
	var form passwordForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(valiator.NotBlank(form.OldPassword), "oldPassword", "This field cannot be blank")
	form.CheckField(valiator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(valiator.MinChars(form.NewPassword, 10), "newPassword", "This field must be at least 10 characters long")
	form.CheckField(valiator.NotBlank(form.ConfirmPassword), "confirmPassword", "This field cannot be blank")
	form.CheckField(valiator.MinChars(form.ConfirmPassword, 10), "confirmPassword", "This field must be at least 10 characters long")
	if form.NewPassword != form.ConfirmPassword {
		form.AddFieldError("confirmPassword", "This field must match the new password")
	}
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "password.tmpl.html", data)
		return
	}
	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	err = app.users.PasswordUpdate(id, form.OldPassword, form.ConfirmPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Incorrect old password")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "password.tmpl.html", data)
			return
		}
		if errors.Is(err, models.ErrNoRecord) {
			app.sessionManager.Put(r.Context(), "flash", "You must be logged in to access this page")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		if errors.Is(err, models.ErrPasswordMustMatch) {
			form.AddFieldError("confirmPassword", "This field must match the new password")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "password.tmpl.html", data)
			return
		}
		if errors.Is(err, models.ErrCanNotUseOldPassword) {
			form.AddNonFieldError(" You can not use your previous password")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "password.tmpl.html", data)
			return
		}
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your password has been updated successfully!")
	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}
