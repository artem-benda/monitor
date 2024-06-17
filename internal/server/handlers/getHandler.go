package handlers

import (
	"net/http"

	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/service"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// MakeGetHandler - создать обработчик метода получения значения указанной метрики
func MakeGetHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/plain")
		metricType, metricName := chi.URLParam(r, "metricType"), chi.URLParam(r, "metricName")
		switch {
		case !model.ValidMetricKind(metricType):
			{
				w.WriteHeader(http.StatusBadRequest)
			}
		case model.ValidMetricKind(metricType) && metricName != "":
			{
				if strVal, ok, err := service.GetMetric(r.Context(), store, metricType, metricName); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				} else if ok {
					w.WriteHeader(http.StatusOK)
					_, err = w.Write([]byte(strVal))
					if err != nil {
						logger.Log.Error("Could not write body", zap.Error(err))
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
		default:
			{
				w.WriteHeader(http.StatusNotFound)
			}
		}
	}
}
