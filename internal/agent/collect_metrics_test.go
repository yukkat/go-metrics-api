package agent

import (
	"testing"

	"github.com/yukkat/go-metrics-api/internal/model"
)

func TestCollectMetrics(t *testing.T) {
	pollCount = 0

	metrics := CollectMetrics()

	if len(metrics) == 0 {
		t.Error("expected metrics, got empty slice")
	}

	var pollMetric model.Metrics
	var randomMetric model.Metrics

	for _, metric := range metrics {
		if metric.ID == "PollCount" {
			pollMetric = metric
		}

		if metric.ID == "RandomValue" {
			randomMetric = metric
		}
	}

	if pollMetric.Delta == nil {
		t.Error("PollCount delta is nil")
	}

	if *pollMetric.Delta != 1 {
		t.Errorf("expected PollCount 1, got %d", *pollMetric.Delta)
	}

	if randomMetric.Value == nil {
		t.Error("RandomValue value is nil")
	}
}

func TestCollectMetricsPollCount(t *testing.T) {
	pollCount = 0

	CollectMetrics()
	metrics := CollectMetrics()

	for _, metric := range metrics {
		if metric.ID == "PollCount" {
			if *metric.Delta != 2 {
				t.Errorf("expected PollCount 2, got %d", *metric.Delta)
			}
			return
		}
	}

	t.Error("PollCount metric not found")
}
