package handlers

import (
	"net/http"

	"github.com/artem-benda/monitor/internal/dto"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/service"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/mailru/easyjson"
)

func MakeUpdateBatchJSONHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")

		dtos := make(dto.MetricsBatch, 0)

		if err := easyjson.UnmarshalFromReader(r.Body, &dtos); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		models := make([]model.MetricKeyWithValue, len(dtos))

		for _, dto := range dtos {

			switch {
			case model.ValidMetricKind(dto.MType) && dto.MType != "" && dto.ID != "":
				switch dto.MType {
				case model.GaugeKind:
					{
						if dto.Value != nil {
							models = append(models, model.MetricKeyWithValue{Kind: dto.MType, Name: dto.ID, Gauge: *dto.Value})
						} else {
							w.WriteHeader(http.StatusBadRequest)
							return
						}
					}
				case model.CounterKind:
					{
						if dto.Delta != nil {
							models = append(models, model.MetricKeyWithValue{Kind: dto.MType, Name: dto.ID, Counter: *dto.Delta})
						} else {
							w.WriteHeader(http.StatusBadRequest)
							return
						}
					}
				default:
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			case !model.ValidMetricKind(dto.MType):
				w.WriteHeader(http.StatusBadRequest)
				return
			case dto.ID == "":
				w.WriteHeader(http.StatusNotFound)
				return
			default:
				w.WriteHeader(http.StatusNotImplemented)
				return
			}
		}

		err := service.UpdateMetrics(r.Context(), store, models)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
