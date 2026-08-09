package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type APIClient struct {
	client  *http.Client
	baseURL string
}

func NewAPIClient(baseURL string, customClient *http.Client) *APIClient {
	if customClient == nil {
		customClient = &http.Client{
			Timeout: 15 * time.Second,
		}
	}
	return &APIClient{
		client:  customClient,
		baseURL: baseURL,
	}
}

// FetchArtistsAndRelations fetches data concurrently.
func (api *APIClient) FetchArtistsAndRelations(ctx context.Context) (map[int]Artist, error) {
	// Create a child context that can be canceled immediately if any goroutine fails
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		artists   []Artist
		relations RelationIndex
		wg        sync.WaitGroup
		once      sync.Once
	)

	// Buffer channel to hold the first error encountered without blocking
	errChan := make(chan error, 2)

	// Helper function to run concurrent tasks safely
	runTask := func(fetchFunc func() error) {
		defer wg.Done()

		// Short-circuit if context is already canceled by another failed task
		if ctx.Err() != nil {
			return
		}

		if err := fetchFunc(); err != nil {
			// Ensure only the FIRST error triggers cancel() and gets recorded
			once.Do(func() {
				errChan <- err
				cancel() // Cancel context to notify sister goroutines
			})
		}
	}

	wg.Add(2)

	go runTask(func() error {
		return api.fetchJSON(ctx, api.baseURL+"/artists", &artists)
	})

	go runTask(func() error {
		return api.fetchJSON(ctx, api.baseURL+"/relation", &relations)
	})

	// Wait for all goroutines to terminate
	wg.Wait()
	close(errChan)

	// If an error was sent to the channel, return it
	if err := <-errChan; err != nil {
		return nil, fmt.Errorf("failed to load external startup data: %w", err)
	}

	// Merge relationships into artists map
	relationMap := make(map[int]map[string][]string, len(relations.Index))
	for _, rel := range relations.Index {
		relationMap[rel.ID] = rel.DatesLocation
	}

	artistByID := make(map[int]Artist, len(artists))
	for _, art := range artists {
		art.DatesLocation = relationMap[art.ID]
		artistByID[art.ID] = art
	}

	return artistByID, nil
}

func (api *APIClient) fetchJSON(ctx context.Context, endpointURL string, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for %s: %w", endpointURL, err)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed for %s: %w", endpointURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, endpointURL)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode JSON response from %s: %w", endpointURL, err)
	}

	return nil
}
