package main

import (
	"net/http"

	"github.com/artem-benda/monitor/internal/gzipper"
	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/server/handlers"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/go-chi/chi/v5"
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
	}

	store, err = storage.NewStorage(config.StoreIntervalSeconds, config.StoreFileName, config.StoreRestoreFromFile)
	if err != nil {
		panic(err)
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
	r.Use(gzipper.GzipMiddleware)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.MakeUpdatePathHandler(store))
	r.Post("/update/", handlers.MakeUpdateJSONHandler(store))
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.MakeGetAllHandler(store))
		r.Get("/value/{metricType}/{metricName}", handlers.MakeGetHandler(store))
		r.Post("/value/", handlers.MakeGetJSONHandler(store))
		r.Get("/ping/", handlers.MakePingDatabaseHandler(dbpool))
	})
	return r
}
