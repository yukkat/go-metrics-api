package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yukkat/go-metrics-api/internal/model"
	"github.com/yukkat/go-metrics-api/internal/repository"
	"github.com/yukkat/go-metrics-api/internal/helpers"
)

func SendMetricJSON(
    client *http.Client,
    serverURL string,
    storage *repository.MemStorage,
) error {
    counters := storage.GetAllCounters()
    gauges := storage.GetAllGauges()

    for _, m := range gauges {
       if err := makeRequestJSON(client, serverURL, model.Gauge, m.Name, m.Value); err != nil {
           return err
       }
    }
    for _, n := range counters {
       if err := makeRequestJSON(client, serverURL, model.Counter, n.Name, n.Value); err != nil {
           return err
       }
    }

    return nil
}

func makeRequestJSON(
    client *http.Client,
    serverURL string,
    mType string,
    name string,
    value string,
    ) error {

    url := fmt.Sprintf(
        "%s/update",
        serverURL,
    )
    var req model.Metrics

    switch mType {
        case model.Gauge:
            val, err := strconv.ParseFloat(value, 64)
            if err != nil {
                return err
            }
           req = model.Metrics{
                ID: name,
                MType: mType,
                Value: &val,
            }

        case model.Counter:
            val, err := strconv.ParseInt(value, 10, 64)
            if err != nil {
                return err
            }
           req = model.Metrics{
                ID: name,
                MType: mType,
                Delta: &val,
            }

        default:
            return fmt.Errorf("metric type %s is not found", mType)
    }
        // Сериализуем структуру в JSON
        body, err := json.Marshal(req)
        if err != nil {
            return fmt.Errorf("failed to encode JSON: %w", err)
        }

       return helpers.SendJSONHelper(client, body, url)

}