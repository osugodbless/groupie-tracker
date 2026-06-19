package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/osugodbless/groupie-tracker/internal/config"
	"github.com/osugodbless/groupie-tracker/internal/handlers"
	"github.com/osugodbless/groupie-tracker/internal/routes"
)

var funcMap = template.FuncMap{
	"add": func(x, y int) int {
		return x + y
	},
}

var tmpl = template.Must(template.New("main").Funcs(funcMap).ParseFiles("templates/index.html", "templates/artistsDetails.html", "templates/tour-dates.html"))

func main() {

	app := &handlers.Application{
		Template: tmpl,
	}

	config.LoadConfig()

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
