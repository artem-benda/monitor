package main

import (
	"net/http"

	"github.com/artem-benda/monitor/internal/gzipper"
	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/server/handlers"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	parseFlags()

	if err := logger.Initialize(config.LogLevel); err != nil {
		panic(err)
	} else {
		defer logger.Log.Sync()
	}

	r := newAppRouter()
	err := http.ListenAndServe(config.Endpoint, r)
	if err != nil {
		panic(err)
	}
}

func newAppRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.LoggerMiddleware)
	r.Use(gzipper.GzipMiddleware)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.MakeUpdatePathHandler(storage.Store))
	r.Post("/update/", handlers.MakeUpdateJSONHandler(storage.Store))
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.MakeGetAllHandler(storage.Store))
		r.Get("/value/{metricType}/{metricName}", handlers.MakeGetHandler(storage.Store))
		r.Post("/value/", handlers.MakeGetJSONHandler(storage.Store))
	})
	return r
}
