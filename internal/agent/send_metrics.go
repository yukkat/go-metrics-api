package agent

import (
    "fmt"
    "net/http"
    "github.com/yukkat/go-metrics-api/internal/model"
    "github.com/yukkat/go-metrics-api/internal/repository"
)

func SendMetric(
    client *http.Client,
    serverURL string,
    storage *repository.MemStorage,
) error {
    counters := storage.GetAllCounters()
    gauges := storage.GetAllGauges()

    for _, m := range gauges {
       if err := makeRequest(client, serverURL, model.Gauge, m.Name, m.Value); err != nil {
           return err
       }
    }
    for _, n := range counters {
       if err := makeRequest(client, serverURL, model.Counter, n.Name, n.Value); err != nil {
           return err
       }
    }

    return nil
}

func makeRequest(
    client *http.Client,
    serverURL string,
    mType string,
    name string,
    value string,
    ) error {

    url := fmt.Sprintf(
        "%s/update/%s/%s/%s",
        serverURL,
        mType,
        name,
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