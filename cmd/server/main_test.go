package main

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	server := newServer()

	tests := []struct {
		name   string
		path   string
		status int
	}{
		{
			name:   "update",
			path:   "/update/gauge/temperature/25.5",
			status: http.StatusOK,
		},
		{
			name:   "not found",
			path:   "/hello",
			status: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, test.path, nil)
			res := httptest.NewRecorder()

			server.ServeHTTP(res, req)

			if res.Code != test.status {
				t.Errorf("expected status %d, got %d", test.status, res.Code)
			}
		})
	}
}

func TestServerAcceptsGzip(t *testing.T) {
	server := newServer()

	var compressed bytes.Buffer
	zw := gzip.NewWriter(&compressed)
	if _, err := zw.Write([]byte(`{"id":"temperature","type":"gauge","value":25.5}`)); err != nil {
		t.Fatalf("cannot compress body: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("cannot close gzip writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/update", &compressed)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/value/gauge/temperature", nil)
	getRes := httptest.NewRecorder()

	server.ServeHTTP(getRes, getReq)

	if strings.TrimSpace(getRes.Body.String()) != "25.5" {
		t.Errorf("expected value 25.5, got %q", getRes.Body.String())
	}
}
