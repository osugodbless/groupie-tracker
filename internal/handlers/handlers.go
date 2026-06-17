package handlers

import (
	"html/template"
	"net/http"

	"github.com/osugodbless/groupie-tracker/internal/config"
)

type Application struct {
	Templates *template.Template
}

func renderTemplate(w http.ResponseWriter, tmpl *Application, data []config.Artist) {
	err := tmpl.Templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	renderTemplate(w, app, nil)
}
