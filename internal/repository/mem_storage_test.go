package repository

import "testing"

func TestSetGauge(t *testing.T) {
	storage := NewMemStorage()

	storage.SetGauge("Alloc", 25.5)

	if storage.gauges["Alloc"] != 25.5 {
		t.Errorf("expected 25.5, got %v", storage.gauges["Alloc"])
	}
}

func TestAddCounter(t *testing.T) {
	storage := NewMemStorage()

	storage.AddCounter("requests", 10)

	if storage.counters["requests"] != 10 {
		t.Errorf("expected 10, got %v", storage.counters["requests"])
	}

	storage.AddCounter("requests", 5)

	if storage.counters["requests"] != 15 {
		t.Errorf("expected 15, got %v", storage.counters["requests"])
	}
}

func TestGetGauge(t *testing.T) {
	storage := NewMemStorage()
	storage.SetGauge("Alloc", 25.5)

	value, ok := storage.GetGauge("Alloc")
	if !ok {
		t.Error("expected gauge to be found")
	}

	if value != 25.5 {
		t.Errorf("expected 25.5, got %v", value)
	}
}

func TestGetGaugeMissing(t *testing.T) {
	storage := NewMemStorage()

	_, ok := storage.GetGauge("nosuch")

	if ok {
		t.Error("expected gauge not to be found")
	}
}

func TestGetCounter(t *testing.T) {
	storage := NewMemStorage()
	storage.AddCounter("requests", 10)

	value, ok := storage.GetCounter("requests")
	if !ok {
		t.Error("expected counter to be found")
	}

	if value != 10 {
		t.Errorf("expected 10, got %v", value)
	}
}

func TestGetCounterMissing(t *testing.T) {
	storage := NewMemStorage()

	_, ok := storage.GetCounter("nosuch")

	if ok {
		t.Error("expected counter not to be found")
	}
}

func TestGetAll(t *testing.T) {
	storage := NewMemStorage()
	storage.SetGauge("Alloc", 25.5)
	storage.AddCounter("requests", 10)

	metrics := storage.GetAll()

	if len(metrics) != 2 {
		t.Errorf("expected 2 metrics, got %d", len(metrics))
	}
}

func TestGetAllEmpty(t *testing.T) {
	storage := NewMemStorage()

	metrics := storage.GetAll()

	if len(metrics) != 0 {
		t.Errorf("expected 0 metrics, got %d", len(metrics))
	}
}
