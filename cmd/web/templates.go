package main

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/IsroilMukhitdinov/snippetbox/internal/models"
)

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	SnippetForm *SnippetForm
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "/pages/*.html"))
	if err != nil {
		return nil, err
	}

	files := []string{
		filepath.Join(dir, "/base.html"),
		filepath.Join(dir, "/partials/nav.html"),
		filepath.Join(dir, "/partials/footer.html"),
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
