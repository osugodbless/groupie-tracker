package routes

import (
	"net/http"

	"github.com/osugodbless/groupie-tracker/internal/handlers"
)

func Routes(app *handlers.Application) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.HomeHandler)
	return mux
}
