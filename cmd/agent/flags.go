package main

import (
    "flag"
)

type Config struct {
    flagRunAddr string
    flagReportInterval int
    flagPollInterval int
}

var defConfig = Config{
    flagRunAddr: "localhost:8080",
    flagPollInterval: 2,
    flagReportInterval: 10,
}


// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() Config {
    var config = Config{
        flagRunAddr: defConfig.flagRunAddr,
        flagPollInterval: defConfig.flagPollInterval,
        flagReportInterval: defConfig.flagReportInterval,
    }
    // регистрируем переменную flagRunAddr
    // как аргумент -a со значением :8080 по умолчанию
    flag.StringVar(&config.flagRunAddr, "a", defConfig.flagRunAddr, "address and port to run server")
    flag.IntVar(&config.flagPollInterval, "p", defConfig.flagPollInterval, "poll interval in seconds")
    flag.IntVar(&config.flagReportInterval, "r", defConfig.flagReportInterval, "report interval in seconds")
    // парсим переданные серверу аргументы в зарегистрированные переменные
    flag.Parse()

    return config;
}