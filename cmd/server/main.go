package main

import (
	"net/http"

	"github.com/artem-benda/monitor/internal/server"
	"github.com/artem-benda/monitor/internal/storage"
)

const (
	updatePath = "/update/"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle(updatePath, http.StripPrefix(updatePath, http.HandlerFunc(server.MakeUpdateHandler(storage.Store))))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
