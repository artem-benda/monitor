package main

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	ServerEndpoint string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	LogLevel       string `env:"LOG_LEVEL"`
}

var config Config

func parseFlags() {
	flag.StringVar(&config.ServerEndpoint, "a", "localhost:8080", "address and port of metrics server")
	flag.IntVar(&config.ReportInterval, "r", 10, "send metrics delay in seconds")
	flag.IntVar(&config.PollInterval, "p", 2, "poll runtime metrics delay in seconds")
	flag.StringVar(&config.LogLevel, "l", "debug", "logging level: debug, info, warn, error, dpanic, panic, fatal")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		panic(err)
	}
}
