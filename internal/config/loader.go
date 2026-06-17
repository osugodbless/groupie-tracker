package config

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var transport = &http.Transport{
	ResponseHeaderTimeout: 5 * time.Second, // Time to wait for server's first response header
	ExpectContinueTimeout: 1 * time.Second, // Time to wait for a response after sending an `Expect: 100-continue` header
}

var client = &http.Client{
	Transport: transport,
	Timeout:   15 * time.Second, // Overall timeout
}

var CompleteBandData []Artist

func LoadConfig() {
	var artists []Artist
	var relations RelationIndex

	// Get requests to external api
	LoadConfigHelper("https://groupietrackers.herokuapp.com/api/artists", &artists)
	LoadConfigHelper("https://groupietrackers.herokuapp.com/api/relation", &relations)

	// Extract band concert info for easy merging with band personal info
	relationMap := make(map[int]map[string][]string)
	for _, rel := range relations.Index {
		relationMap[rel.ID] = rel.DatesLocation
	}

	// Merge band information with their concert dates together
	CompleteBandData = make([]Artist, len(artists)) // Allocate enough space to hold the complete band data
	for i, art := range artists {
		CompleteBandData[i] = art
		CompleteBandData[i].DatesLocation = relationMap[art.ID]
	}

	// return CompleteBandData
}

func LoadConfigHelper(endpointUrl string, target any) {

	resp, err := client.Get(endpointUrl)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(target)
	// Check for error
	if err != nil {
		log.Fatalf("Error decoding data: %v", err)
	}

}
