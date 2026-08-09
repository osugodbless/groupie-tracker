package handlers

import (
	"bytes"
	"cmp"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/osugodbless/groupie-tracker/internal/config"
)

// ArtistService defines the data access contract for easier mocking and testing.
type ArtistService interface {
	GetByID(id int) (config.Artist, error)
	GetAll() map[int]config.Artist
	Filter(filters Filters) map[int]config.Artist
}

// BandArtistService implements ArtistService safely with a read lock.
type BandArtistService struct {
	mu      sync.RWMutex
	artists map[int]config.Artist
}

func NewBandArtistService(data map[int]config.Artist) *BandArtistService {
	return &BandArtistService{
		artists: data,
	}
}

func (s *BandArtistService) GetByID(id int) (config.Artist, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	artist, ok := s.artists[id]
	if !ok {
		return config.Artist{}, errors.New("artist not found")
	}
	return artist, nil
}

func (s *BandArtistService) GetAll() map[int]config.Artist {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cp := make(map[int]config.Artist, len(s.artists))
	for k, v := range s.artists {
		cp[k] = v
	}
	return cp
}

func (s *BandArtistService) Filter(f Filters) map[int]config.Artist {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[int]config.Artist)

	var fromTime, toTime time.Time
	var errFrom, errTo error
	if f.FirstAlbum.FromDate != "" {
		fromTime, errFrom = time.Parse("2006-01-02", f.FirstAlbum.FromDate)
	}
	if f.FirstAlbum.ToDate != "" {
		toTime, errTo = time.Parse("2006-01-02", f.FirstAlbum.ToDate)
	}

	for id, artist := range s.artists {
		// Search Query
		if f.Search != "" && !strings.Contains(strings.ToLower(artist.Name), strings.ToLower(f.Search)) {
			continue
		}

		// Number of Members
		if f.NumOfMembers > 0 && len(artist.Members) != f.NumOfMembers {
			continue
		}

		// Concert Location
		if f.ConcertLoc != "" {
			rawLoc := strings.ToLower(f.ConcertLoc)
			rawLoc = strings.ReplaceAll(rawLoc, ", ", "-")
			rawLoc = strings.ReplaceAll(rawLoc, " ", "_")
			if _, ok := artist.DatesLocation[rawLoc]; !ok {
				continue
			}
		}

		// First Album Date Range
		if errFrom == nil && !fromTime.IsZero() || errTo == nil && !toTime.IsZero() {
			albumDate, err := time.Parse("02-01-2006", artist.FirstAlbum)
			if err != nil {
				continue
			}
			if !fromTime.IsZero() && albumDate.Before(fromTime) {
				continue
			}
			if !toTime.IsZero() && albumDate.After(toTime) {
				continue
			}
		}

		result[id] = artist
	}

	return result
}

type FirstAlbum struct {
	FromDate string
	ToDate   string
}

type Filters struct {
	Search       string
	Sort         string
	FirstAlbum   FirstAlbum
	NumOfMembers int
	ConcertLoc   string
}

type Application struct {
	Templates     *template.Template
	Logger        *slog.Logger
	ArtistService ArtistService
}

func NewApplication(tmpl *template.Template, logger *slog.Logger, service ArtistService) *Application {
	return &Application{
		Templates:     tmpl,
		Logger:        logger,
		ArtistService: service,
	}
}

func (app *Application) renderTemplate(w http.ResponseWriter, status int, page string, data any) {
	buf := new(bytes.Buffer)
	err := app.Templates.ExecuteTemplate(buf, page, data)
	if err != nil {
		app.Logger.Error("template rendering failed", "template", page, "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := buf.WriteTo(w); err != nil {
		app.Logger.Error("failed writing buffer to response writer", "error", err)
	}
}

// HomeHandler handles GET /
func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, http.StatusOK, "base.tmpl", app.ArtistService.GetAll())
}

// ArtistsHandler handles GET /artists
func (app *Application) ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, http.StatusOK, "artists-page", app.ArtistService.GetAll())
}

// GetArtistHandler handles GET /artists/{id}
func (app *Application) GetArtistHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	artist, err := app.ArtistService.GetByID(id)
	if err != nil {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	app.renderTemplate(w, http.StatusOK, "artistsDetails.tmpl", artist)
}

// TourDatesHandler handles GET /artists/{id}/tour-dates
func (app *Application) TourDatesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	artist, err := app.ArtistService.GetByID(id)
	if err != nil {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	app.renderTemplate(w, http.StatusOK, "tour-dates.tmpl", artist)
}

// FilterArtists handles GET /artists/filter
func (app *Application) FilterArtists(w http.ResponseWriter, r *http.Request) {
	var members int
	if m := r.FormValue("members"); m != "" {
		parsed, err := strconv.Atoi(m)
		if err != nil {
			http.Error(w, "Invalid number of members", http.StatusBadRequest)
			return
		}
		members = parsed
	}

	filters := Filters{
		Search: r.FormValue("search"),
		Sort:   r.FormValue("sort"),
		FirstAlbum: FirstAlbum{
			FromDate: r.FormValue("from_date"),
			ToDate:   r.FormValue("to_date"),
		},
		NumOfMembers: members,
		ConcertLoc:   r.FormValue("concert-loc"),
	}

	filteredMap := app.ArtistService.Filter(filters)

	// Convert to slice for sorting
	artists := make([]config.Artist, 0, len(filteredMap))
	for _, artist := range filteredMap {
		artists = append(artists, artist)
	}

	sortArtists(artists, filters.Sort)

	app.renderTemplate(w, http.StatusOK, "artists:grid", artists)
}

func sortArtists(artists []config.Artist, sortBy string) {
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
	case "creation-asc":
		slices.SortFunc(artists, func(a, b config.Artist) int {
			return cmp.Compare(a.CreationYear, b.CreationYear)
		})
	case "creation-desc":
		slices.SortFunc(artists, func(a, b config.Artist) int {
			return cmp.Compare(b.CreationYear, a.CreationYear)
		})
	}
}
