package handlers_test

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/osugodbless/groupie-tracker/internal/config"
	"github.com/osugodbless/groupie-tracker/internal/handlers"
)

var firstMockArtist = config.Artist{
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

var secondMockArtist = config.Artist{
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

var testTemplate = `
		{{define "base.html"}}<div>Base Template</div>{{end}}
		{{define "artistsDetails.html"}}<div>Artist Details: {{.Name}}</div>{{end}}
		{{define "tour-dates.html"}}<div>Tour Dates: {{.Name}}</div>{{end}}
	`

func setup() (*handlers.Application, map[int]config.Artist) {
	tmpl := template.Must(template.New("test").Parse(testTemplate))

	data := map[int]config.Artist{
		1: firstMockArtist,
		7: secondMockArtist,
	}
	config.ArtistByID = data

	app := &handlers.Application{
		Templates: tmpl,
	}
	return app, data
}

func TestHomeHandler(t *testing.T) {
	app, _ := setup()

	tests := []struct {
		name   string
		method string
		path   string
		status int
	}{
		{name: "GET returns 200", method: "GET", path: "/", status: http.StatusOK},
		{name: "POST returns 405", method: "POST", path: "/", status: http.StatusMethodNotAllowed},
		{name: "PUT returns 405", method: "PUT", path: "/", status: http.StatusMethodNotAllowed},
		{name: "DELETE returns 405", method: "DELETE", path: "/", status: http.StatusMethodNotAllowed},
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
	app, _ := setup()

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	app.HomeHandler(w, r)

	body := w.Body.String()
	if !strings.Contains(body, "Base Template") {
		t.Errorf("expected response to contain 'Base Template', got: %s", body)
	}
}

func TestArtistHandler(t *testing.T) {
	app, _ := setup()

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
			name:      "POST returns 405",
			method:    "POST",
			path:      "/artist/1",
			pathValue: "1",
			status:    http.StatusMethodNotAllowed,
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

			app.ArtistHandler(w, r)

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
	app, _ := setup()

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
			name:      "POST returns 405",
			method:    "POST",
			path:      "/artist/1/tour-data",
			pathValue: "1",
			status:    http.StatusMethodNotAllowed,
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
	config.ArtistByID = map[int]config.Artist{
		1: firstMockArtist,
	}

	tmpl := template.Must(template.New("test").Parse(`{{define "something.html"}}ok{{end}}`))

	app := &handlers.Application{
		Templates: tmpl,
	}

	r := httptest.NewRequest("GET", "/artist/1", nil)
	r.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	app.ArtistHandler(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for template execution error, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Internal Server Error") {
		t.Errorf("expected 'Internal Server Error' in body, got: %s", w.Body.String())
	}
}

func TestTourDatesRenderTemplateError(t *testing.T) {
	config.ArtistByID = map[int]config.Artist{
		1: firstMockArtist,
	}

	tmpl := template.Must(template.New("test").Parse(`{{define "something.html"}}ok{{end}}`))

	app := &handlers.Application{
		Templates: tmpl,
	}

	r := httptest.NewRequest("GET", "/artist/1/tour-data", nil)
	r.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	app.TourDatesHandler(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for template execution error, got %d", w.Code)
	}
}

func TestGetArtistByIDNotFound(t *testing.T) {
	setup()

	r := httptest.NewRequest("GET", "/artist/42", nil)
	r.SetPathValue("id", "42")
	w := httptest.NewRecorder()

	app := &handlers.Application{
		Templates: template.Must(template.New("test").Parse(testTemplate)),
	}

	app.ArtistHandler(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Artist not found") {
		t.Errorf("expected 'Artist not found' in body, got: %s", body)
	}
}

func TestArtistHandlerResponseBody(t *testing.T) {
	app, _ := setup()

	r := httptest.NewRequest("GET", "/artist/1", nil)
	r.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	app.ArtistHandler(w, r)

	body := w.Body.String()
	if !strings.Contains(body, "Artist Details: Mock Artist") {
		t.Errorf("expected body to contain artist details, got: %s", body)
	}
}

func TestTourDatesHandlerResponseBody(t *testing.T) {
	app, _ := setup()

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
	config.ArtistByID = map[int]config.Artist{}

	tmpl := template.Must(template.New("test").Parse(testTemplate))
	app := &handlers.Application{Templates: tmpl}

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	app.HomeHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 with empty artist map, got %d", w.Code)
	}
}

func TestHomeHandlerRenderTemplateError(t *testing.T) {
	config.ArtistByID = map[int]config.Artist{
		1: firstMockArtist,
	}

	tmpl := template.Must(template.New("test").Parse(`{{define "something.html"}}ok{{end}}`))

	app := &handlers.Application{
		Templates: tmpl,
	}

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	app.HomeHandler(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for template execution error, got %d", w.Code)
	}
}

func TestArtistHandlerMethodNotAllowed(t *testing.T) {
	app, _ := setup()

	methods := []string{"POST", "PUT", "DELETE", "PATCH"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			r := httptest.NewRequest(method, "/artist/1", nil)
			r.SetPathValue("id", "1")
			w := httptest.NewRecorder()

			app.ArtistHandler(w, r)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405 for %s, got %d", method, w.Code)
			}
		})
	}
}

func TestTourDatesHandlerMethodNotAllowed(t *testing.T) {
	app, _ := setup()

	methods := []string{"POST", "PUT", "DELETE", "PATCH"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			r := httptest.NewRequest(method, "/artist/1/tour-data", nil)
			r.SetPathValue("id", "1")
			w := httptest.NewRecorder()

			app.TourDatesHandler(w, r)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405 for %s, got %d", method, w.Code)
			}
		})
	}
}
