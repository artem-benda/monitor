package main

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Endpoint             string `env:"ADDRESS"`
	LogLevel             string `env:"LOG_LEVEL"`
	StoreIntervalSeconds int    `env:"STORE_INTERVAL"`
	StoreFileName        string `env:"FILE_STORAGE_PATH"`
	StoreRestoreFromFile bool   `env:"RESTORE"`
}

var config Config

func parseFlags() {
	flag.StringVar(&config.Endpoint, "a", "localhost:8080", "address and port of metrics server")
	flag.StringVar(&config.LogLevel, "l", "debug", "logging level: debug, info, warn, error, dpanic, panic, fatal")
	flag.IntVar(&config.StoreIntervalSeconds, "i", 300, "Period in seconds to save current metrics into file")
	flag.StringVar(&config.StoreFileName, "f", "/tmp/metrics-db.json", "file path to save current metrics values to")
	flag.BoolVar(&config.StoreRestoreFromFile, "r", true, "should restore metrics values from file on startup")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		panic(err)
	}
}
