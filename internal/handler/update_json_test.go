package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/yukkat/go-metrics-api/internal/repository"
)

func TestUpdateHandlerJSON(t *testing.T) {
	storage := repository.NewMemStorage()
	handler := NewUpdateHandlerJSON(storage)

	r := chi.NewRouter()
	r.Post("/update", handler.ServeHTTP)

	tests := []struct {
		name   string
		method string
		body   string
		status int
	}{
		{
			name:   "gauge",
			method: http.MethodPost,
			body:   `{"id":"temperature","type":"gauge","value":25.5}`,
			status: http.StatusOK,
		},
		{
			name:   "counter",
			method: http.MethodPost,
			body:   `{"id":"requests","type":"counter","delta":10}`,
			status: http.StatusOK,
		},
		{
			name:   "wrong method",
			method: http.MethodGet,
			body:   `{"id":"temperature","type":"gauge","value":25.5}`,
			status: http.StatusMethodNotAllowed,
		},
		{
			name:   "invalid json",
			method: http.MethodPost,
			body:   `{"id":"temperature",`,
			status: http.StatusBadRequest,
		},
		{
			name:   "empty body",
			method: http.MethodPost,
			body:   ``,
			status: http.StatusBadRequest,
		},
		{
			name:   "unknown type",
			method: http.MethodPost,
			body:   `{"id":"temperature","type":"unknown","value":25.5}`,
			status: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, "/update", strings.NewReader(test.body))
			res := httptest.NewRecorder()

			r.ServeHTTP(res, req)

			if res.Code != test.status {
				t.Errorf("expected status %d, got %d", test.status, res.Code)
			}
		})
	}
}

func TestUpdateHandlerJSONSetsGauge(t *testing.T) {
	tests := []struct {
		name  string
		body  string
		value float64
	}{
		{
			name:  "usual value",
			body:  `{"id":"temperature","type":"gauge","value":25.5}`,
			value: 25.5,
		},
		{
			name:  "zero value",
			body:  `{"id":"temperature","type":"gauge","value":0}`,
			value: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := repository.NewMemStorage()
			handler := NewUpdateHandlerJSON(storage)

			req := httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(test.body))
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			value, ok := storage.GetGauge("temperature")
			if !ok {
				t.Fatalf("expected gauge temperature to be set")
			}

			if value != test.value {
				t.Errorf("expected gauge %v, got %v", test.value, value)
			}
		})
	}
}

func TestUpdateHandlerJSONAddsCounter(t *testing.T) {
	storage := repository.NewMemStorage()
	storage.AddCounter("requests", 5)
	handler := NewUpdateHandlerJSON(storage)

	req := httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(`{"id":"requests","type":"counter","delta":10}`))
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	value, ok := storage.GetCounter("requests")
	if !ok {
		t.Fatalf("expected counter requests to be set")
	}

	if value != 15 {
		t.Errorf("expected counter 15, got %v", value)
	}
}

func TestUpdateHandlerJSONContentType(t *testing.T) {
	storage := repository.NewMemStorage()
	handler := NewUpdateHandlerJSON(storage)

	req := httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(`{"id":"temperature","type":"gauge","value":25.5}`))
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	contentType := res.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json, got %q", contentType)
	}
}

func TestUpdateHandlerJSONWithoutValue(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		status int
	}{
		{
			name:   "gauge without value",
			body:   `{"id":"temperature","type":"gauge"}`,
			status: http.StatusBadRequest,
		},
		{
			name:   "counter without delta",
			body:   `{"id":"requests","type":"counter"}`,
			status: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := repository.NewMemStorage()
			handler := NewUpdateHandlerJSON(storage)

			req := httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(test.body))
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			if res.Code != test.status {
				t.Errorf("expected status %d, got %d", test.status, res.Code)
			}
		})
	}
}
