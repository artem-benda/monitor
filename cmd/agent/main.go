package main

import (
	"fmt"
	"time"

	"github.com/artem-benda/monitor/internal/client"
	"github.com/artem-benda/monitor/internal/client/service"
	"github.com/artem-benda/monitor/internal/client/storage"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/go-resty/resty/v2"
)

func main() {
	parseFlags()
	resty := resty.New()
	resty.SetTimeout(30 * time.Second)

	serverEndpointURL := fmt.Sprintf("http://%s", serverEndpoint)

	var metrics map[model.Metric]string

	go func() {
		for {
			metrics = service.ReadMetrics(storage.CounterStore)
			time.Sleep(time.Duration(pollInterval) * time.Second)
		}
	}()

	for {
		client.SendAllMetrics(resty, serverEndpointURL, metrics)
		storage.CounterStore.Reset()
		time.Sleep(time.Duration(reportInterval) * time.Second)
	}
}
