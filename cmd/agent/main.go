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
	"go.uber.org/zap"
)

func main() {
	parseFlags()

	if err := logger.Initialize(config.LogLevel); err != nil {
		panic(err)
	} else {
		defer logger.Log.Sync()
	}

	client := resty.New()
	serverEndpointURL := fmt.Sprintf("http://%s", config.ServerEndpoint)
	client.SetBaseURL(serverEndpointURL)
	client.SetTimeout(30 * time.Second)
	client.OnAfterResponse(logger.NewRestyResponseLogger())

	var metrics map[model.MetricKey]model.MetricValue

	go func() {
		for {
			metrics = service.ReadMetrics(storage.CounterStore)
			time.Sleep(time.Duration(config.PollInterval) * time.Second)
		}
	}()

	for {
		err := requests.SendAllMetrics(client, metrics)

		if err != nil {
			logger.Log.Debug("error sending metrics batch", zap.Error(err))
		}

		storage.CounterStore.Reset()
		time.Sleep(time.Duration(config.ReportInterval) * time.Second)
	}
}
