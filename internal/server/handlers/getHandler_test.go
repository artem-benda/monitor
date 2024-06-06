package handlers

import (
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func ExampleMakeGetHandler() {
	memStorage, err := storage.NewMemStorage(100, "", false)
	if err != nil {
		panic(err)
	}

	handler := MakeGetHandler(memStorage)

	r := chi.NewRouter()
	r.Get("/value/{metricType}/{metricName}", handler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
