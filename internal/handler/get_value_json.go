package handler;

import (
    "net/http"
    "encoding/json"

	"go.uber.org/zap"

    "github.com/yukkat/go-metrics-api/internal/repository"
    "github.com/yukkat/go-metrics-api/internal/model"
    "github.com/yukkat/go-metrics-api/internal/logger"
)

type ValueHandlerJSON struct {
	storage repository.Storage
}

func NewValueHandlerJSON(storage repository.Storage) *ValueHandlerJSON {
	return &ValueHandlerJSON{storage: storage}
}

func (h *ValueHandlerJSON) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var req model.Metrics
        dec := json.NewDecoder(r.Body)
        if err := dec.Decode(&req); err != nil {
            logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
            http.Error(w, "invalid request body", http.StatusBadRequest)
            return
        }

        var resp model.Metrics

	switch req.MType {
        case model.Gauge:
            value, ok := h.storage.GetGauge(req.ID)
            if !ok {
                http.NotFound(w, r)
                return
            }

            resp = model.Metrics{
                ID: req.ID,
                MType: model.Gauge,
                Value: &value,
            }

        case model.Counter:
            value, ok := h.storage.GetCounter(req.ID)

            if !ok {
                http.NotFound(w, r)
                return
            }

            resp = model.Metrics{
                ID: req.ID,
                MType: model.Counter,
                Delta: &value,
            }

        default:
            http.NotFound(w, r)
            return
	}

	// заполняем модель ответа
    w.Header().Set("Content-Type", "application/json")

    // сериализуем ответ сервера
    enc := json.NewEncoder(w)
    if err := enc.Encode(resp); err != nil {
        logger.Log.Debug("error encoding response", zap.Error(err))
        return
    }
    logger.Log.Debug("sending HTTP 200 response")
}