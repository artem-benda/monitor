package main

import (
	"context"
	"net/http"

	_ "net/http/pprof" // подключаем пакет pprof

	"github.com/artem-benda/monitor/internal/gzipper"
	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/retry"
	"github.com/artem-benda/monitor/internal/server/errors"
	"github.com/artem-benda/monitor/internal/server/handlers"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/artem-benda/monitor/internal/signer"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

var store storage.Storage
var dbpool *pgxpool.Pool

func main() {
	parseFlags()
	var err error

	if err = logger.Initialize(config.LogLevel); err != nil {
		panic(err)
	} else {
		defer logger.Log.Sync()
	}

	if config.DatabaseDSN != "" {
		dbpool = newConnectionPool(config.DatabaseDSN)
		initDB(dbpool)
		defer dbpool.Close()

		retryController := retry.NewRetryController(errors.ErrStorageConnection{})
		store = storage.NewDBStorage(dbpool, retryController)
	} else {
		store, err = storage.NewMemStorage(config.StoreIntervalSeconds, config.StoreFileName, config.StoreRestoreFromFile)
		if err != nil {
			panic(err)
		}
	}

	r := newAppRouter()
	err = http.ListenAndServe(config.Endpoint, r)
	if err != nil {
		panic(err)
	}
}

func newAppRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.LoggerMiddleware)
	r.Use(signer.CreateVerifyAndSignMiddleware([]byte(config.Key)))
	r.Use(gzipper.GzipMiddleware)

	r.Mount("/debug", middleware.Profiler())

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.MakeGetAllHandler(store))
		r.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.MakeUpdatePathHandler(store))
		r.Post("/update/", handlers.MakeUpdateJSONHandler(store))
		r.Post("/updates/", handlers.MakeUpdateBatchJSONHandler(store))
		r.Get("/ping", handlers.MakePingDatabaseHandler(dbpool))
		r.Get("/value/{metricType}/{metricName}", handlers.MakeGetHandler(store))
		r.Post("/value/", handlers.MakeGetJSONHandler(store))
	})
	return r
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
