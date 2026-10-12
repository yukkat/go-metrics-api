package main

import (
    "flag"
    "time"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string
var pollInterval time.Duration
var reportInterval time.Duration

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
    // регистрируем переменную flagRunAddr
    // как аргумент -a со значением :8080 по умолчанию
    flag.StringVar(&flagRunAddr, "a", defHost, "address and port to run server")
    flag.DurationVar(&pollInterval, "r", defPollInterval, "address and port to run server")
    flag.DurationVar(&reportInterval, "p", defReportInterval, "address and port to run server")
    // парсим переданные серверу аргументы в зарегистрированные переменные
    flag.Parse()
}