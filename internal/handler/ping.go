package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

type PingHandler struct {
	db *sql.DB
}

func NewPingHandler(db *sql.DB) *PingHandler {
	return &PingHandler{db: db}
}

func (h *PingHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if h.db == nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}
