package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/yukkat/go-metrics-api/internal/agent"
	"github.com/yukkat/go-metrics-api/internal/repository"
)

func main() {
	config := parseFlags()

	host := fmt.Sprintf("http://%s", config.flagRunAddr)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var storage = repository.NewMemStorage()

	pollTicker := time.NewTicker(time.Duration(config.flagPollInterval*int(time.Second)))
	reportTicker := time.NewTicker(time.Duration(config.flagReportInterval*int(time.Second)))

	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			agent.CollectMetrics(storage)

		case <-reportTicker.C:
			agent.SendMetric(client, host, storage)
		}
	}
}