package main

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Endpoint string `env:"ADDRESS"`
}

var config Config

func parseFlags() {
	flag.StringVar(&config.Endpoint, "a", "localhost:8080", "address and port of metrics server")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		panic(err)
	}
}
