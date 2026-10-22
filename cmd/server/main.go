package main

import (
	"net/http"

	"github.com/yukkat/go-metrics-api/internal/handler"
	"github.com/yukkat/go-metrics-api/internal/logger"
	"github.com/yukkat/go-metrics-api/internal/repository"
	"github.com/go-chi/chi/v5"
)

func newServer() http.Handler {
	storage := repository.NewMemStorage()
	updateHandler := handler.NewUpdateHandler(storage)
	updateHandlerJSON := handler.NewUpdateHandlerJSON(storage)
	valueHandler := handler.NewValueHandler(storage)
	valueHandlerJSON := handler.NewValueHandlerJSON(storage)
	allMetricsHandler := handler.NewAllMetricsHandler(storage)

	r := chi.NewRouter()

	r.Use(logger.Handle)

	r.Get("/", allMetricsHandler.ServeHTTP)
	r.Post("/update/{type}/{name}/{value}", updateHandler.ServeHTTP)
	r.Post("/update/", updateHandlerJSON.ServeHTTP)
	r.Post("/update", updateHandlerJSON.ServeHTTP)
	r.Get("/value/{type}/{name}", valueHandler.ServeHTTP)
	r.Post("/value/", valueHandlerJSON.ServeHTTP)
	r.Post("/value", valueHandlerJSON.ServeHTTP)

	return r
}

func main() {
	parseFlags()

	if err := logger.Initialize(flagLogLevel); err != nil {
		panic(err)
	}

	server := newServer()

	if err := http.ListenAndServe(flagRunAddr, server); err != nil {
		panic(err)
	}
}