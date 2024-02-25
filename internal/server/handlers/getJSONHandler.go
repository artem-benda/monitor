package handlers

import (
	"log"
	"net/http"

	"github.com/artem-benda/monitor/internal/dto"
	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/service"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

func MakeGetJSONHandler(store storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("GetJSONHandler, method = %s, path = %s", r.Method, r.URL.Path)
		w.Header().Add("Content-type", "application/json")

		metrics := &dto.Metrics{}
		if err := easyjson.UnmarshalFromReader(r.Body, metrics); err != nil {
			logger.Log.Info("Error parsing request body", zap.Error(err))
			http.Error(w, "Error parsing request body", http.StatusBadRequest)
		}

		switch {
		case !model.ValidMetricKind(metrics.MType):
			{
				http.Error(w, "Metric type not supported", http.StatusBadRequest)
			}
		case metrics.MType == model.GaugeKind && metrics.ID != "":
			{
				if floatVal, ok := service.GetGaugeMetric(store, metrics.ID); ok {
					w.WriteHeader(http.StatusOK)
					metrics.Value = &floatVal
					easyjson.MarshalToHTTPResponseWriter(metrics, w)
				} else {
					http.Error(w, "{}", http.StatusNotFound)
				}
			}
		case metrics.MType == model.CounterKind && metrics.ID != "":
			{
				if intVal, ok := service.GetCounterMetric(store, metrics.ID); ok {
					w.WriteHeader(http.StatusOK)
					metrics.Delta = &intVal
					easyjson.MarshalToHTTPResponseWriter(metrics, w)
				} else {
					http.Error(w, "{}", http.StatusNotFound)
				}
			}
		default:
			{
				http.Error(w, "", http.StatusNotFound)
			}
		}
	}

	return logger.WithRequestLogger(handlerFunc)
}
