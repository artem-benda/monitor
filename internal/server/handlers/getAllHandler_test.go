package handlers

import (
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func ExampleMakeGetAllHandler() {
	memStorage, err := storage.NewMemStorage(100, "", false)
	if err != nil {
		panic(err)
	}

	handler := MakeGetAllHandler(memStorage)

	r := chi.NewRouter()
	r.Get("/", handler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
