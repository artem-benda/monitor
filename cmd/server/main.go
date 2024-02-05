package main

import (
	"net/http"

	"github.com/artem-benda/monitor/internal/server"
	"github.com/artem-benda/monitor/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", server.MakeUpdateHandler(storage.Store))
	r.Route("/", func(r chi.Router) {
		r.Get("/", server.MakeGetAllHandler(storage.Store))
		r.Get("/value/{metricType}/{metricName}", server.MakeGetHandler(storage.Store))
	})
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
