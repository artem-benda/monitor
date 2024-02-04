package main

import (
	"net/http"
	"time"

	"github.com/artem-benda/monitor/internal/client"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/service"
	"github.com/artem-benda/monitor/internal/storage"
)

func main() {
	httpClient := &http.Client{Timeout: 30 * time.Second}

	for {
		var metrics map[model.Metric]string
		for i := 0; i < 5; i++ {
			metrics = service.ReadMetrics(storage.CounterStore)
			time.Sleep(2 * time.Second)
		}
		client.SendAllMetrics(httpClient, "http://localhost:8080/", metrics)
		storage.CounterStore.Reset()
	}
}
