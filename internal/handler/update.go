package handler

import (
	"net/http"
	"strconv"
	"strings"

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
	if req.Method != http.MethodPost {
		res.Header().Set("Allow", http.MethodPost)
		http.Error(res, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metricType, name, rawValue, status := parseUpdatePath(req.URL.Path)
	if status != http.StatusOK {
		if status == http.StatusNotFound {
			http.NotFound(res, req)
			return
		}
		http.Error(res, "invalid metric path", status)
		return
	}

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

func parseUpdatePath(path string) (metricType, name, value string, status int) {
	if path == "/update" {
		return "", "", "", http.StatusNotFound
	}
	if !strings.HasPrefix(path, "/update/") {
		return "", "", "", http.StatusNotFound
	}

	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) < 3 || parts[0] != "update" {
		return "", "", "", http.StatusNotFound
	}
	if parts[2] == "" {
		return "", "", "", http.StatusNotFound
	}
	if len(parts) != 4 {
		return "", "", "", http.StatusBadRequest
	}

	return parts[1], parts[2], parts[3], http.StatusOK
}
