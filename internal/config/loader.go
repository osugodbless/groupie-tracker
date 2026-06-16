package config

import (
	"encoding/json"
	"fmt"
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
	Timeout:   15 * time.Second, // Still advisable to set an overall timeout
}

var CompleteArtistsData []Artist

func LoadConfig() []Artist {
	var artists []Artist
	var relations RelationIndex

	LoadConfigHelper("https://groupietrackers.herokuapp.com/api/artists", &artists)
	fmt.Printf("Loaded %d Artists\n", len(artists))

	LoadConfigHelper("https://groupietrackers.herokuapp.com/api/relation", &relations)
	fmt.Printf("Loaded %d Relation profiles\n", len(relations.Index))

	relationMap := make(map[int]map[string][]string)
	for _, rel := range relations.Index {
		relationMap[rel.ID] = rel.DatesLocation
	}

	CompleteArtistsData = make([]Artist, len(artists))
	for i, art := range artists {
		CompleteArtistsData[i] = art
		CompleteArtistsData[i].DatesLocation = relationMap[art.ID]
	}

	return CompleteArtistsData
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
