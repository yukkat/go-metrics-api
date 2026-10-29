package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGzipMiddlewareCompressesResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello gzip"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	res := httptest.NewRecorder()

	GzipMiddleware(handler).ServeHTTP(res, req)

	if res.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected Content-Encoding gzip, got %q", res.Header().Get("Content-Encoding"))
	}

	zr, err := gzip.NewReader(res.Body)
	if err != nil {
		t.Fatalf("cannot create gzip reader: %v", err)
	}
	defer zr.Close()

	body, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("cannot read response body: %v", err)
	}

	if string(body) != "hello gzip" {
		t.Errorf("expected body hello gzip, got %q", string(body))
	}
}

func TestGzipMiddlewareDecompressesRequest(t *testing.T) {
	var gotBody string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("cannot read request body: %v", err)
			return
		}

		gotBody = string(body)
		w.WriteHeader(http.StatusOK)
	})

	var compressed bytes.Buffer
	zw := gzip.NewWriter(&compressed)
	if _, err := zw.Write([]byte("hello gzip")); err != nil {
		t.Fatalf("cannot compress body: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("cannot close gzip writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", &compressed)
	req.Header.Set("Content-Encoding", "gzip")
	res := httptest.NewRecorder()

	GzipMiddleware(handler).ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	if gotBody != "hello gzip" {
		t.Errorf("expected body hello gzip, got %q", gotBody)
	}
}

func TestGzipMiddlewareWithoutGzip(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader("hello"))
	res := httptest.NewRecorder()

	GzipMiddleware(handler).ServeHTTP(res, req)

	if res.Header().Get("Content-Encoding") != "" {
		t.Errorf("expected empty Content-Encoding, got %q", res.Header().Get("Content-Encoding"))
	}

	if res.Body.String() != "hello" {
		t.Errorf("expected body hello, got %q", res.Body.String())
	}
}

func TestGzipMiddlewareInvalidRequestGzip(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not gzip"))
	req.Header.Set("Content-Encoding", "gzip")
	res := httptest.NewRecorder()

	GzipMiddleware(handler).ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}
