package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// 2 is the call depth. 2 means the line that called the function that called Output. if we didn't use 2, the line that called Output would be the line that would be printed in the log.
	app.errLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("couldn't find %s in template cache", page)
		app.serverError(w, err)
		return
	}
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
	w.WriteHeader(status)
	buf.WriteTo(w)

}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func formValidator(r http.Request) (*snippetCreateForm, error) {
	if err := r.ParseForm(); err != nil {
		return &snippetCreateForm{}, err
	}
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	errFields := make(map[string]string)
	if err != nil {
		errFields["expires"] = "Please select a valid expiry time"
	}
	if strings.TrimSpace(title) == "" {
		errFields["title"] = "Title field cannot be empty"
	} else if utf8.RuneCountInString(title) > 100 {
		errFields["title"] = "Title field cannot be more than 100 characters"
	}
	if strings.TrimSpace(content) == "" {
		errFields["content"] = "Content field cannot be empty"
	}
	if !(expires == 1 || expires == 7 || expires == 365) {
		errFields["expires"] = "Please select a valid expiry time"
	}
	form := snippetCreateForm{
		Title:       title,
		Content:     content,
		Expires:     expires,
		FieldErrors: errFields,
	}
	if len(errFields) > 0 {
		return &form, nil
	}
	data := snippetCreateForm{
		Title:   title,
		Content: content,
		Expires: expires,
	}
	return &data, nil
}
