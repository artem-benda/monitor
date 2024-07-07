package main

import (
	"fmt"
	"time"

	"github.com/artem-benda/monitor/internal/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func mustCreateRestyClient() *resty.Client {
	client := resty.New()
	serverEndpointURL := fmt.Sprintf("http://%s", config.ServerEndpoint)

	logger.Log.Debug("Starting with base URL", zap.String("baseURL", serverEndpointURL))

	client.SetBaseURL(serverEndpointURL)
	client.SetTimeout(30 * time.Second)
	client.OnAfterResponse(logger.NewRestyResponseLogger())
	client.Header.Add("X-Real-IP", mustGetLocalIPAddr(config.ServerEndpoint).String())
}
