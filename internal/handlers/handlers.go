package handlers

import (
	"bytes"
	"cmp"
	"errors"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/osugodbless/groupie-tracker/internal/config"
)

type Application struct {
	Templates *template.Template
}

type Filters struct {
	Search       string
	Sort         string
	NumOfMembers string
	ConcertLoc   string
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

func (app *Application) FilterArtists(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	filters := Filters{
		Search:       r.FormValue("search"),
		Sort:         r.FormValue("sort"),
		NumOfMembers: r.FormValue("members"),
		ConcertLoc:   r.FormValue("concert-loc"),
	}

	if filters.Search != "" {
		matches := searchArtistsByName(filters.Search)
		renderTemplate(w, app, "artists:grid", matches)
		return
	}

	if filters.Sort != "" {
		artistsToSort := make([]config.Artist, 0, len(config.ArtistByID))
		for _, artist := range config.ArtistByID {
			artistsToSort = append(artistsToSort, artist)
		}
		sortedArtists := sortArtists(artistsToSort, filters.Sort)
		renderTemplate(w, app, "artists:grid", sortedArtists)
		return
	}

	renderTemplate(w, app, "artists:grid", config.ArtistByID)
}

func getArtistByID(id int) (config.Artist, error) {
	artist, ok := config.ArtistByID[id]
	if ok {
		return artist, nil
	}
	return config.Artist{}, errors.New("Artist not found")
}

func searchArtistsByName(query string) map[int]config.Artist {
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

	return matches
}

func sortArtists(artists []config.Artist, sortBy string) []config.Artist {
	sortBy = strings.ToLower(strings.TrimSpace(sortBy))
	switch sortBy {
	case "name-asc":
		slices.SortFunc(artists, func(a, b config.Artist) int {
			return cmp.Compare(a.Name, b.Name)
		})

	case "name-desc":
		slices.SortFunc(artists, func(a, b config.Artist) int {
			return cmp.Compare(b.Name, a.Name)
		})
	case "members":
		// Sorting logic for number of members can be implemented here
	case "creationdate":
		// Sorting logic for creation date can be implemented here
	default:
		log.Printf("Unknown sort option: %s", sortBy)
	}

	return artists

}
