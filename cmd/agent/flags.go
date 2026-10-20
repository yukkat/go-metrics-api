package main

import (
    "flag"
    "os"
    "strconv"
    "fmt"
)

type Config struct {
    flagRunAddr string
    flagReportInterval int
    flagPollInterval int
    flagLogLevel string
}

var defConfig = Config{
    flagRunAddr: "localhost:8080",
    flagPollInterval: 2,
    flagReportInterval: 10,
    flagLogLevel: "info",
}


// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() Config {
    var config = Config{
        flagRunAddr: defConfig.flagRunAddr,
        flagPollInterval: defConfig.flagPollInterval,
        flagReportInterval: defConfig.flagReportInterval,
        flagLogLevel: defConfig.flagLogLevel,
    }

    // регистрируем переменную flagRunAddr
    // как аргумент -a со значением :8080 по умолчанию
    flag.StringVar(&config.flagRunAddr, "a", defConfig.flagRunAddr, "address and port to run server")
    flag.StringVar(&config.flagLogLevel, "l", defConfig.flagLogLevel, "address and port to run server")
    flag.IntVar(&config.flagPollInterval, "p", defConfig.flagPollInterval, "poll interval in seconds")
    flag.IntVar(&config.flagReportInterval, "r", defConfig.flagReportInterval, "report interval in seconds")
    // парсим переданные серверу аргументы в зарегистрированные переменные
    flag.Parse()

    var _ = readEnvs(&config);

    return config;
}

func readEnvs(c *Config) error {
        // для случаев, когда в переменной окружения ADDRESS присутствует непустое значение,
        // переопределим адрес запуска сервера,
        // даже если он был передан через аргумент командной строки
        if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
            c.flagRunAddr = envRunAddr
        }
        if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
            c.flagLogLevel = envLogLevel
        }
        if reportInterval := os.Getenv("REPORT_INTERVAL"); reportInterval != "" {
            value, err := strconv.Atoi(reportInterval)
            if err != nil {
                return fmt.Errorf("invalid REPORT_INTERVAL %q: %w", reportInterval, err)
            }

            c.flagReportInterval = value
        }
        if pollInterval := os.Getenv("POLL_INTERVAL"); pollInterval != "" {
            value, err := strconv.Atoi(pollInterval)
            if err != nil {
                return fmt.Errorf("invalid POLL_INTERVAL %q: %w", pollInterval, err)
            }

            c.flagPollInterval = value
        }
        return nil
}