package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/yukkat/go-metrics-api/internal/handler"
	"github.com/yukkat/go-metrics-api/internal/logger"
	"github.com/yukkat/go-metrics-api/internal/middleware"
	"github.com/yukkat/go-metrics-api/internal/repository"
)

func newServer() http.Handler {
	return newServerWithStorage(repository.NewMemStorage(), nil)
}

func newServerWithStorage(storage *repository.MemStorage, db *sql.DB) http.Handler {
	updateHandler := handler.NewUpdateHandler(storage)
	updateHandlerJSON := handler.NewUpdateHandlerJSON(storage)
	valueHandler := handler.NewValueHandler(storage)
	valueHandlerJSON := handler.NewValueHandlerJSON(storage)
	allMetricsHandler := handler.NewAllMetricsHandler(storage)
	pingHandler := handler.NewPingHandler(db)

	r := chi.NewRouter()

	r.Use(logger.Handle)
	r.Use(middleware.GzipMiddleware)

	r.Get("/", allMetricsHandler.ServeHTTP)
	r.Get("/ping", pingHandler.ServeHTTP)
	r.Post("/update/{type}/{name}/{value}", updateHandler.ServeHTTP)
	r.Post("/update/", updateHandlerJSON.ServeHTTP)
	r.Post("/update", updateHandlerJSON.ServeHTTP)
	r.Get("/value/{type}/{name}", valueHandler.ServeHTTP)
	r.Post("/value/", valueHandlerJSON.ServeHTTP)
	r.Post("/value", valueHandlerJSON.ServeHTTP)

	return r
}

func openDB(dsn string) *sql.DB {
	if dsn == "" {
		return nil
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Log.Error("failed to open database", zap.Error(err))
		return nil
	}

	return db
}

func startSaveTicker(
	ctx context.Context,
	storage *repository.MemStorage,
	path string,
	interval int,
) {
	if path == "" || interval <= 0 {
		return
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	go func() {
		defer ticker.Stop()

		for range ticker.C {
			if err := storage.SaveToFile(path); err != nil {
				logger.Log.Error("failed to save metrics", zap.Error(err))
			}
		}
	}()
}

func main() {
	parseFlags()

	if err := logger.Initialize(flagLogLevel); err != nil {
		panic(err)
	}

	db := openDB(flagDatabaseDSN)
	if db != nil {
		defer db.Close()
	}

	storage := repository.NewMemStorage()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if flagRestore && flagFileStoragePath != "" {
		if err := storage.LoadFromFile(flagFileStoragePath); err != nil {
			logger.Log.Error("failed to load metrics", zap.Error(err))
		}
	}

	if flagStoreInterval <= 0 {
		storage.ConfigureFileStorage(flagFileStoragePath, true)
	} else {
		startSaveTicker(ctx, storage, flagFileStoragePath, flagStoreInterval)
	}

	server := newServerWithStorage(storage, db)

	if err := http.ListenAndServe(flagRunAddr, server); err != nil {
		panic(err)
	}
}
