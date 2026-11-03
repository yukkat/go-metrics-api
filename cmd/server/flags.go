package main

import (
	"flag"
	"os"
	"strconv"
)

const (
	defaultRunAddr         = "localhost:8080"
	defaultLogLevel        = "info"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "metrics-storage.json"
	defaultRestore         = false
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var (
	flagRunAddr         string
	flagLogLevel        string
	flagStoreInterval   int
	flagFileStoragePath string
	flagRestore         bool
)

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", defaultRunAddr, "address and port to run server")
	flag.StringVar(&flagLogLevel, "l", defaultLogLevel, "log level")
	flag.IntVar(&flagStoreInterval, "i", defaultStoreInterval, "interval in seconds to save metrics to file (0 means synchronous saving)")
	flag.StringVar(&flagFileStoragePath, "f", defaultFileStoragePath, "path to file with saved metrics")
	flag.BoolVar(&flagRestore, "r", defaultRestore, "load saved metrics from file on startup")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	// для случаев, когда в переменной окружения RUN_ADDR присутствует непустое значение,
	// переопределим адрес запуска сервера,
	// даже если он был передан через аргумент командной строки
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		value, err := strconv.Atoi(envStoreInterval)
		if err == nil {
			flagStoreInterval = value
		}
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		flagFileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		value, err := strconv.ParseBool(envRestore)
		if err == nil {
			flagRestore = value
		}
	}
}
