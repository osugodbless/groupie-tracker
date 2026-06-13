package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var transport = &http.Transport{
	ResponseHeaderTimeout: 5 * time.Second, // Time to wait for server's first response header
	ExpectContinueTimeout: 1 * time.Second, // Time to wait for a response after sending an `Expect: 100-continue` header
}

var client = &http.Client{
	Transport: transport,
	Timeout:   15 * time.Second, // Still advisable to set an overall timeout
}

var appData Data

func LoadConfigHelper(endpointName, endpointUrl string) {
	endpointName = strings.ToLower(endpointName)

	resp, err := client.Get(endpointUrl)
	if err != nil {
		log.Fatalf("Error fetching %v: %v", endpointName, err)
	}
	defer resp.Body.Close()

	switch endpointName {

	// Decode Artists
	case "artists":
		err = json.NewDecoder(resp.Body).Decode(&appData.Artists)

	// Decode Locations
	case "locations":
		var locationIndex LocationIndex
		err = json.NewDecoder(resp.Body).Decode(&locationIndex)
		appData.Locations = locationIndex.Index

	// Decode Dates
	case "dates":
		var dateIndex DateIndex
		err = json.NewDecoder(resp.Body).Decode(&dateIndex)
		appData.Dates = dateIndex.Index

	// Decode Relations
	case "relations":
		var relationIndex RelationIndex
		err = json.NewDecoder(resp.Body).Decode(&relationIndex)
		appData.Relations = relationIndex.Index

	}

	// Check for error
	if err != nil {
		log.Fatalf("Error decoding %v: %v", endpointName, err)
	}

}

func LoadConfig() {
	LoadConfigHelper("artists", "https://groupietrackers.herokuapp.com/api/artists")
	fmt.Printf("Loaded %d Artists\n", len(appData.Artists))

	LoadConfigHelper("locations", "https://groupietrackers.herokuapp.com/api/locations")
	fmt.Printf("Loaded %d Location profiles\n", len(appData.Locations))

	LoadConfigHelper("dates", "https://groupietrackers.herokuapp.com/api/dates")
	fmt.Printf("Loaded %d Date profiles\n", len(appData.Dates))

	LoadConfigHelper("relations", "https://groupietrackers.herokuapp.com/api/relation")
	fmt.Printf("Loaded %d Relation profiles\n", len(appData.Relations))
}
