package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/IsroilMukhitdinov/snippetbox/internal/models"
	"github.com/IsroilMukhitdinov/snippetbox/internal/validators"
	"github.com/julienschmidt/httprouter"
)

type SnippetForm struct {
	Title   string
	Content string
	Expires int
	validators.Validator
}

func (app *application) home(response http.ResponseWriter, request *http.Request) {

	snippets, err := app.snippetModel.Latest()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			snippets = []*models.Snippet{}
		} else {
			app.serverError(response, err)
			return
		}
	}

	app.render(response, http.StatusOK, "home.html", &templateData{
		Snippets: snippets,
	})
}

func (app *application) snippetView(response http.ResponseWriter, request *http.Request) {

	params := httprouter.ParamsFromContext(request.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(response)
		return
	}

	snippet, err := app.snippetModel.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(response)
			return
		} else {
			app.serverError(response, err)
			return
		}
	}

	app.render(response, http.StatusOK, "view.html", &templateData{
		Snippet: snippet,
	})
}

func (app *application) snippetCreatePost(response http.ResponseWriter, request *http.Request) {

	err := request.ParseForm()
	if err != nil {
		app.clientError(response, http.StatusBadRequest)
		return
	}

	title := request.PostForm.Get("title")
	content := request.PostForm.Get("content")
	expires, err := strconv.Atoi(request.PostForm.Get("expires"))
	if err != nil {
		app.clientError(response, http.StatusBadRequest)
		return
	}

	form := &SnippetForm{
		Title:   title,
		Content: content,
		Expires: expires,
	}

	form.CheckField(validators.NotBlank(title), "title", "This field cannot be blank")
	form.CheckField(validators.NotBlank(content), "content", "This field cannot be blank")
	form.CheckField(validators.MaxLength(title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validators.PermittedValues(expires, 1, 30, 365), "expires", "This field must equal 1, 30 or 365")

	if !form.Valid() {
		app.render(response, http.StatusUnprocessableEntity, "create.html", &templateData{
			SnippetForm: form,
		})

		return
	}

	id, err := app.snippetModel.Insert(title, content, expires)
	if err != nil {
		app.serverError(response, err)
		return
	}

	http.Redirect(response, request, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetCreate(response http.ResponseWriter, request *http.Request) {
	app.render(response, http.StatusOK, "create.html", &templateData{})
}
