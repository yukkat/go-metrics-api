package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yukkat/go-metrics-api/internal/repository"
)

func TestSendMetric(t *testing.T) {
	count := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	storage := repository.NewMemStorage()
	storage.SetGauge("temperature", 10.5)
	storage.AddCounter("requests", 3)

	err := SendMetric(server.Client(), server.URL, storage)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 requests, got %d", count)
	}
}

func TestSendMetricEmptyStorage(t *testing.T) {
	count := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	storage := repository.NewMemStorage()

	err := SendMetric(server.Client(), server.URL, storage)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if count != 0 {
		t.Errorf("expected 0 requests, got %d", count)
	}
}
