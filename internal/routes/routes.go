package routes

import (
	"net/http"

	"github.com/osugodbless/groupie-tracker/internal/handlers"
)

func Routes(app *handlers.Application) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.HomeHandler)
	mux.HandleFunc("GET /artist/{id}", app.ArtistHandler)
	mux.HandleFunc("GET /artist/{id}/tour-data", app.TourDatesHandler)
	return mux
}
