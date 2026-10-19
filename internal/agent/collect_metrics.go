package agent

import (
	"math/rand/v2"
	"runtime"
	"github.com/yukkat/go-metrics-api/internal/model"
	"github.com/yukkat/go-metrics-api/internal/repository"
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

func CollectMetrics(s *repository.MemStorage) {
    updatePollCount()
    randomValue := getRandomValue()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

		s.SetGauge("Alloc", float64(m.Alloc))
		s.SetGauge("BuckHashSys", float64(m.BuckHashSys))
		s.SetGauge("Frees", float64(m.Frees))
		s.SetGauge("GCCPUFraction", m.GCCPUFraction)
		s.SetGauge("GCSys", float64(m.GCSys))
		s.SetGauge("HeapAlloc", float64(m.HeapAlloc))
		s.SetGauge("HeapIdle", float64(m.HeapIdle))
		s.SetGauge("HeapInuse", float64(m.HeapInuse))
		s.SetGauge("HeapObjects", float64(m.HeapObjects))
		s.SetGauge("HeapReleased", float64(m.HeapReleased))
		s.SetGauge("HeapSys", float64(m.HeapSys))
		s.SetGauge("LastGC", float64(m.LastGC))
		s.SetGauge("Lookups", float64(m.Lookups))
		s.SetGauge("MCacheInuse", float64(m.MCacheInuse))
		s.SetGauge("MCacheSys", float64(m.MCacheSys))
		s.SetGauge("MSpanInuse", float64(m.MSpanInuse))
		s.SetGauge("MSpanSys", float64(m.MSpanSys))
		s.SetGauge("Mallocs", float64(m.Mallocs))
		s.SetGauge("NextGC", float64(m.NextGC))
		s.SetGauge("NumForcedGC", float64(m.NumForcedGC))
		s.SetGauge("NumGC", float64(m.NumGC))
		s.SetGauge("OtherSys", float64(m.OtherSys))
		s.SetGauge("PauseTotalNs", float64(m.PauseTotalNs))
		s.SetGauge("StackInuse", float64(m.StackInuse))
		s.SetGauge("StackSys", float64(m.StackSys))
		s.SetGauge("Sys", float64(m.Sys))
		s.SetGauge("TotalAlloc", float64(m.TotalAlloc))
		s.SetGauge("RandomValue", randomValue)
		s.AddCounter("PollCount", pollCount)
}