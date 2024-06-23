package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	ServerEndpoint     string `env:"ADDRESS" json:"address"`
	LogLevel           string `env:"LOG_LEVEL"`
	Key                string `env:"KEY"`
	RSAPubKeyBase64    string `env:"CRYPTO_KEY" json:"crypto_key"`
	ReportInterval     int    `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval       int    `env:"POLL_INTERVAL" json:"poll_interval"`
	MaxParallelWorkers int    `env:"RATE_LIMIT"`
}

var config Config

func parseFlags() {
	// Проверим наличие флага -c/-config
	var configFilenameFlag string
	flag.StringVar(&configFilenameFlag, "c", "", "path to configuration file in JSON format")
	if configFilenameFlag == "" {
		flag.StringVar(&configFilenameFlag, "config", "", "path to configuration file in JSON format")
	}

	// Если задан путь к конфигурационному файлу - парсим его
	if configFilenameFlag != "" {
		fileBytes, err := os.ReadFile(configFilenameFlag)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(fileBytes, &config)
		if err != nil {
			panic(err)
		}
	}

	flag.StringVar(&config.ServerEndpoint, "a", "localhost:8080", "address and port of metrics server")
	flag.IntVar(&config.ReportInterval, "r", 10, "send metrics delay in seconds")
	flag.IntVar(&config.PollInterval, "p", 2, "poll runtime metrics delay in seconds")
	flag.StringVar(&config.LogLevel, "v", "debug", "logging level: debug, info, warn, error, dpanic, panic, fatal")
	flag.StringVar(&config.Key, "k", "", "if set, header with signature will be added to requests")
	flag.IntVar(&config.MaxParallelWorkers, "l", 2, "max parallel workers, to limit parallel metrics requests")
	flag.StringVar(&config.RSAPubKeyBase64, "crypto-key", "", "RSA base64 public key, used to encrypt request body, if set")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		panic(err)
	}
}
