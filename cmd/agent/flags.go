package main

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

var (
	serverEndpoint string
	reportInterval int
	pollInterval   int
)

func parseFlags() {
	flag.StringVar(&serverEndpoint, "a", "localhost:8080", "address and port of metrics server")
	flag.IntVar(&reportInterval, "r", 10, "send metrics delay in seconds")
	flag.IntVar(&pollInterval, "p", 2, "poll runtime metrics delay in seconds")
	flag.Parse()

	var envConfig struct {
		ServerEndpoint string `env:"ADDRESS"`
		ReportInterval int    `env:"REPORT_INTERVAL"`
		PollInterval   int    `env:"POLL_INTERVAL"`
	}

	env.Parse(&envConfig)
	if envConfig.ServerEndpoint != "" {
		serverEndpoint = envConfig.ServerEndpoint
	}
	if envConfig.ReportInterval != 0 {
		reportInterval = envConfig.ReportInterval
	}
	if envConfig.PollInterval != 0 {
		pollInterval = envConfig.PollInterval
	}
}
