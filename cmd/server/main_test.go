package main

import (
	"net/http"
	"net/http/httptest"
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