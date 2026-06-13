package main

import (
	"fmt"

	"github.com/osugodbless/groupie-tracker/internal/config"
)

func main() {
	config.LoadConfig()
	fmt.Println("Config Data loaded successfully")
}
