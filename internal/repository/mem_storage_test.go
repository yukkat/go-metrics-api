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

func TestSaveAndLoadFromFile(t *testing.T) {
	path := t.TempDir() + "/metrics.json"

	storage := NewMemStorage()
	storage.SetGauge("Alloc", 25.5)
	storage.AddCounter("requests", 10)

	if err := storage.SaveToFile(path); err != nil {
		t.Fatalf("cannot save metrics: %v", err)
	}

	restored := NewMemStorage()
	if err := restored.LoadFromFile(path); err != nil {
		t.Fatalf("cannot load metrics: %v", err)
	}

	value, ok := restored.GetGauge("Alloc")
	if !ok {
		t.Fatal("expected gauge Alloc to be loaded")
	}
	if value != 25.5 {
		t.Errorf("expected gauge 25.5, got %v", value)
	}

	counter, ok := restored.GetCounter("requests")
	if !ok {
		t.Fatal("expected counter requests to be loaded")
	}
	if counter != 10 {
		t.Errorf("expected counter 10, got %v", counter)
	}
}

func TestSyncSaveAfterUpdate(t *testing.T) {
	path := t.TempDir() + "/metrics.json"

	storage := NewMemStorage()
	storage.ConfigureFileStorage(path, true)
	storage.SetGauge("temperature", 25.5)

	restored := NewMemStorage()
	if err := restored.LoadFromFile(path); err != nil {
		t.Fatalf("cannot load metrics: %v", err)
	}

	value, ok := restored.GetGauge("temperature")
	if !ok {
		t.Fatal("expected gauge temperature to be saved")
	}
	if value != 25.5 {
		t.Errorf("expected gauge 25.5, got %v", value)
	}
}
