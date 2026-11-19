package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGzipTransportCompressesRequest(t *testing.T) {
	var gotBody string
	var gotEncoding string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotEncoding = r.Header.Get("Content-Encoding")

		zr, err := gzip.NewReader(r.Body)
		if err != nil {
			t.Errorf("cannot create gzip reader: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer zr.Close()

		body, err := io.ReadAll(zr)
		if err != nil {
			t.Errorf("cannot read request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		gotBody = string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &http.Client{
		Transport: &GzipTransport{},
	}

	body := `{"id":"temperature","type":"gauge","value":25.5}`
	resp, err := client.Post(server.URL, "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if gotEncoding != "gzip" {
		t.Errorf("expected Content-Encoding gzip, got %q", gotEncoding)
	}

	if gotBody != body {
		t.Errorf("expected body %s, got %s", body, gotBody)
	}
}
