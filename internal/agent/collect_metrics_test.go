package agent

import (
	"testing"

	"github.com/yukkat/go-metrics-api/internal/repository"
)

func TestCollectMetrics(t *testing.T) {
	pollCount = 0
	storage := repository.NewMemStorage()

	CollectMetrics(storage)

	if _, ok := storage.GetGauge("Alloc"); !ok {
		t.Error("expected gauge Alloc, not found")
	}

	if _, ok := storage.GetGauge("RandomValue"); !ok {
		t.Error("expected gauge RandomValue, not found")
	}

	value, ok := storage.GetCounter("PollCount")
	if !ok {
		t.Error("expected counter PollCount, not found")
	}

	if value != 1 {
		t.Errorf("expected PollCount 1, got %d", value)
	}
}

func TestCollectMetricsPollCount(t *testing.T) {
	pollCount = 0
	storage := repository.NewMemStorage()

	CollectMetrics(storage)
	CollectMetrics(storage)

	if pollCount != 2 {
		t.Errorf("expected pollCount 2, got %d", pollCount)
	}
}
