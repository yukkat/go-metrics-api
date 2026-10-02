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
