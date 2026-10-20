package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestHandle(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)
	old := Log
	Log = zap.New(core)
	t.Cleanup(func() { Log = old })

	handler := Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("hello"))
	}))

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/alloc/1", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, res.Code)
	}

	if logs.Len() != 2 {
		t.Fatalf("expected 2 log entries, got %d", logs.Len())
	}

	reqEntry := logs.All()[0]
	if got := reqEntry.ContextMap()["uri"]; got != "/update/gauge/alloc/1" {
		t.Errorf("unexpected uri: %v", got)
	}
	if got := reqEntry.ContextMap()["method"]; got != http.MethodPost {
		t.Errorf("unexpected method: %v", got)
	}
	if _, ok := reqEntry.ContextMap()["duration"]; !ok {
		t.Error("duration is missing in request log")
	}

	respEntry := logs.All()[1]
	if got := respEntry.ContextMap()["status"]; got != int64(http.StatusCreated) {
		t.Errorf("unexpected status: %v", got)
	}
	if got := respEntry.ContextMap()["size"]; got != int64(5) {
		t.Errorf("unexpected size: %v", got)
	}
}
