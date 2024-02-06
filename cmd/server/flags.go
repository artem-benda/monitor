package main

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

var endpoint string

func parseFlags() {
	flag.StringVar(&endpoint, "a", "localhost:8080", "address and port of metrics server")
	flag.Parse()

	var envConfig struct {
		Endpoint string `env:"ADDRESS"`
	}

	if err := env.Parse(&envConfig); err != nil {
		panic(err)
	}
	if envConfig.Endpoint != "" {
		endpoint = envConfig.Endpoint
	}
}
