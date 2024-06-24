package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Endpoint             string `env:"ADDRESS" json:"address"`
	LogLevel             string `env:"LOG_LEVEL"`
	StoreFileName        string `env:"FILE_STORAGE_PATH" json:"store_file"`
	DatabaseDSN          string `env:"DATABASE_DSN" json:"database_dsn"`
	Key                  string `env:"KEY"`
	RSAPrivKeyBase64     string `env:"CRYPTO_KEY" json:"crypto_key"`
	StoreIntervalSeconds int    `env:"STORE_INTERVAL" json:"store_interval"`
	StoreRestoreFromFile bool   `env:"RESTORE" json:"restore"`
}

var config Config

func parseFlags() {
	// Проверим наличие флага -c/-config
	var configFilenameFlag string
	cmdArgs := os.Args[1:]
	for ind, arg := range cmdArgs {
		// если текущий аргумент - ключ - проверим, что есть следующий аргумент и считаем его путем к конфигурационному файлу
		if arg == "-c" || arg == "-config" && ind+1 < len(cmdArgs) {
			configFilenameFlag = cmdArgs[ind+1]
		}
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

	flag.StringVar(&config.Endpoint, "a", "localhost:8080", "address and port of metrics server")
	flag.StringVar(&config.LogLevel, "l", "debug", "logging level: debug, info, warn, error, dpanic, panic, fatal")
	flag.IntVar(&config.StoreIntervalSeconds, "i", 300, "Period in seconds to save current metrics into file")
	flag.StringVar(&config.StoreFileName, "f", "/tmp/metrics-db.json", "file path to save current metrics values to")
	flag.BoolVar(&config.StoreRestoreFromFile, "r", true, "should restore metrics values from file on startup")
	flag.StringVar(&config.DatabaseDSN, "d", "", "Database connection URL in pgx format, for ex. postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10")
	flag.StringVar(&config.Key, "k", "", "if set, signature in header for POST requests will be validated")
	flag.StringVar(&config.RSAPrivKeyBase64, "crypto-key", "", "RSA base64 private key, used to decrypt agent's request body, if set")
	flag.StringVar(&configFilenameFlag, "c", "", "path to configuration file in JSON format")
	flag.StringVar(&configFilenameFlag, "config", "", "path to configuration file in JSON format")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		panic(err)
	}
}
