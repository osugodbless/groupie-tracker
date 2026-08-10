package handlers_test

import (
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/osugodbless/groupie-tracker/internal/client"
	"github.com/osugodbless/groupie-tracker/internal/handlers"
)

var (
	firstMockArtist = client.Artist{
		ID:           1,
		Name:         "Mock Artist",
		Members:      []string{"Alice", "Bob"},
		CreationYear: 2020,
		FirstAlbum:   "2021-01-01",
		Image:        "https://example.com/img.jpg",
		DatesLocation: map[string][]string{
			"New_York":    {"2023-01-01", "2023-01-02"},
			"Los_Angeles": {"2023-02-01"},
		},
	}

	secondMockArtist = client.Artist{
		ID:           7,
		Name:         "Mock Artist Two",
		Members:      []string{"Charlie"},
		CreationYear: 2015,
		FirstAlbum:   "2016-05-10",
		Image:        "https://example.com/img2.jpg",
		DatesLocation: map[string][]string{
			"London": {"2023-03-01"},
		},
	}
)

var testTemplate = `
		{{define "base.tmpl"}}<div>Base Template</div>{{end}}
		{{define "artistsDetails.tmpl"}}<div>Artist Details: {{.Name}}</div>{{end}}
		{{define "tour-dates.tmpl"}}<div>Tour Dates: {{.Name}}</div>{{end}}
	`

func setup() *handlers.Application {
	tmpl := template.Must(template.New("test").Parse(testTemplate))
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	data := map[int]client.Artist{
		1: firstMockArtist,
		7: secondMockArtist,
	}
	service := handlers.NewBandArtistService(data)

	app := handlers.NewApplication(tmpl, logger, service)

	return app
}

func TestHomeHandler(t *testing.T) {
	app := setup()

	tests := []struct {
		name   string
		method string
		path   string
		status int
	}{
		{name: "GET returns 200", method: "GET", path: "/", status: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			app.HomeHandler(w, r)

			if w.Code != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, w.Code)
			}
		})
	}
}

func TestHomeHandlerResponseBody(t *testing.T) {
	app := setup()

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	app.HomeHandler(w, r)

	body := w.Body.String()
	if !strings.Contains(body, "Base Template") {
		t.Errorf("expected response to contain 'Base Template', got: %s", body)
	}
}

func TestArtistHandler(t *testing.T) {
	app := setup()

	tests := []struct {
		name      string
		method    string
		path      string
		pathValue string
		status    int
		contains  string
	}{
		{
			name:      "GET valid artist returns 200",
			method:    "GET",
			path:      "/artist/1",
			pathValue: "1",
			status:    http.StatusOK,
			contains:  "Mock Artist",
		},
		{
			name:      "non-numeric ID returns 400",
			method:    "GET",
			path:      "/artist/abc",
			pathValue: "abc",
			status:    http.StatusBadRequest,
		},
		{
			name:      "non-existent artist returns 404",
			method:    "GET",
			path:      "/artist/9999",
			pathValue: "9999",
			status:    http.StatusNotFound,
		},
		{
			name:      "second artist returns 200",
			method:    "GET",
			path:      "/artist/7",
			pathValue: "7",
			status:    http.StatusOK,
			contains:  "Mock Artist Two",
		},
		{
			name:      "empty string ID returns 400",
			method:    "GET",
			path:      "/artist/",
			pathValue: "",
			status:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.path, nil)
			r.SetPathValue("id", tt.pathValue)
			w := httptest.NewRecorder()

			app.GetArtistHandler(w, r)

			if w.Code != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, w.Code)
			}

			if tt.contains != "" && !strings.Contains(w.Body.String(), tt.contains) {
				t.Errorf("expected body to contain %q, got: %s", tt.contains, w.Body.String())
			}
		})
	}
}

func TestTourDatesHandler(t *testing.T) {
	app := setup()

	tests := []struct {
		name      string
		method    string
		path      string
		pathValue string
		status    int
		contains  string
	}{
		{
			name:      "GET valid artist returns 200",
			method:    "GET",
			path:      "/artist/1/tour-data",
			pathValue: "1",
			status:    http.StatusOK,
			contains:  "Tour Dates",
		},
		{
			name:      "non-numeric ID returns 400",
			method:    "GET",
			path:      "/artist/abc/tour-data",
			pathValue: "abc",
			status:    http.StatusBadRequest,
		},
		{
			name:      "non-existent artist returns 404",
			method:    "GET",
			path:      "/artist/9999/tour-data",
			pathValue: "9999",
			status:    http.StatusNotFound,
		},
		{
			name:      "second artist returns 200",
			method:    "GET",
			path:      "/artist/7/tour-data",
			pathValue: "7",
			status:    http.StatusOK,
			contains:  "Tour Dates",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.path, nil)
			r.SetPathValue("id", tt.pathValue)
			w := httptest.NewRecorder()

			app.TourDatesHandler(w, r)

			if w.Code != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, w.Code)
			}

			if tt.contains != "" && !strings.Contains(w.Body.String(), tt.contains) {
				t.Errorf("expected body to contain %q, got: %s", tt.contains, w.Body.String())
			}
		})
	}
}

func TestRenderTemplateError(t *testing.T) {
	app := setup()
	app.Templates = template.Must(template.New("test").Parse(`{{define "something.tmpl"}}ok{{end}}`))

	r := httptest.NewRequest("GET", "/artist/1", nil)
	r.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	app.GetArtistHandler(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for template execution error, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Internal Server Error") {
		t.Errorf("expected 'Internal Server Error' in body, got: %s", w.Body.String())
	}
}

func TestTourDatesRenderTemplateError(t *testing.T) {
	app := setup()
	app.Templates = template.Must(template.New("test").Parse(`{{define "something.tmpl"}}ok{{end}}`))

	r := httptest.NewRequest("GET", "/artist/1/tour-data", nil)
	r.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	app.TourDatesHandler(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for template execution error, got %d", w.Code)
	}
}

func TestGetArtistByIDNotFound(t *testing.T) {
	app := setup()

	r := httptest.NewRequest("GET", "/artist/42", nil)
	r.SetPathValue("id", "42")
	w := httptest.NewRecorder()

	app.GetArtistHandler(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Artist not found") {
		t.Errorf("expected 'Artist not found' in body, got: %s", body)
	}
}

func TestArtistHandlerResponseBody(t *testing.T) {
	app := setup()

	r := httptest.NewRequest("GET", "/artist/1", nil)
	r.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	app.GetArtistHandler(w, r)

	body := w.Body.String()
	if !strings.Contains(body, "Artist Details: Mock Artist") {
		t.Errorf("expected body to contain artist details, got: %s", body)
	}
}

func TestTourDatesHandlerResponseBody(t *testing.T) {
	app := setup()

	r := httptest.NewRequest("GET", "/artist/1/tour-data", nil)
	r.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	app.TourDatesHandler(w, r)

	body := w.Body.String()
	if !strings.Contains(body, "Tour Dates: Mock Artist") {
		t.Errorf("expected body to contain tour dates, got: %s", body)
	}
}

func TestHomeHandlerEmptyArtistMap(t *testing.T) {
	app := setup()
	app.ArtistService = handlers.NewBandArtistService(map[int]client.Artist{}) // Empty artist map

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	app.HomeHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 with empty artist map, got %d", w.Code)
	}
}

func TestHomeHandlerRenderTemplateError(t *testing.T) {
	app := setup()
	app.Templates = template.Must(template.New("test").Parse(`{{define "something.tmpl"}}ok{{end}}`))

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	app.HomeHandler(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for template execution error, got %d", w.Code)
	}
}
