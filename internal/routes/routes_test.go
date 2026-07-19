package routes_test

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osugodbless/groupie-tracker/internal/config"
	"github.com/osugodbless/groupie-tracker/internal/handlers"
	"github.com/osugodbless/groupie-tracker/internal/routes"
)

var firstMockArtist = config.Artist{
	ID:   1,
	Name: "Mock Artist",
	DatesLocation: map[string][]string{
		"2023-01-01": {"New York", "Los Angeles"},
	},
}

var secondMockArtist = config.Artist{
	ID:   7,
	Name: "Mock Artist",
	DatesLocation: map[string][]string{
		"2023-01-01": {"New York", "Los Angeles"},
	},
}

var testTemplate = `
		{{define "base.html"}}<div>Base Template</div>{{end}}
		{{define "artistsDetails.html"}}<div>Artist Details: {{.Name}}</div>{{end}}
		{{define "tour-dates.html"}}<div>Tour Dates: {{.Name}}</div>{{end}}
	`

func TestRoutes(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse(testTemplate))

	config.ArtistByID = map[int]config.Artist{
		1: firstMockArtist,
		7: secondMockArtist,
	}

	app := &handlers.Application{
		Templates: tmpl,
	}

	mux := routes.Routes(app)

	tests := []struct {
		name   string
		method string
		path   string
		status int
	}{
		{"Test Home Route", "GET", "/", http.StatusOK},
		{"Test Artist Route", "GET", "/artist/1", http.StatusOK},
		{"Test Tour Dates Route", "GET", "/artist/1/tour-data", http.StatusOK},
		{"Test Non-existent Artist Route", "GET", "/artist/9999", http.StatusNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, test.path, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, r)

			res := w.Result()

			if res.StatusCode != test.status {
				t.Errorf("expected status %v; got %v", test.status, res.StatusCode)
			}
		})
	}
}
