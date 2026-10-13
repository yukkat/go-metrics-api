package handler;

import (
    "github.com/yukkat/go-metrics-api/internal/repository"
    "net/http"
    "fmt"
    "github.com/go-chi/chi/v5"
)

type ValueHandler struct {
	storage repository.Storage
}

func NewValueHandler(storage repository.Storage) *ValueHandler {
	return &ValueHandler{storage: storage}
}

func (h *ValueHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    metricType := chi.URLParam(req, "type")
	metricName := chi.URLParam(req, "name")

	switch metricType {
	case "gauge":
		value, ok := h.storage.GetGauge(metricName)
		if !ok {
			http.NotFound(res, req)
			return
		}

		res.WriteHeader(http.StatusOK)
		fmt.Fprint(res, value)

	case "counter":
		value, ok := h.storage.GetCounter(metricName)

    fmt.Printf("metricName=%q, value=%d, ok=%t\n", metricName, value, ok)
		if !ok {
			http.NotFound(res, req)
			return
		}

		res.WriteHeader(http.StatusOK)
		fmt.Fprint(res, value)

	default:
		http.NotFound(res, req)
	}
}