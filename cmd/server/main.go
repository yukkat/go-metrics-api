package main

import (
	"net/http"
	"strings"

	"github.com/yukkat/go-metrics-api/internal/handler"
	"github.com/yukkat/go-metrics-api/internal/repository"
)

func main() {
	storage := repository.NewMemStorage()
	updateHandler := handler.NewUpdateHandler(storage)

	server := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		if path == "/update" || strings.HasPrefix(path, "/update/") {
			updateHandler.ServeHTTP(res, req)
			return
		}

		http.NotFound(res, req)
	})

	if err := http.ListenAndServe(":8080", server); err != nil {
		panic(err)
	}
}
