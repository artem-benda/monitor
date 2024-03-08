package handlers

import (
	"net/http"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/service"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func MakeGetHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
