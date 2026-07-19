package config

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type APIClient struct {
	Client  *http.Client
	BaseURL string
}

var ArtistByID map[int]Artist

func (api *APIClient) LoadConfig() {
	var artists []Artist
	var relations RelationIndex

	// Get requests to external api concurrently
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		api.LoadConfigHelper(api.BaseURL+"/artists", &artists)
	}()

	go func() {
		defer wg.Done()
		api.LoadConfigHelper(api.BaseURL+"/relation", &relations)
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

func (api *APIClient) LoadConfigHelper(endpointUrl string, target any) {

	resp, err := api.Client.Get(endpointUrl)
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
