package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/yukkat/go-metrics-api/internal/handler"
	"github.com/yukkat/go-metrics-api/internal/logger"
	"github.com/yukkat/go-metrics-api/internal/middleware"
	"github.com/yukkat/go-metrics-api/internal/repository"
)

func newServer() http.Handler {
	return newServerWithStorage(repository.NewMemStorage())
}

func newServerWithStorage(storage *repository.MemStorage) http.Handler {
	updateHandler := handler.NewUpdateHandler(storage)
	updateHandlerJSON := handler.NewUpdateHandlerJSON(storage)
	valueHandler := handler.NewValueHandler(storage)
	valueHandlerJSON := handler.NewValueHandlerJSON(storage)
	allMetricsHandler := handler.NewAllMetricsHandler(storage)

	r := chi.NewRouter()

	r.Use(logger.Handle)
	r.Use(middleware.GzipMiddleware)

	r.Get("/", allMetricsHandler.ServeHTTP)
	r.Post("/update/{type}/{name}/{value}", updateHandler.ServeHTTP)
	r.Post("/update/", updateHandlerJSON.ServeHTTP)
	r.Post("/update", updateHandlerJSON.ServeHTTP)
	r.Get("/value/{type}/{name}", valueHandler.ServeHTTP)
	r.Post("/value/", valueHandlerJSON.ServeHTTP)
	r.Post("/value", valueHandlerJSON.ServeHTTP)

	return r
}

func startSaveTicker(storage *repository.MemStorage, path string, interval int) {
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

	storage := repository.NewMemStorage()

	if flagRestore && flagFileStoragePath != "" {
		if err := storage.LoadFromFile(flagFileStoragePath); err != nil {
			logger.Log.Error("failed to load metrics", zap.Error(err))
		}
	}

	if flagStoreInterval <= 0 {
		storage.ConfigureFileStorage(flagFileStoragePath, true)
	} else {
		startSaveTicker(storage, flagFileStoragePath, flagStoreInterval)
	}

	server := newServerWithStorage(storage)

	if err := http.ListenAndServe(flagRunAddr, server); err != nil {
		panic(err)
	}
}
