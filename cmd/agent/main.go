package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/yukkat/go-metrics-api/internal/agent"
)

func main() {
	parseFlags()

	host := fmt.Sprintf("http://%s", flagRunAddr)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	pollTicker := time.NewTicker(time.Duration(flagPollInterval))
	reportTicker := time.NewTicker(time.Duration(flagReportInterval))

	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			agent.CollectMetrics()

		case <-reportTicker.C:
			agent.CollectAndSend(client, host)
		}
	}
}