package main

import (
	"net/http"
	"time"
	"github.com/yukkat/go-metrics-api/internal/agent"
)

var pollInterval = 2 * time.Second
var timeout = 5 * time.Second
var host = "http://localhost:8080"

func main() {
    client := &http.Client{
        Timeout: timeout,
    }

    for {
        agent.CollectAndSend(client, host)
        time.Sleep(pollInterval)
    }
}