package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yukkat/go-metrics-api/internal/repository"
)

func TestAllMetricsHandler(t *testing.T) {
	storage := repository.NewMemStorage()
	handler := NewAllMetricsHandler(storage)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}

func TestAllMetricsHandlerWithMetrics(t *testing.T) {
	storage := repository.NewMemStorage()
	storage.SetGauge("temperature", 25.5)
	storage.AddCounter("requests", 10)

	handler := NewAllMetricsHandler(storage)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	body := res.Body.String()

	if !strings.Contains(body, "temperature: 25.5") {
		t.Errorf("expected body to contain temperature: 25.5, got %q", body)
	}

	if !strings.Contains(body, "requests: 10") {
		t.Errorf("expected body to contain requests: 10, got %q", body)
	}
}
