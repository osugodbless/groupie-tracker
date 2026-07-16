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

func TestRoutes(t *testing.T) {
	tmpl := template.Must(template.ParseFiles("../../testdata/base.html"))

	config.LoadConfig()

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
		{"Test Home Page", "GET", "/", http.StatusOK},
		{"Test Artist Page", "GET", "/artist/1", http.StatusOK},
		{"Test Tour Dates Page", "GET", "/artist/1/tour-data", http.StatusOK},
		{"Test Non-existent Artist", "GET", "/artist/9999", http.StatusNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, test.path, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, r)

			res := w.Result()

			if res.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", res.StatusCode)
			}
		})
	}
}
