package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/yukkat/go-metrics-api/internal/agent"
)

var pollInterval = 2 * time.Second
var timeout = 5 * time.Second

func main() {
    parseFlags()

    host := fmt.Sprintf(`http://localhost%s`, flagRunAddr)
    client := &http.Client{
        Timeout: timeout,
    }
    fmt.Println(host)

    for {
        agent.CollectAndSend(client, host)
        time.Sleep(pollInterval)
    }
}