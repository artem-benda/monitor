package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/artem-benda/monitor/internal/service"
	"github.com/artem-benda/monitor/internal/storage"
)

func MakeGetAllHandler(store storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("GetHandler, method = %s, path = %s", r.Method, r.URL.Path)
		w.Header().Add("Content-type", "text/html")
		for metricKey, strVal := range service.GetAllMetrics(store) {
			w.Write([]byte("<p>"))
			w.Write([]byte(fmt.Sprintf("%s: %s", metricKey.Name, strVal)))
			w.Write([]byte("</p>"))
		}
	}
}
