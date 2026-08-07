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
	"time"

	"github.com/osugodbless/groupie-tracker/internal/config"
)

type Application struct {
	Templates *template.Template
}

type Filters struct {
	Search       string
	Sort         string
	FirstAlbum   FirstAlbum
	NumOfMembers string
	ConcertLoc   string
}

type FirstAlbum struct {
	from_date string
	to_date   string
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
		Search: r.FormValue("search"),
		Sort:   r.FormValue("sort"),
		FirstAlbum: FirstAlbum{
			from_date: r.FormValue("from_date"),
			to_date:   r.FormValue("to_date"),
		},
		NumOfMembers: r.FormValue("members"),
		ConcertLoc:   r.FormValue("concert-loc"),
	}

	if filters.Search != "" {
		matches := searchArtistsByName(filters.Search)
		renderTemplate(w, app, "artists:grid", matches)
		return
	}

	if filters.Sort != "" {
		if filters.Sort == "default" {
			renderTemplate(w, app, "artists:grid", config.ArtistByID)
			return
		}
		artistsToSort := make([]config.Artist, 0, len(config.ArtistByID))
		for _, artist := range config.ArtistByID {
			artistsToSort = append(artistsToSort, artist)
		}
		sortArtists(&artistsToSort, filters.Sort)
		renderTemplate(w, app, "artists:grid", artistsToSort)
		return
	}

	if filters.FirstAlbum.from_date != "" || filters.FirstAlbum.to_date != "" {
		matches := filterArtistsByFirstAlbum(filters.FirstAlbum.from_date, filters.FirstAlbum.to_date)
		renderTemplate(w, app, "artists:grid", matches)
		return
	}

	if filters.NumOfMembers != "" {
		numOfMembers, err := strconv.Atoi(filters.NumOfMembers)
		if err != nil {
			http.Error(w, "Invalid number of members", http.StatusBadRequest)
			return
		}
		matches := filterArtistsByNumOfMembers(numOfMembers)
		renderTemplate(w, app, "artists:grid", matches)
		return
	}

	if filters.ConcertLoc != "" {
		matches := filterArtistsByConcertLocation(filters.ConcertLoc)
		renderTemplate(w, app, "artists:grid", matches)
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

	matches = make(map[int]config.Artist)
	for id, artist := range config.ArtistByID {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
			matches[id] = artist
		}
	}
	return matches
}

func sortArtists(artists *[]config.Artist, sortBy string) {
	sortBy = strings.ToLower(strings.TrimSpace(sortBy))
	// Sort by sorting options
	switch sortBy {
	case "name-asc":
		slices.SortFunc(*artists, func(a, b config.Artist) int {
			return cmp.Compare(a.Name, b.Name)
		})
	case "name-desc":
		slices.SortFunc(*artists, func(a, b config.Artist) int {
			return cmp.Compare(b.Name, a.Name)
		})
	case "creation-asc":
		slices.SortFunc(*artists, func(a, b config.Artist) int {
			return cmp.Compare(a.CreationYear, b.CreationYear)
		})
	case "creation-desc":
		slices.SortFunc(*artists, func(a, b config.Artist) int {
			return cmp.Compare(b.CreationYear, a.CreationYear)
		})
	}
}

func filterArtistsByFirstAlbum(fromDate, toDate string) map[int]config.Artist {
	from, err1 := time.Parse("2006-01-02", fromDate)
	to, err2 := time.Parse("2006-01-02", toDate)

	if err1 != nil || err2 != nil {
		return nil
	}

	matches := make(map[int]config.Artist)

	for id, artist := range config.ArtistByID {
		artistFirstAlbum, err := time.Parse("02-01-2006", artist.FirstAlbum)
		if err != nil {
			continue
		}
		if artistFirstAlbum.After(from) || artistFirstAlbum.Equal(from) {
			if artistFirstAlbum.Before(to) || artistFirstAlbum.Equal(to) {
				matches[id] = artist
			}
		}
	}

	return matches
}

func filterArtistsByNumOfMembers(numOfMembers int) map[int]config.Artist {
	matches := make(map[int]config.Artist)

	for id, artist := range config.ArtistByID {
		if len(artist.Members) == numOfMembers {
			matches[id] = artist
		}
	}

	return matches
}

func filterArtistsByConcertLocation(loc string) map[int]config.Artist {
	rawLoc := strings.ToLower(loc)
	rawLoc = strings.ReplaceAll(loc, ", ", "-")
	rawLoc = strings.ReplaceAll(rawLoc, " ", "_")
	matches := make(map[int]config.Artist)
	for id, artist := range config.ArtistByID {
		_, ok := artist.DatesLocation[rawLoc]
		if ok {
			matches[id] = artist
		}
	}

	return matches
}
