package handlers

import (
	"net/http"

	"github.com/artem-benda/monitor/internal/dto"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/service"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/mailru/easyjson"
)

func MakeUpdateJSONHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")

		metrics := &dto.Metrics{}
		if err := easyjson.UnmarshalFromReader(r.Body, metrics); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch {
		case !model.ValidMetricKind(metrics.MType):
			w.WriteHeader(http.StatusBadRequest)
			return
		case metrics.ID == "":
			w.WriteHeader(http.StatusNotFound)
			return
		case metrics.MType == "":
			w.WriteHeader(http.StatusNotImplemented)
			return
		case metrics.MType == model.GaugeKind && metrics.Value == nil:
			w.WriteHeader(http.StatusBadRequest)
			return
		case metrics.MType == model.CounterKind && metrics.Delta == nil:
			w.WriteHeader(http.StatusBadRequest)
			return
		case metrics.MType == model.GaugeKind:
			var err error
			*metrics.Value, err = service.UpdateAndGetGaugeMetric(r.Context(), store, metrics.ID, *metrics.Value)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
				easyjson.MarshalToHTTPResponseWriter(metrics, w)
			}
		case metrics.MType == model.CounterKind:
			var err error
			*metrics.Delta, err = service.UpdateAndGetCounterMetric(r.Context(), store, metrics.ID, *metrics.Delta)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
				easyjson.MarshalToHTTPResponseWriter(metrics, w)
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
