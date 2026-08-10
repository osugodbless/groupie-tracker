package main

import (
	"context"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/osugodbless/groupie-tracker/internal/client"
	"github.com/osugodbless/groupie-tracker/internal/handlers"
	"github.com/osugodbless/groupie-tracker/internal/routes"
)

var transport = &http.Transport{
	ResponseHeaderTimeout: 5 * time.Second, // Time to wait for server's first response header
	ExpectContinueTimeout: 1 * time.Second, // Time to wait for a response after sending an `Expect: 100-continue` header
}

var funcMap = template.FuncMap{
	"add": func(x, y int) int {
		return x + y
	},
	"clean": func(s string) string {
		result := strings.ReplaceAll(s, "-", ", ")
		return strings.ReplaceAll(result, "_", " ")
	},
}

var baseTmpl = template.Must(template.New("base").Funcs(funcMap).ParseFiles("templates/base.tmpl", "templates/index.tmpl", "templates/artistsDetails.tmpl", "templates/tour-dates.tmpl", "templates/artists.tmpl"))

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	c := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second, // Overall timeout for the entire request
	}
	client := client.NewAPIClient("https://groupietrackers.herokuapp.com/api", c)

	// Create a context with a timeout to ensure the entire operation doesn't exceed a certain duration
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	logger.Info("Loading band data from external API...")

	// Fetch artists and relations concurrently
	artists, err := client.FetchArtistsAndRelations(ctx)
	if err != nil {
		logger.Error("Failed to load band data", "error", err)
		os.Exit(1)
	}

	logger.Info("Successfully loaded band data", "count", len(artists))

	// Initialize the BandArtistService with the loaded artists
	service := handlers.NewBandArtistService(artists)
	app := handlers.NewApplication(baseTmpl, logger, service)

	mux := routes.Routes(app)

	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}

	logger.Info("Server started", "Address", server.Addr)
	log.Fatal(server.ListenAndServe())
}
