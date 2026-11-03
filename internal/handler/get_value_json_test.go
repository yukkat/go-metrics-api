package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/yukkat/go-metrics-api/internal/model"
	"github.com/yukkat/go-metrics-api/internal/repository"
)

func TestValueHandlerJSON(t *testing.T) {
	storage := repository.NewMemStorage()
	storage.SetGauge("temperature", 25.5)
	storage.AddCounter("requests", 10)

	handler := NewValueHandlerJSON(storage)

	r := chi.NewRouter()
	r.Get("/value", handler.ServeHTTP)

	tests := []struct {
		name   string
		method string
		body   string
		status int
	}{
		{
			name:   "gauge",
			method: http.MethodGet,
			body:   `{"id":"temperature","type":"gauge"}`,
			status: http.StatusOK,
		},
		{
			name:   "counter",
			method: http.MethodGet,
			body:   `{"id":"requests","type":"counter"}`,
			status: http.StatusOK,
		},
		{
			name:   "missing gauge",
			method: http.MethodGet,
			body:   `{"id":"nosuch","type":"gauge"}`,
			status: http.StatusNotFound,
		},
		{
			name:   "missing counter",
			method: http.MethodGet,
			body:   `{"id":"nosuch","type":"counter"}`,
			status: http.StatusNotFound,
		},
		{
			name:   "unknown type",
			method: http.MethodGet,
			body:   `{"id":"temperature","type":"unknown"}`,
			status: http.StatusNotFound,
		},
		{
			name:   "invalid json",
			method: http.MethodGet,
			body:   `{"id":`,
			status: http.StatusBadRequest,
		},
		{
			name:   "empty body",
			method: http.MethodGet,
			body:   ``,
			status: http.StatusBadRequest,
		},
		{
			name:   "wrong method",
			method: http.MethodPost,
			body:   `{"id":"temperature","type":"gauge"}`,
			status: http.StatusMethodNotAllowed,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, "/value", strings.NewReader(test.body))
			res := httptest.NewRecorder()

			r.ServeHTTP(res, req)

			if res.Code != test.status {
				t.Errorf("expected status %d, got %d", test.status, res.Code)
			}
		})
	}
}

func TestValueHandlerJSONGaugeBody(t *testing.T) {
	storage := repository.NewMemStorage()
	storage.SetGauge("temperature", 25.5)
	storage.SetGauge("zero", 0)

	handler := NewValueHandlerJSON(storage)

	tests := []struct {
		name  string
		body  string
		id    string
		value float64
	}{
		{
			name:  "usual value",
			body:  `{"id":"temperature","type":"gauge"}`,
			id:    "temperature",
			value: 25.5,
		},
		{
			name:  "zero value",
			body:  `{"id":"zero","type":"gauge"}`,
			id:    "zero",
			value: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/value", strings.NewReader(test.body))
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			var got model.Metrics
			err := json.Unmarshal(res.Body.Bytes(), &got)
			if err != nil {
				t.Fatalf("cannot unmarshal response %q: %v", res.Body.String(), err)
			}

			if got.ID != test.id {
				t.Errorf("expected id %s, got %s", test.id, got.ID)
			}

			if got.MType != model.Gauge {
				t.Errorf("expected type %s, got %s", model.Gauge, got.MType)
			}

			if got.Value == nil {
				t.Fatalf("expected value to be set")
			}

			if *got.Value != test.value {
				t.Errorf("expected value %v, got %v", test.value, *got.Value)
			}
		})
	}
}

func TestValueHandlerJSONCounterBody(t *testing.T) {
	storage := repository.NewMemStorage()
	storage.AddCounter("requests", 10)

	handler := NewValueHandlerJSON(storage)

	req := httptest.NewRequest(http.MethodGet, "/value", strings.NewReader(`{"id":"requests","type":"counter"}`))
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	var got model.Metrics
	err := json.Unmarshal(res.Body.Bytes(), &got)
	if err != nil {
		t.Fatalf("cannot unmarshal response %q: %v", res.Body.String(), err)
	}

	if got.ID != "requests" {
		t.Errorf("expected id requests, got %s", got.ID)
	}

	if got.MType != model.Counter {
		t.Errorf("expected type %s, got %s", model.Counter, got.MType)
	}

	if got.Delta == nil {
		t.Fatalf("expected delta to be set")
	}

	if *got.Delta != 10 {
		t.Errorf("expected delta 10, got %v", *got.Delta)
	}
}

func TestValueHandlerJSONContentType(t *testing.T) {
	storage := repository.NewMemStorage()
	storage.SetGauge("temperature", 25.5)

	handler := NewValueHandlerJSON(storage)

	req := httptest.NewRequest(http.MethodGet, "/value", strings.NewReader(`{"id":"temperature","type":"gauge"}`))
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	contentType := res.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json, got %q", contentType)
	}
}
