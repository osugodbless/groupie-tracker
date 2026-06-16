package main

import (
	"fmt"

	"github.com/osugodbless/groupie-tracker/internal/config"
)

func main() {
	artistsData := config.LoadConfig()
	fmt.Println("Config Data loaded successfully")
	fmt.Println(artistsData)

}
