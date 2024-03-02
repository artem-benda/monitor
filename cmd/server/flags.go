package main

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Endpoint string `env:"ADDRESS"`
	LogLevel string `env:"LOG_LEVEL"`
}

var config Config

func parseFlags() {
	flag.StringVar(&config.Endpoint, "a", "localhost:8080", "address and port of metrics server")
	flag.StringVar(&config.LogLevel, "l", "debug", "logging level: debug, info, warn, error, dpanic, panic, fatal")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		panic(err)
	}
}
