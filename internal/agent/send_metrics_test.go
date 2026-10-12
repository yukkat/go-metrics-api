package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yukkat/go-metrics-api/internal/model"
)

func TestSendMetric(t *testing.T) {
	value := 10.5
	delta := int64(10)

	tests := []struct {
		name    string
		metric  model.Metrics
		status  int
		wantErr bool
	}{
		{
			name: "gauge",
			metric: model.Metrics{
				ID:    "temperature",
				MType: model.Gauge,
				Value: &value,
			},
			status:  http.StatusOK,
			wantErr: false,
		},
		{
			name: "counter",
			metric: model.Metrics{
				ID:    "requests",
				MType: model.Counter,
				Delta: &delta,
			},
			status:  http.StatusOK,
			wantErr: false,
		},
		{
			name: "invalid metric",
			metric: model.Metrics{
				ID:    "temperature",
				MType: "unknown",
			},
			status:  http.StatusOK,
			wantErr: true,
		},
		{
			name: "server error",
			metric: model.Metrics{
				ID:    "temperature",
				MType: model.Gauge,
				Value: &value,
			},
			status:  http.StatusInternalServerError,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.status)
			}))
			defer server.Close()

			err := SendMetric(server.Client(), server.URL, test.metric)

			if test.wantErr && err == nil {
				t.Error("expected error, got nil")
			}

			if !test.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
