package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"binai.net/internal/models"
	"binai.net/ui"
)

type templateData struct {
	CurrentYear     int
	Lot             *models.Lot
	Lots            []models.Lot
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
	// added
	Metadata models.Metadata
	Userdata struct {
		Name  string
		Email string
		Lots  []string
	}
	Title     string
	Regions   []string
	CountLots int
	PageSize  []int
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

func seq(from, to int) []int {
	var nums []int
	for i := from; i <= to; i++ {
		nums = append(nums, i)
	}
	return nums
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"seq":       seq,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
