package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/yukkat/go-metrics-api/internal/repository"
)

func TestValueHandler(t *testing.T) {
	storage := repository.NewMemStorage()
	storage.SetGauge("temperature", 25.5)
	storage.AddCounter("requests", 10)

	handler := NewValueHandler(storage)

	r := chi.NewRouter()
	r.Get("/value/{type}/{name}", handler.ServeHTTP)

	tests := []struct {
		name   string
		path   string
		status int
	}{
		{
			name:   "existing gauge",
			path:   "/value/gauge/temperature",
			status: http.StatusOK,
		},
		{
			name:   "existing counter",
			path:   "/value/counter/requests",
			status: http.StatusOK,
		},
		{
			name:   "missing gauge",
			path:   "/value/gauge/nosuch",
			status: http.StatusNotFound,
		},
		{
			name:   "missing counter",
			path:   "/value/counter/nosuch",
			status: http.StatusNotFound,
		},
		{
			name:   "unknown type",
			path:   "/value/something/temperature",
			status: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.path, nil)
			res := httptest.NewRecorder()

			r.ServeHTTP(res, req)

			if res.Code != test.status {
				t.Errorf("expected status %d, got %d", test.status, res.Code)
			}
		})
	}
}

func TestValueHandlerBody(t *testing.T) {
	storage := repository.NewMemStorage()
	storage.SetGauge("temperature", 25.5)

	handler := NewValueHandler(storage)

	r := chi.NewRouter()
	r.Get("/value/{type}/{name}", handler.ServeHTTP)

	req := httptest.NewRequest(http.MethodGet, "/value/gauge/temperature", nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	if res.Body.String() != "25.5" {
		t.Errorf("expected body 25.5, got %q", res.Body.String())
	}
}
