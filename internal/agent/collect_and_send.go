package agent

import (
    "net/http"
    "log"
)

func CollectAndSend(client *http.Client, host string) {
    metrics := CollectMetrics()

    for _, metric := range metrics {
        if err := SendMetric(
            client,
            host,
            metric,
        ); err != nil {
            log.Printf(
                "failed to send metric %s: %v",
                metric.ID,
                err,
            )
        }
    }
}

