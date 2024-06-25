package handlers

import (
	"log"
	"net/http"

	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func ExampleMakeUpdateBatchJSONHandler() {
	memStorage, _, err := storage.NewMemStorage(100, "", false)
	if err != nil {
		panic(err)
	}

	handler := MakeUpdateBatchJSONHandler(memStorage)

	r := chi.NewRouter()
	r.Post("/", handler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
