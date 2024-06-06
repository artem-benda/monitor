package handlers

import (
	"fmt"
	"net/http"

	"github.com/artem-benda/monitor/internal/server/service"
	"github.com/artem-benda/monitor/internal/server/storage"
)

// MakeGetAllHandler - создать обработчик метода получения списка всех актуальных значений метрик
func MakeGetAllHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/html")
		w.WriteHeader(http.StatusOK)
		metrics, err := service.GetAllMetrics(r.Context(), store)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for metricKey, strVal := range metrics {
			w.Write([]byte("<p>"))
			w.Write([]byte(fmt.Sprintf("%s: %s", metricKey.Name, strVal)))
			w.Write([]byte("</p>"))
		}
	}
}
