package handlers

import (
	"net/http"

	"github.com/artem-benda/monitor/internal/dto"
	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/service"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

// MakeGetJSONHandler - создать обработчик метода получения списка всех актуальных значений метрик в виде JSON
func MakeGetJSONHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")

		metrics := &dto.Metrics{}
		if err := easyjson.UnmarshalFromReader(r.Body, metrics); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch {
		case !model.ValidMetricKind(metrics.MType):
			{
				w.WriteHeader(http.StatusBadRequest)
			}
		case metrics.MType == model.GaugeKind && metrics.ID != "":
			{
				if floatVal, ok, err := service.GetGaugeMetric(r.Context(), store, metrics.ID); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				} else if ok {
					w.WriteHeader(http.StatusOK)
					metrics.Value = &floatVal
					_, _, err = easyjson.MarshalToHTTPResponseWriter(metrics, w)
					if err != nil {
						logger.Log.Error("Error writing json body", zap.Error(err))
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
		case metrics.MType == model.CounterKind && metrics.ID != "":
			{
				if intVal, ok, err := service.GetCounterMetric(r.Context(), store, metrics.ID); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				} else if ok {
					w.WriteHeader(http.StatusOK)
					metrics.Delta = &intVal
					_, _, err = easyjson.MarshalToHTTPResponseWriter(metrics, w)
					if err != nil {
						logger.Log.Error("Could not write json body", zap.Error(err))
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
