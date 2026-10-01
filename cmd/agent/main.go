package main

import (
	"log"
	"net/http"
	"time"
	"github.com/yukkat/go-metrics-api/internal/agent"
)

var pollInterval = 2 * time.Second
var timeout = 5 * time.Second

func main() {
    client := &http.Client{
        Timeout: timeout,
    }

    for {
        metrics := agent.CollectMetrics()

        for _, metric := range metrics {
            if err := agent.SendMetric(
                client,
                "http://localhost:8080",
                metric,
            ); err != nil {
                log.Printf(
                    "failed to send metric %s: %v",
                    metric.ID,
                    err,
                )
            }
        }

        time.Sleep(pollInterval)
    }
}