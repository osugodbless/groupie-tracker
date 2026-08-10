package routes_test

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/osugodbless/groupie-tracker/internal/client"
	"github.com/osugodbless/groupie-tracker/internal/handlers"
	"github.com/osugodbless/groupie-tracker/internal/routes"
)

var (
	firstMockArtist = client.Artist{
		ID:   1,
		Name: "First Artist",
		DatesLocation: map[string][]string{
			"2023-01-01": {"New York", "Los Angeles"},
		},
	}

	secondMockArtist = client.Artist{
		ID:   7,
		Name: "Seventh Artist",
		DatesLocation: map[string][]string{
			"2023-01-01": {"New York", "Los Angeles"},
		},
	}
)

var testTemplate = `
		{{define "base.tmpl"}}<div>Base Template</div>{{end}}
		{{define "artistsDetails.tmpl"}}<div>Artist Details: {{.Name}}</div>{{end}}
		{{define "tour-dates.tmpl"}}<div>Tour Dates: {{.Name}}</div>{{end}}
		{{define "artists:grid"}}<div>Artist Grid</div>{{end}}
		{{define "artists-page"}}<div>Artists Page</div>{{end}}
	`

func TestRoutes(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse(testTemplate))

	artists := map[int]client.Artist{
		1: firstMockArtist,
		7: secondMockArtist,
	}

	service := handlers.NewBandArtistService(artists)
	app := handlers.NewApplication(tmpl, nil, service)
	mux := routes.Routes(app)

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		status int
	}{
		{"Test Home Route", "GET", "/", "", http.StatusOK},
		{"Test Artists Route", "GET", "/artists", "", http.StatusOK},
		{"Test Search Artists Route", "GET", "/artists/findx?search=First", "", http.StatusOK},
		{"Test Artist Route", "GET", "/artists/1", "", http.StatusOK},
		{"Test Tour Dates Route", "GET", "/artists/1/tour-data", "", http.StatusOK},
		{"Test Non-existent Artist Route", "GET", "/artists/9999", "", http.StatusNotFound},
		{"Test Home POST returns 405", "POST", "/", "", http.StatusMethodNotAllowed},
		{"Test Artist POST returns 405", "POST", "/artists/1", "", http.StatusMethodNotAllowed},
		{"Test Tour Dates POST returns 405", "POST", "/artists/1/tour-data", "", http.StatusMethodNotAllowed},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, r)

			res := w.Result()

			if res.StatusCode != test.status {
				t.Errorf("expected status %v; got %v", test.status, res.StatusCode)
			}
		})
	}
}
