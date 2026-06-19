package handlers

import (
	"html/template"
	"net/http"

	"github.com/osugodbless/groupie-tracker/internal/config"
)

type Application struct {
	Template *template.Template
}

func renderTemplate(w http.ResponseWriter, app *Application, data []config.Artist) {
	err := app.Template.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	renderTemplate(w, app, config.CompleteArtistsData)
}
