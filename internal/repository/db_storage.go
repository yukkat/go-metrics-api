package repository

import (
	"database/sql"
	"strconv"

	"go.uber.org/zap"

	"github.com/yukkat/go-metrics-api/internal/logger"
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(db *sql.DB) *DBStorage {
	return &DBStorage{db: db}
}

func (s *DBStorage) SetGauge(name string, value float64) {
	_, err := s.db.Exec(
		`INSERT INTO gauges (name, value) VALUES ($1, $2)
		 ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value`,
		name, value,
	)
	if err != nil {
		logger.Log.Error("failed to set gauge", zap.String("name", name), zap.Error(err))
	}
}

func (s *DBStorage) AddCounter(name string, value int64) {
	_, err := s.db.Exec(
		`INSERT INTO counters (name, value) VALUES ($1, $2)
		 ON CONFLICT (name) DO UPDATE SET value = counters.value + EXCLUDED.value`,
		name, value,
	)
	if err != nil {
		logger.Log.Error("failed to add counter", zap.String("name", name), zap.Error(err))
	}
}

func (s *DBStorage) GetAll() []Metric {
	gauges := s.GetAllGauges()
	counters := s.GetAllCounters()

	metrics := make([]Metric, 0, len(gauges)+len(counters))
	metrics = append(metrics, gauges...)
	metrics = append(metrics, counters...)

	return metrics
}

func (s *DBStorage) GetAllGauges() []Metric {
	rows, err := s.db.Query(`SELECT name, value FROM gauges ORDER BY name`)
	if err != nil {
		logger.Log.Error("failed to get gauges", zap.Error(err))
		return []Metric{}
	}
	defer rows.Close()

	metrics := make([]Metric, 0)
	for rows.Next() {
		var name string
		var value float64
		if err := rows.Scan(&name, &value); err != nil {
			logger.Log.Error("failed to scan gauge", zap.Error(err))
			return metrics
		}

		metrics = append(metrics, Metric{
			Name:  name,
			Value: strconv.FormatFloat(value, 'f', -1, 64),
		})
	}
	if err := rows.Err(); err != nil {
		logger.Log.Error("failed to read gauges", zap.Error(err))
	}

	return metrics
}

func (s *DBStorage) GetAllCounters() []Metric {
	rows, err := s.db.Query(`SELECT name, value FROM counters ORDER BY name`)
	if err != nil {
		logger.Log.Error("failed to get counters", zap.Error(err))
		return []Metric{}
	}
	defer rows.Close()

	metrics := make([]Metric, 0)
	for rows.Next() {
		var name string
		var value int64
		if err := rows.Scan(&name, &value); err != nil {
			logger.Log.Error("failed to scan counter", zap.Error(err))
			return metrics
		}

		metrics = append(metrics, Metric{
			Name:  name,
			Value: strconv.FormatInt(value, 10),
		})
	}
	if err := rows.Err(); err != nil {
		logger.Log.Error("failed to read counters", zap.Error(err))
	}

	return metrics
}

func (s *DBStorage) GetGauge(name string) (float64, bool) {
	var value float64
	err := s.db.QueryRow(`SELECT value FROM gauges WHERE name = $1`, name).Scan(&value)
	if err == sql.ErrNoRows {
		return 0, false
	}
	if err != nil {
		logger.Log.Error("failed to get gauge", zap.String("name", name), zap.Error(err))
		return 0, false
	}

	return value, true
}

func (s *DBStorage) GetCounter(name string) (int64, bool) {
	var value int64
	err := s.db.QueryRow(`SELECT value FROM counters WHERE name = $1`, name).Scan(&value)
	if err == sql.ErrNoRows {
		return 0, false
	}
	if err != nil {
		logger.Log.Error("failed to get counter", zap.String("name", name), zap.Error(err))
		return 0, false
	}

	return value, true
}
