package routes

import (
	"net/http"

	"github.com/osugodbless/groupie-tracker/assets"
	"github.com/osugodbless/groupie-tracker/internal/handlers"
)

// Routes sets up the routes for the application and returns an http.Handler.
func Routes(app *handlers.Application) http.Handler {
	fs := http.FileServerFS(assets.StaticFS)
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static", fs))
	mux.HandleFunc("GET /{$}", app.HomeHandler)
	mux.HandleFunc("GET /artists", app.ArtistsHandler)
	mux.HandleFunc("GET /artists/filter", app.FilterArtists)
	mux.HandleFunc("GET /artists/{id}", app.GetArtistHandler)
	mux.HandleFunc("GET /artists/{id}/tour-data", app.TourDatesHandler)

	return mux
}
