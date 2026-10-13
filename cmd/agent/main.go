package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/yukkat/go-metrics-api/internal/agent"
)


func main() {
    parseFlags()

    host := fmt.Sprintf(`http://%s`, flagRunAddr)

    client := &http.Client{
        Timeout: reportInterval,
    }
    fmt.Println(host)

    for {
        agent.CollectAndSend(client, host)
        time.Sleep(pollInterval)
    }
}