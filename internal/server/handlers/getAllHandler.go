package handlers

import (
	"fmt"
	"net/http"

	"github.com/artem-benda/monitor/internal/server/service"
	"github.com/artem-benda/monitor/internal/server/storage"
)

func MakeGetAllHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/html")
		w.WriteHeader(http.StatusOK)
		for metricKey, strVal := range service.GetAllMetrics(store) {
			w.Write([]byte("<p>"))
			w.Write([]byte(fmt.Sprintf("%s: %s", metricKey.Name, strVal)))
			w.Write([]byte("</p>"))
		}
	}
}
