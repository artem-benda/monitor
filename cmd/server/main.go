package main

import (
	"net/http"

	"github.com/artem-benda/monitor/internal/server/handlers"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	parseFlags()

	if logger, err := zap.NewDevelopment(); err != nil {
		panic(err)
	} else {
		defer logger.Sync()
	}

	r := newAppRouter()
	err := http.ListenAndServe(config.Endpoint, r)
	if err != nil {
		panic(err)
	}
}

func newAppRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.MakeUpdateHandler(storage.Store))
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.MakeGetAllHandler(storage.Store))
		r.Get("/value/{metricType}/{metricName}", handlers.MakeGetHandler(storage.Store))
	})
	return r
}
