package handlers

import (
	"fmt"
	"net/http"

	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/server/service"
	"github.com/artem-benda/monitor/internal/server/storage"
	"go.uber.org/zap"
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
			_, err = w.Write([]byte(fmt.Sprintf("<p>%s: %s</p>", metricKey.Name, strVal)))
			if err != nil {
				logger.Log.Error("Could not write body", zap.Error(err))
			}
		}
	}
}
