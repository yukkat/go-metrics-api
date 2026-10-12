package repository

import (
    "strconv"
)

type Storage interface {
	SetGauge(name string, value float64)
	AddCounter(name string, value int64)

	GetAll() []Metric
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
}

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

type Metric struct {
	Name  string
	Value string
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

func (s *MemStorage) GetAll() []Metric {
	metrics := make([]Metric, 0, len(s.gauges)+len(s.counters))

	for name, value := range s.gauges {
		metrics = append(metrics, Metric{
			Name:  name,
			Value: strconv.FormatFloat(value, 'f', -1, 64),
		})
	}

	for name, value := range s.counters {
		metrics = append(metrics, Metric{
			Name:  name,
			Value: strconv.FormatInt(value, 10),
		})
	}

	return metrics
}

func (s *MemStorage) GetGauge(name string) (float64, bool) {
	value, ok := s.gauges[name]
	return value, ok
}

func (s *MemStorage) GetCounter(name string) (int64, bool) {
	value, ok := s.counters[name]
	return value, ok
}