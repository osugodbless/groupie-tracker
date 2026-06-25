package handlers

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/osugodbless/groupie-tracker/internal/config"
)

type Application struct {
	Template *template.Template
}

func renderTemplate(w http.ResponseWriter, app *Application, tmpl string, data any) {
	err := app.Template.ExecuteTemplate(w, tmpl+".html", data)
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
	renderTemplate(w, app, "index", config.ArtistByID)
}

func (app *Application) ArtistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	artist, err := getArtistByID(id)

	if err != nil {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	renderTemplate(w, app, "artistsDetails", artist)
}

func (app *Application) TourDatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	artist, err := getArtistByID(id)

	if err != nil {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	renderTemplate(w, app, "tour-dates", artist)
}

func getArtistByID(id int) (config.Artist, error) {
	artist, ok := config.ArtistByID[id]
	if ok {
		return artist, nil
	}
	return config.Artist{}, errors.New("Artist not found")
}
