package main

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	ServerEndpoint     string `env:"ADDRESS"`
	LogLevel           string `env:"LOG_LEVEL"`
	Key                string `env:"KEY"`
	ReportInterval     int    `env:"REPORT_INTERVAL"`
	PollInterval       int    `env:"POLL_INTERVAL"`
	MaxParallelWorkers int    `env:"RATE_LIMIT"`
}

var config Config

func parseFlags() {
	flag.StringVar(&config.ServerEndpoint, "a", "localhost:8080", "address and port of metrics server")
	flag.IntVar(&config.ReportInterval, "r", 10, "send metrics delay in seconds")
	flag.IntVar(&config.PollInterval, "p", 2, "poll runtime metrics delay in seconds")
	flag.StringVar(&config.LogLevel, "v", "debug", "logging level: debug, info, warn, error, dpanic, panic, fatal")
	flag.StringVar(&config.Key, "k", "", "if set, header with signature will be added to requests")
	flag.IntVar(&config.MaxParallelWorkers, "l", 2, "max parallel workers, to limit parallel metrics requests")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		panic(err)
	}
}
