package handler

import (
	"net/http"
	"strconv"
    "fmt"
	"github.com/go-chi/chi/v5"
	"github.com/yukkat/go-metrics-api/internal/repository"
)

const (
	gaugeType   = "gauge"
	counterType = "counter"
)

type UpdateHandler struct {
	storage repository.Storage
}

func NewUpdateHandler(storage repository.Storage) *UpdateHandler {
	return &UpdateHandler{storage: storage}
}

func (h *UpdateHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "type")
	name := chi.URLParam(req, "name")
	rawValue := chi.URLParam(req, "value")
	fmt.Println(metricType, name, rawValue)

	switch metricType {
	case gaugeType:
		value, err := strconv.ParseFloat(rawValue, 64)
		if err != nil {
			http.Error(res, "invalid gauge value", http.StatusBadRequest)
			return
		}

		h.storage.SetGauge(name, value)

	case counterType:
		value, err := strconv.ParseInt(rawValue, 10, 64)
		if err != nil {
			http.Error(res, "invalid counter value", http.StatusBadRequest)
			return
		}

		h.storage.AddCounter(name, value)

	default:
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
}