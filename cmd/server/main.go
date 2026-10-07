package main

import (
	"net/http"

	"github.com/yukkat/go-metrics-api/internal/handler"
	"github.com/yukkat/go-metrics-api/internal/repository"
	"github.com/go-chi/chi/v5"
)

func newServer() http.Handler {
	storage := repository.NewMemStorage()
	updateHandler := handler.NewUpdateHandler(storage)
	valueHandler := handler.NewValueHandler(storage)
	allMetricsHandler := handler.NewAllMetricsHandler(storage)

	r := chi.NewRouter()

    r.Get("/", allMetricsHandler.ServeHTTP)
	r.Post("/update/{type}/{name}/{value}", updateHandler.ServeHTTP)
	r.Get("/value/{type}/{name}", valueHandler.ServeHTTP)

	return r
}

func main() {
	server := newServer()

	if err := http.ListenAndServe(":8080", server); err != nil {
		panic(err)
	}
}