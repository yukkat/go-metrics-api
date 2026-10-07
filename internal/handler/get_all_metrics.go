package handler;

import (
    "github.com/yukkat/go-metrics-api/internal/repository"
    "net/http"
    "fmt"
)

type AllMetricsHandler struct {
	storage repository.Storage
}

func NewAllMetricsHandler(storage repository.Storage) *AllMetricsHandler {
	return &AllMetricsHandler{storage: storage}
}

func (h *AllMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    metrics := h.storage.GetAll()

    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    fmt.Fprintln(w, "<html><body>")

    for _, metric := range metrics {
        fmt.Fprintf(w, "<p>%s: %s</p>\n", metric.Name, metric.Value)
    }

    fmt.Fprintln(w, "</body></html>")
}