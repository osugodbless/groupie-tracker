package config_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/osugodbless/groupie-tracker/internal/config"
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
func TestLoadConfig(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
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

	apiClient := &config.APIClient{
		Client:  client,
		BaseURL: "https://example.com/api",
	}

	apiClient.LoadConfig()

	if len(config.ArtistByID) == 0 {
		t.Errorf("Expected artists to be populated, but the map was empty")
	}

	for id, artist := range config.ArtistByID {
		if id != 1 {
			t.Errorf("Expected artist ID 1, got %d", id)
		}
		if artist.ID != 1 {
			t.Errorf("Expected artist ID 1, got %d", artist.ID)
		}
	}

}

func TestAPIClient(t *testing.T) {
	client := &config.APIClient{
		Client:  &http.Client{Timeout: 10 * time.Second},
		BaseURL: "https://example.com/api",
	}

	if client.BaseURL != "https://example.com/api" {
		t.Errorf("Expected BaseURL to be 'https://example.com/api', got '%s'", client.BaseURL)
	}
	if client.Client.Timeout != 10*time.Second {
		t.Errorf("Expected Client timeout to be 10 seconds, got %v", client.Client.Timeout)
	}
}
