package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yukkat/go-metrics-api/internal/repository"
)

func TestUpdateHandler(t *testing.T) {
	storage := repository.NewMemStorage()
	handler := NewUpdateHandler(storage)

	tests := []struct {
		name   string
		method string
		path   string
		status int
	}{
		{
			name:   "gauge",
			method: http.MethodPost,
			path:   "/update/gauge/temperature/25.5",
			status: http.StatusOK,
		},
		{
			name:   "counter",
			method: http.MethodPost,
			path:   "/update/counter/requests/10",
			status: http.StatusOK,
		},
		{
			name:   "wrong method",
			method: http.MethodGet,
			path:   "/update/gauge/temperature/25.5",
			status: http.StatusMethodNotAllowed,
		},
		{
			name:   "invalid gauge value",
			method: http.MethodPost,
			path:   "/update/gauge/temperature/hello",
			status: http.StatusBadRequest,
		},
		{
			name:   "invalid path",
			method: http.MethodPost,
			path:   "/something/gauge/temperature/25.5",
			status: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.path, nil)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			if res.Code != test.status {
				t.Errorf("expected status %d, got %d", test.status, res.Code)
			}
		})
	}
}
