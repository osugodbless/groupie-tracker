package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/osugodbless/groupie-tracker/internal/config"
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

	app := &handlers.Application{
		Templates: baseTmpl,
	}

	apiClient := &config.APIClient{
		Client: &http.Client{
			Transport: transport,
			Timeout:   15 * time.Second, // Still advisable to set an overall timeout
		},
		BaseURL: "https://groupietrackers.herokuapp.com/api",
	}

	apiClient.LoadConfig()

	mux := routes.Routes(app)

	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}

	fmt.Printf("Server started on %v", server.Addr)
	log.Fatal(server.ListenAndServe())
}
