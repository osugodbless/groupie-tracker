package client_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/osugodbless/groupie-tracker/internal/client"
)

type RoundTripperFunc func(req *http.Request) *http.Response

func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripperFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}
func TestFetchArtistsAndRelations(t *testing.T) {
	testClient := NewTestClient(func(req *http.Request) *http.Response {
		switch req.URL.String() {
		case "https://example.com/api/artists":
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`[{"id": 1}]`)),
				Header:     make(http.Header),
			}
		case "https://example.com/api/relation":
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"index": [{"id": 1, "datesLocations": {"New_York": ["2023-01-01"]}}]}`)),
				Header:     make(http.Header),
			}
		default:
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
			}
		}
	})

	apiClient := client.NewAPIClient("https://example.com/api", testClient)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	artists, err := apiClient.FetchArtistsAndRelations(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(artists) == 0 {
		t.Errorf("Expected artists to be populated, but the map was empty")
	}

	for id, artist := range artists {
		if id != 1 {
			t.Errorf("Expected artist ID 1, got %d", id)
		}
		if artist.ID != 1 {
			t.Errorf("Expected artist ID 1, got %d", artist.ID)
		}
	}

}

func TestNewAPIClient(t *testing.T) {
	c := client.NewAPIClient("https://example.com/api", &http.Client{Timeout: 10 * time.Second})

	if c.BaseURL != "https://example.com/api" {
		t.Errorf("Expected BaseURL to be 'https://example.com/api', got '%s'", c.BaseURL)
	}
	if c.Client.Timeout != 10*time.Second {
		t.Errorf("Expected Client timeout to be 10 seconds, got %v", c.Client.Timeout)
	}
}
