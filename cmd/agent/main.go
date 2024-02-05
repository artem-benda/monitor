package main

import (
	"time"

	"github.com/artem-benda/monitor/internal/client"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/service"
	"github.com/artem-benda/monitor/internal/storage"
	"github.com/go-resty/resty/v2"
)

func main() {
	resty := resty.New()
	resty.SetTimeout(30 * time.Second)

	for {
		var metrics map[model.Metric]string
		for i := 0; i < 5; i++ {
			metrics = service.ReadMetrics(storage.CounterStore)
			time.Sleep(2 * time.Second)
		}
		client.SendAllMetrics(resty, "http://localhost:8080/", metrics)
		storage.CounterStore.Reset()
	}
}
