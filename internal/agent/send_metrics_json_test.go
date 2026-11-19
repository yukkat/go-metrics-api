package agent

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/yukkat/go-metrics-api/internal/model"
	"github.com/yukkat/go-metrics-api/internal/repository"
)

func TestSendMetricJSON(t *testing.T) {
	var mu sync.Mutex
	var bodies []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/update" {
			t.Errorf("expected path /update, got %s", r.URL.Path)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected content type application/json, got %s", r.Header.Get("Content-Type"))
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("cannot read request body: %v", err)
		}

		mu.Lock()
		bodies = append(bodies, string(body))
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	storage := repository.NewMemStorage()
	storage.SetGauge("temperature", 10.5)
	storage.AddCounter("requests", 3)

	err := SendMetricJSON(server.Client(), server.URL, storage)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(bodies) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(bodies))
	}

	gaugeBody := `{"id":"temperature","type":"gauge","value":10.5}`
	if bodies[0] != gaugeBody {
		t.Errorf("expected gauge body %s, got %s", gaugeBody, bodies[0])
	}

	counterBody := `{"id":"requests","type":"counter","delta":3}`
	if bodies[1] != counterBody {
		t.Errorf("expected counter body %s, got %s", counterBody, bodies[1])
	}
}

func TestSendMetricJSONEmptyStorage(t *testing.T) {
	var mu sync.Mutex
	count := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	storage := repository.NewMemStorage()

	err := SendMetricJSON(server.Client(), server.URL, storage)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if count != 0 {
		t.Errorf("expected 0 requests, got %d", count)
	}
}

func TestSendMetricJSONServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	storage := repository.NewMemStorage()
	storage.SetGauge("temperature", 10.5)

	err := SendMetricJSON(server.Client(), server.URL, storage)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestMakeRequestJSONUnknownType(t *testing.T) {
	err := makeRequestJSON(http.DefaultClient, "http://localhost", "unknown", "temperature", "10.5")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestMakeRequestJSONInvalidValue(t *testing.T) {
	err := makeRequestJSON(http.DefaultClient, "http://localhost", model.Gauge, "temperature", "hello")
	if err == nil {
		t.Errorf("expected error for gauge, got nil")
	}

	err = makeRequestJSON(http.DefaultClient, "http://localhost", model.Counter, "requests", "hello")
	if err == nil {
		t.Errorf("expected error for counter, got nil")
	}
}
