package routes

import (
	"net/http"

	"github.com/osugodbless/groupie-tracker/internal/handlers"
)

func Routes(app *handlers.Application) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.HomeHandler)
	mux.HandleFunc("GET /artists", app.ArtistsHandler)
	mux.HandleFunc("GET /artists/{id}", app.GetArtistHandler)
	mux.HandleFunc("GET /artists/{id}/tour-data", app.TourDatesHandler)
	mux.HandleFunc("POST /artists/search", app.SearchArtists)
	return mux
}
