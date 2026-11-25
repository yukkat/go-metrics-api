package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/yukkat/go-metrics-api/internal/repository"
)

func TestNewStorageFallsBackToMemory(t *testing.T) {
	oldDSN, oldPath := flagDatabaseDSN, flagFileStoragePath
	oldRestore, oldInterval := flagRestore, flagStoreInterval
	defer func() {
		flagDatabaseDSN, flagFileStoragePath = oldDSN, oldPath
		flagRestore, flagStoreInterval = oldRestore, oldInterval
	}()

	flagDatabaseDSN = ""
	flagFileStoragePath = ""
	flagRestore = false
	flagStoreInterval = 300

	storage, db := newStorage(context.Background())
	if db != nil {
		t.Error("expected no database")
	}
	if _, ok := storage.(*repository.MemStorage); !ok {
		t.Errorf("expected memory storage, got %T", storage)
	}
}

func TestNewStorageFallsBackToFile(t *testing.T) {
	oldDSN, oldPath := flagDatabaseDSN, flagFileStoragePath
	oldRestore, oldInterval := flagRestore, flagStoreInterval
	defer func() {
		flagDatabaseDSN, flagFileStoragePath = oldDSN, oldPath
		flagRestore, flagStoreInterval = oldRestore, oldInterval
	}()

	path := filepath.Join(t.TempDir(), "metrics.json")

	flagDatabaseDSN = ""
	flagFileStoragePath = path
	flagRestore = false
	flagStoreInterval = 0

	storage, db := newStorage(context.Background())
	if db != nil {
		t.Error("expected no database")
	}
	if _, ok := storage.(*repository.MemStorage); !ok {
		t.Fatalf("expected memory storage with file saving, got %T", storage)
	}

	storage.SetGauge("temperature", 25.5)
	if _, err := os.ReadFile(path); err != nil {
		t.Errorf("expected metrics file to be created: %v", err)
	}
}

func TestNewStorageFallsBackOnDatabaseError(t *testing.T) {
	oldDSN, oldPath := flagDatabaseDSN, flagFileStoragePath
	oldRestore, oldInterval := flagRestore, flagStoreInterval
	defer func() {
		flagDatabaseDSN, flagFileStoragePath = oldDSN, oldPath
		flagRestore, flagStoreInterval = oldRestore, oldInterval
	}()

	flagDatabaseDSN = "unknown://db"
	flagFileStoragePath = ""
	flagRestore = false
	flagStoreInterval = 300

	storage, db := newStorage(context.Background())
	if db != nil {
		t.Error("expected no database")
	}
	if _, ok := storage.(*repository.MemStorage); !ok {
		t.Errorf("expected memory storage, got %T", storage)
	}
}
