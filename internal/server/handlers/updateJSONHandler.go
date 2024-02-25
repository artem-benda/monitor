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

func MakeUpdateJSONHandler(store storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("UpdateJSONHandler, method = %s, path = %s", r.Method, r.URL.Path)
		w.Header().Add("Content-type", "application/json")

		metrics := &dto.Metrics{}
		if err := easyjson.UnmarshalFromReader(r.Body, metrics); err != nil {
			logger.Log.Info("Error parsing request body", zap.Error(err))
			http.Error(w, "Error parsing request body", http.StatusBadRequest)
		}

		switch {
		case model.ValidMetricKind(metrics.MType) && metrics.MType != "" && metrics.ID != "":
			switch metrics.MType {
			case model.GaugeKind:
				{
					if metrics.Value != nil {
						service.UpdateGaugeMetric(store, metrics.ID, *metrics.Value)
						w.WriteHeader(http.StatusOK)
					} else {
						http.Error(w, "Metric value not set", http.StatusBadRequest)
					}
				}
			case model.CounterKind:
				{
					if metrics.Delta != nil {
						service.UpdateCounterMetric(store, metrics.ID, *metrics.Delta)
						w.WriteHeader(http.StatusOK)
					} else {
						http.Error(w, "Metric value not set", http.StatusBadRequest)
					}
				}
			default:
				http.Error(w, "Bad metric type", http.StatusBadRequest)
			}
		case !model.ValidMetricKind(metrics.MType):
			http.Error(w, "Metric type not supported", http.StatusBadRequest)
		case metrics.ID == "":
			http.Error(w, "Metric name cannot be empty", http.StatusNotFound)
		default:
			http.Error(w, "Method unimplemented", http.StatusNotImplemented)
		}
	}

	return logger.WithRequestLogger(handlerFunc)
}
