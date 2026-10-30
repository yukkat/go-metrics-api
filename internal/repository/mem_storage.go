package repository

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"

	"github.com/yukkat/go-metrics-api/internal/model"
)

type Storage interface {
	SetGauge(name string, value float64)
	AddCounter(name string, value int64)

	GetAll() []Metric
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	GetAllGauges() []Metric
	GetAllCounters() []Metric
}

type MemStorage struct {
	mu       sync.RWMutex
	gauges   map[string]float64
	counters map[string]int64

	filePath string
	syncSave bool
	saveMu   sync.Mutex
}

type Metric struct {
	Name  string
	Value string
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

// ConfigureFileStorage включает сохранение метрик в файл.
// Если syncSave равно true, файл сохраняется после каждого обновления.
func (s *MemStorage) ConfigureFileStorage(path string, syncSave bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.filePath = path
	s.syncSave = syncSave
}

func (s *MemStorage) SetGauge(name string, value float64) {
	s.mu.Lock()
	if s.gauges == nil {
		s.gauges = make(map[string]float64)
	}
	s.gauges[name] = value
	s.mu.Unlock()

	s.saveAfterUpdate()
}

func (s *MemStorage) AddCounter(name string, value int64) {
	s.mu.Lock()
	if s.counters == nil {
		s.counters = make(map[string]int64)
	}
	s.counters[name] += value
	s.mu.Unlock()

	s.saveAfterUpdate()
}

func (s *MemStorage) GetAll() []Metric {
	gauges := s.GetAllGauges()
	counters := s.GetAllCounters()
	metrics := make([]Metric, 0, len(gauges)+len(counters))

	metrics = append(metrics, gauges...)
	metrics = append(metrics, counters...)

	return metrics
}

func (s *MemStorage) GetAllGauges() []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metrics := make([]Metric, 0, len(s.gauges))

	for name, value := range s.gauges {
		metrics = append(metrics, Metric{
			Name:  name,
			Value: strconv.FormatFloat(value, 'f', -1, 64),
		})
	}

	return metrics
}

func (s *MemStorage) GetAllCounters() []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metrics := make([]Metric, 0, len(s.counters))

	for name, value := range s.counters {
		metrics = append(metrics, Metric{
			Name:  name,
			Value: strconv.FormatInt(value, 10),
		})
	}

	return metrics
}

func (s *MemStorage) GetGauge(name string) (float64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.gauges[name]
	return value, ok
}

func (s *MemStorage) GetCounter(name string) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.counters[name]
	return value, ok
}

// SaveToFile сохраняет текущие значения метрик в файл в формате JSON.
func (s *MemStorage) SaveToFile(path string) error {
	s.saveMu.Lock()
	defer s.saveMu.Unlock()

	s.mu.RLock()
	metrics := make([]model.Metrics, 0, len(s.gauges)+len(s.counters))

	for name, value := range s.gauges {
		value := value
		metrics = append(metrics, model.Metrics{
			ID:    name,
			MType: model.Gauge,
			Value: &value,
		})
	}

	for name, value := range s.counters {
		value := value
		metrics = append(metrics, model.Metrics{
			ID:    name,
			MType: model.Counter,
			Delta: &value,
		})
	}
	s.mu.RUnlock()

	data, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// LoadFromFile загружает значения метрик из файла.
func (s *MemStorage) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var metrics []model.Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.gauges == nil {
		s.gauges = make(map[string]float64)
	}
	if s.counters == nil {
		s.counters = make(map[string]int64)
	}

	for _, metric := range metrics {
		switch metric.MType {
		case model.Gauge:
			if metric.Value != nil {
				s.gauges[metric.ID] = *metric.Value
			}
		case model.Counter:
			if metric.Delta != nil {
				s.counters[metric.ID] = *metric.Delta
			}
		}
	}

	return nil
}

func (s *MemStorage) saveAfterUpdate() {
	s.mu.RLock()
	path := s.filePath
	syncSave := s.syncSave
	s.mu.RUnlock()

	if syncSave && path != "" {
		_ = s.SaveToFile(path)
	}
}
