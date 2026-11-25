package repository

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestDBStorage(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("cannot open database: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE gauges (name TEXT PRIMARY KEY, value DOUBLE PRECISION NOT NULL)`); err != nil {
		t.Fatalf("cannot create gauges: %v", err)
	}
	if _, err := db.Exec(`CREATE TABLE counters (name TEXT PRIMARY KEY, value BIGINT NOT NULL)`); err != nil {
		t.Fatalf("cannot create counters: %v", err)
	}

	storage := NewDBStorage(db)

	storage.SetGauge("Alloc", 25.5)
	value, ok := storage.GetGauge("Alloc")
	if !ok || value != 25.5 {
		t.Errorf("expected gauge 25.5, got %v, found=%t", value, ok)
	}

	storage.SetGauge("Alloc", 30.5)
	value, ok = storage.GetGauge("Alloc")
	if !ok || value != 30.5 {
		t.Errorf("expected gauge 30.5, got %v, found=%t", value, ok)
	}

	storage.AddCounter("requests", 10)
	storage.AddCounter("requests", 5)
	counter, ok := storage.GetCounter("requests")
	if !ok || counter != 15 {
		t.Errorf("expected counter 15, got %d, found=%t", counter, ok)
	}

	metrics := storage.GetAll()
	if len(metrics) != 2 {
		t.Fatalf("expected 2 metrics, got %d", len(metrics))
	}
}
