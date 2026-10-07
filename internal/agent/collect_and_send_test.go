package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCollectAndSend(t *testing.T) {
	count := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	CollectAndSend(server.Client(), server.URL)

	if count == 0 {
		t.Error("expected metrics to be sent")
	}
}
