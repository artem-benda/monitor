package main

import (
	"fmt"
	"time"

	"github.com/artem-benda/monitor/internal/client/requests"
	"github.com/artem-benda/monitor/internal/client/service"
	"github.com/artem-benda/monitor/internal/client/storage"
	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/go-resty/resty/v2"
)

func main() {
	parseFlags()

	if err := logger.Initialize(config.LogLevel); err != nil {
		panic(err)
	} else {
		defer logger.Log.Sync()
	}

	resty := resty.New()
	serverEndpointURL := fmt.Sprintf("http://%s", config.ServerEndpoint)
	resty.SetBaseURL(serverEndpointURL)
	resty.SetTimeout(30 * time.Second)

	var metrics map[model.Metric]any

	go func() {
		for {
			metrics = service.ReadMetrics(storage.CounterStore)
			time.Sleep(time.Duration(config.PollInterval) * time.Second)
		}
	}()

	for {
		requests.SendAllMetrics(resty, metrics)
		storage.CounterStore.Reset()
		time.Sleep(time.Duration(config.ReportInterval) * time.Second)
	}
}
