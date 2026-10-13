package main

import (
    "flag"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string
var flagReportInterval int
var flagPollInterval int

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
    // регистрируем переменную flagRunAddr
    // как аргумент -a со значением :8080 по умолчанию
    flag.StringVar(&flagRunAddr, "a", defHost, "address and port to run server")
    flag.IntVar(&flagReportInterval, "r", defReportInterval, "report interval in seconds")
    flag.IntVar(&flagPollInterval, "p", defPollInterval, "poll interval in seconds")
    // парсим переданные серверу аргументы в зарегистрированные переменные
    flag.Parse()
}