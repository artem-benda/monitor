package handlers

import (
	"log"
	"net/http"

	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func ExampleMakeGetHandler() {
	memStorage, _, err := storage.NewMemStorage(100, "", false)
	if err != nil {
		panic(err)
	}

	handler := MakeGetHandler(memStorage)

	r := chi.NewRouter()
	r.Get("/value/{metricType}/{metricName}", handler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
