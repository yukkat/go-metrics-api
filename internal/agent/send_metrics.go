package agent

import (
    "fmt"
    "net/http"
    "strconv"
    "github.com/yukkat/go-metrics-api/internal/model"
)

func SendMetric(
    client *http.Client,
    serverURL string,
    metric model.Metrics,
) error {
    var value string

    if metric.MType == "gauge" && metric.Value != nil {
        value = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
    } else if metric.MType == "counter" && metric.Delta != nil {
        value = strconv.FormatInt(*metric.Delta, 10)
    } else {
        return fmt.Errorf("invalid metric: %+v", metric)
    }

    url := fmt.Sprintf(
        "%s/update/%s/%s/%s",
        serverURL,
        metric.MType,
        metric.ID,
        value,
    )

    req, err := http.NewRequest(http.MethodPost, url, nil)
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "text/plain")

    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return fmt.Errorf("server returned status %s", resp.Status)
    }

    return nil
}