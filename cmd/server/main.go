package main

import (
	"net/http"

	"github.com/artem-benda/monitor/internal/server"
)

const (
	updatePath = "/update/"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle(updatePath, http.StripPrefix(updatePath, http.HandlerFunc(server.UpdateHandler)))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
