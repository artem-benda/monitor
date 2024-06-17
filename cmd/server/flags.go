package main

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Endpoint             string `env:"ADDRESS"`
	LogLevel             string `env:"LOG_LEVEL"`
	StoreFileName        string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN          string `env:"DATABASE_DSN"`
	Key                  string `env:"KEY"`
	RSAPrivKeyBase64     string `env:"CRYPTO_KEY"`
	StoreIntervalSeconds int    `env:"STORE_INTERVAL"`
	StoreRestoreFromFile bool   `env:"RESTORE"`
}

var config Config

func parseFlags() {
	flag.StringVar(&config.Endpoint, "a", "localhost:8080", "address and port of metrics server")
	flag.StringVar(&config.LogLevel, "l", "debug", "logging level: debug, info, warn, error, dpanic, panic, fatal")
	flag.IntVar(&config.StoreIntervalSeconds, "i", 300, "Period in seconds to save current metrics into file")
	flag.StringVar(&config.StoreFileName, "f", "/tmp/metrics-db.json", "file path to save current metrics values to")
	flag.BoolVar(&config.StoreRestoreFromFile, "r", true, "should restore metrics values from file on startup")
	flag.StringVar(&config.DatabaseDSN, "d", "", "Database connection URL in pgx format, for ex. postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10")
	flag.StringVar(&config.Key, "k", "", "if set, signature in header for POST requests will be validated")
	flag.StringVar(&config.RSAPrivKeyBase64, "crypto-key", "", "RSA base64 private key, used to decrypt agent's request body, if set")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		panic(err)
	}
}
