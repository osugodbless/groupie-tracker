package config

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
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

var ArtistByID map[int]Artist

func LoadConfig() {
	var artists []Artist
	var relations RelationIndex

	// Get requests to external api concurrently
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		LoadConfigHelper("https://groupietrackers.herokuapp.com/api/artists", &artists)
	}()

	go func() {
		defer wg.Done()
		LoadConfigHelper("https://groupietrackers.herokuapp.com/api/relation", &relations)
	}()

	wg.Wait()

	// Extract artist concert info for easy merging with band personal info
	relationMap := make(map[int]map[string][]string)
	for _, rel := range relations.Index {
		relationMap[rel.ID] = rel.DatesLocation
	}

	// Merge artist information with their concert dates together
	ArtistByID = make(map[int]Artist, len(artists))

	for _, art := range artists {
		art.DatesLocation = relationMap[art.ID]
		ArtistByID[art.ID] = art
	}
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
