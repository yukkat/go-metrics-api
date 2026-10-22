package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/yukkat/go-metrics-api/internal/logger"
	"github.com/yukkat/go-metrics-api/internal/model"
	"github.com/yukkat/go-metrics-api/internal/repository"
)


type UpdateHandlerJSON struct {
	storage repository.Storage
}

func NewUpdateHandlerJSON(storage repository.Storage) *UpdateHandlerJSON {
	return &UpdateHandlerJSON{storage: storage}
}

func (h *UpdateHandlerJSON) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var req model.Metrics
    dec := json.NewDecoder(r.Body)
    if err := dec.Decode(&req); err != nil {
        logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

	switch req.MType {
        case model.Gauge:
            if req.Value == nil {
                http.Error(w, "invalid gauge value", http.StatusBadRequest)
                return
            }

            h.storage.SetGauge(req.ID, *req.Value)

        case model.Counter:
            if req.Delta == nil {
                http.Error(w, "invalid counter value", http.StatusBadRequest)
                return
            }

            h.storage.AddCounter(req.ID, *req.Delta)

        default:
            http.Error(w, "invalid metric type", http.StatusBadRequest)
            return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}