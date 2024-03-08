package handlers

import (
	"log"
	"net/http"

	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/service"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func MakeGetHandler(store storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("GetHandler, method = %s, path = %s", r.Method, r.URL.Path)
		w.Header().Add("Content-type", "text/plain")
		metricType, metricName := chi.URLParam(r, "metricType"), chi.URLParam(r, "metricName")
		switch {
		case !model.ValidMetricKind(metricType):
			{
				http.Error(w, "Metric type not supported", http.StatusBadRequest)
			}
		case model.ValidMetricKind(metricType) && metricName != "":
			{
				if strVal, ok := service.GetMetric(store, metricType, metricName); ok {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(strVal))
				} else {
					http.Error(w, "", http.StatusNotFound)
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
