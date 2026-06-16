package main

import (
	"html/template"

	"github.com/osugodbless/groupie-tracker/internal/config"
)

var templates = template.Must(template.ParseFiles("templates/index.html"))

func main() {
	artistsData := config.LoadConfig()

}
