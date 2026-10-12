package repository

type Storage interface {
	SetGauge(name string, value float64)
	AddCounter(name string, value int64)
}

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

var _ Storage = (*MemStorage)(nil)

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (s *MemStorage) SetGauge(name string, value float64) {
    if s.gauges == nil {
		s.gauges = make(map[string]float64)
	}
	s.gauges[name] = value
}

func (s *MemStorage) AddCounter(name string, value int64) {
	if s.counters == nil {
		s.counters = make(map[string]int64)
	}
	s.counters[name] += value
}
