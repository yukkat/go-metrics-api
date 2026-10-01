package agent

import (
	"math/rand/v2"
	"runtime"
	"github.com/yukkat/go-metrics-api/internal/model"
)

func gauge(id string, value float64) model.Metrics {
	return model.Metrics{
		ID:    id,
		MType: "gauge",
		Value: &value,
	}
}

func counter(id string, delta int64) model.Metrics {
    return model.Metrics{
        ID: id,
        MType: "counter",
        Delta: &delta,
    }
}

var pollCount int64 = 0;
func updatePollCount(){
    pollCount += 1
}

func getRandomValue() float64 {
    return rand.Float64()
}

func CollectMetrics() []model.Metrics {
    updatePollCount()
    randomValue := getRandomValue()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return []model.Metrics{
		gauge("Alloc", float64(m.Alloc)),
		gauge("BuckHashSys", float64(m.BuckHashSys)),
		gauge("Frees", float64(m.Frees)),
		gauge("GCCPUFraction", m.GCCPUFraction),
		gauge("GCSys", float64(m.GCSys)),
		gauge("HeapAlloc", float64(m.HeapAlloc)),
		gauge("HeapIdle", float64(m.HeapIdle)),
		gauge("HeapInuse", float64(m.HeapInuse)),
		gauge("HeapObjects", float64(m.HeapObjects)),
		gauge("HeapReleased", float64(m.HeapReleased)),
		gauge("HeapSys", float64(m.HeapSys)),
		gauge("LastGC", float64(m.LastGC)),
		gauge("Lookups", float64(m.Lookups)),
		gauge("MCacheInuse", float64(m.MCacheInuse)),
		gauge("MCacheSys", float64(m.MCacheSys)),
		gauge("MSpanInuse", float64(m.MSpanInuse)),
		gauge("MSpanSys", float64(m.MSpanSys)),
		gauge("Mallocs", float64(m.Mallocs)),
		gauge("NextGC", float64(m.NextGC)),
		gauge("NumForcedGC", float64(m.NumForcedGC)),
		gauge("NumGC", float64(m.NumGC)),
		gauge("OtherSys", float64(m.OtherSys)),
		gauge("PauseTotalNs", float64(m.PauseTotalNs)),
		gauge("StackInuse", float64(m.StackInuse)),
		gauge("StackSys", float64(m.StackSys)),
		gauge("Sys", float64(m.Sys)),
		gauge("TotalAlloc", float64(m.TotalAlloc)),
		counter("PollCount", pollCount),
		gauge("RandomValue", randomValue),
	}
}