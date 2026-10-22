package main

import (
	"fmt"
	"net/http"
	"time"

    "go.uber.org/zap"

	"github.com/yukkat/go-metrics-api/internal/agent"
	"github.com/yukkat/go-metrics-api/internal/repository"
	"github.com/yukkat/go-metrics-api/internal/logger"
)

func main() {
	config := parseFlags()

	host := fmt.Sprintf("http://%s", config.flagRunAddr)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	if err := logger.Initialize(config.flagLogLevel); err != nil {
        panic(err)
    }

    logger.Log.Info("Running agent", zap.String("address", config.flagRunAddr))

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
			if err := agent.SendMetricJSON(client, host, storage); err != nil {
				logger.Log.Error("failed to send metrics", zap.Error(err))
			}
		}
	}
}