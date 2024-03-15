package handlers

import (
	"net/http"

	"github.com/artem-benda/monitor/internal/dto"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/service"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/mailru/easyjson"
)

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
					easyjson.MarshalToHTTPResponseWriter(metrics, w)
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
					easyjson.MarshalToHTTPResponseWriter(metrics, w)
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
