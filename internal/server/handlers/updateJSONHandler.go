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
		case model.ValidMetricKind(metrics.MType) && metrics.MType != "" && metrics.ID != "":
			switch metrics.MType {
			case model.GaugeKind:
				{
					if metrics.Value != nil {
						var err error
						*metrics.Value, err = service.UpdateAndGetGaugeMetric(r.Context(), store, metrics.ID, *metrics.Value)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
						} else {
							w.WriteHeader(http.StatusOK)
							easyjson.MarshalToHTTPResponseWriter(metrics, w)
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
					}
				}
			case model.CounterKind:
				{
					if metrics.Delta != nil {
						var err error
						*metrics.Delta, err = service.UpdateAndGetCounterMetric(r.Context(), store, metrics.ID, *metrics.Delta)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
						} else {
							w.WriteHeader(http.StatusOK)
							easyjson.MarshalToHTTPResponseWriter(metrics, w)
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
					}
				}
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
		case !model.ValidMetricKind(metrics.MType):
			w.WriteHeader(http.StatusBadRequest)
		case metrics.ID == "":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}
