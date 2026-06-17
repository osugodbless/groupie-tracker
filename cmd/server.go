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

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func main() {

	app := &handlers.Application{
		Templates: tmpl,
	}

	config.LoadConfig()

	mux := routes.Routes(app)

	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}

	fmt.Printf("Server is running on %s", server.Addr)
	log.Fatal(server.ListenAndServe())

}
