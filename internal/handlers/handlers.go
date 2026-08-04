package handlers

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/osugodbless/groupie-tracker/internal/config"
)

type Application struct {
	Templates *template.Template
}

func renderTemplate(w http.ResponseWriter, app *Application, contentFile string, data any) {
	buf := new(bytes.Buffer)
	err := app.Templates.ExecuteTemplate(buf, contentFile, data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	renderTemplate(w, app, "base.tmpl", config.ArtistByID)
}

func (app *Application) ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	renderTemplate(w, app, "artists-page", config.ArtistByID)
}

func (app *Application) GetArtistHandler(w http.ResponseWriter, r *http.Request) {
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

	renderTemplate(w, app, "artistsDetails.tmpl", artist)
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

	renderTemplate(w, app, "tour-dates.tmpl", artist)
}

func (app *Application) SearchArtists(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.FormValue("query")

	var matches map[int]config.Artist

	if query == "" {
		matches = config.ArtistByID
	} else {
		matches = make(map[int]config.Artist)
		for id, artist := range config.ArtistByID {
			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
				matches[id] = artist
			}
		}
	}
	renderTemplate(w, app, "artists:grid", matches)
}

func getArtistByID(id int) (config.Artist, error) {
	artist, ok := config.ArtistByID[id]
	if ok {
		return artist, nil
	}
	return config.Artist{}, errors.New("Artist not found")
}
