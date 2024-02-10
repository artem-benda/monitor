package main

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	ServerEndpoint string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

var config Config

func parseFlags() {
	flag.StringVar(&config.ServerEndpoint, "a", "localhost:8080", "address and port of metrics server")
	flag.IntVar(&config.ReportInterval, "r", 10, "send metrics delay in seconds")
	flag.IntVar(&config.PollInterval, "p", 2, "poll runtime metrics delay in seconds")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		panic(err)
	}
}
