package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) render(response http.ResponseWriter, status int, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(response, fmt.Errorf("template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td))
	if err != nil {
		app.serverError(response, err)
		return
	}

	response.WriteHeader(status)
	buf.WriteTo(response)
}

func (app *application) addDefaultData(td *templateData) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CurrentYear = time.Now().Year()

	if td.SnippetForm == nil {
		td.SnippetForm = &SnippetForm{
			Expires: 1,
		}
	}

	return td
}

func (app *application) serverError(response http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errLog.Output(2, trace)

	http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(response http.ResponseWriter, status int) {
	http.Error(response, http.StatusText(status), status)
}

func (app *application) notFound(response http.ResponseWriter) {
	app.clientError(response, http.StatusNotFound)
}
