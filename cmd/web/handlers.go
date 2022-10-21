package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/IsroilMukhitdinov/snippetbox/internal/models"
)

func (app *application) home(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		app.notFound(response)
		return
	}

	snippets, err := app.snippetModel.Latest()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			snippets = []*models.Snippet{}
		} else {
			app.serverError(response, err)
			return
		}
	}

	app.render(response, "home.html", &templateData{
		Snippets: snippets,
	})
}

func (app *application) snippetView(response http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
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

	app.render(response, "view.html", &templateData{
		Snippet: snippet,
	})
}

func (app *application) snippetCreate(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		response.Header().Set("Allow", http.MethodPost)
		app.clientError(response, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := app.snippetModel.Insert(title, content, expires)
	if err != nil {
		app.serverError(response, err)
		return
	}

	http.Redirect(response, request, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
