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
	flash := app.sessionManager.PopString(r.Context(), "flash")
	fmt.Printf("this is the session data: %v", flash)
	data := app.newTemplateData(r)
	data.Snippet = snippet
	data.Flash = flash
	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "snippet.form.tmpl.html", data)
}

type snippetCreateForm struct {
	Title              string `form:"title"`
	Content            string `form:"content"`
	Expires            int    `form:"expires"`
	valiator.Validator `form:"-"`
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
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
