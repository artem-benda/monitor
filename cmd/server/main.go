package main

import (
	"context"

	_ "net/http/pprof" // подключаем пакет pprof

	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/retry"
	"github.com/artem-benda/monitor/internal/server/errors"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/artem-benda/monitor/internal/version"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	store        storage.Storage
	dbpool       *pgxpool.Pool
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	version.PrintVersion(buildVersion, buildDate, buildCommit)
	parseFlags()
	var err error

	if err = logger.Initialize(config.LogLevel); err != nil {
		panic(err)
	} else {
		defer func() {
			err = logger.Log.Sync()
			if err != nil {
				logger.Log.Error("Could not sync log data", zap.Error(err))
			}
		}()
	}

	var flushStorage func() error

	if config.DatabaseDSN != "" {
		dbpool = newConnectionPool(config.DatabaseDSN)
		initDB(dbpool)
		defer dbpool.Close()

		retryController := retry.NewRetryController(errors.ErrStorageConnection{})
		store = storage.NewDBStorage(dbpool, retryController)
	} else {
		store, flushStorage, err = storage.NewMemStorage(config.StoreIntervalSeconds, config.StoreFileName, config.StoreRestoreFromFile)
		if err != nil {
			panic(err)
		}
	}

	if config.UseGRPC {
		mustRunGrpcServer(store, dbpool, flushStorage)
	} else {
		mustRunRestServer(flushStorage)
	}
}

func initDB(dbpool *pgxpool.Pool) {
	createTblMetrics := `CREATE TABLE IF NOT EXISTS metrics(
		mtype text NOT NULL,
		mname text NOT NULL,
		gauge double precision,
		counter bigint,
		PRIMARY KEY (mtype, mname)
		)`
	_, err := dbpool.Exec(context.Background(), createTblMetrics)
	if err != nil {
		panic(err)
	}
}
