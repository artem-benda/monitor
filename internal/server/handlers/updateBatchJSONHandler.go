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

func MakeUpdateBatchJSONHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")

		dtos := make(dto.MetricsBatch, 0)

		if err := easyjson.UnmarshalFromReader(r.Body, &dtos); err != nil {
			logger.Log.Debug("error unmarshalling dtos", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(dtos) == 0 {
			w.WriteHeader(http.StatusOK)
			return
		}

		models := make([]model.MetricKeyWithValue, len(dtos))

		for _, dto := range dtos {

			switch {
			case !model.ValidMetricKind(dto.MType):
				w.WriteHeader(http.StatusBadRequest)
				return
			case dto.ID == "":
				w.WriteHeader(http.StatusNotFound)
				return
			case dto.MType == "":
				w.WriteHeader(http.StatusNotImplemented)
				return
			case dto.MType == model.GaugeKind && dto.Value == nil:
				w.WriteHeader(http.StatusBadRequest)
				return
			case dto.MType == model.CounterKind && dto.Delta == nil:
				w.WriteHeader(http.StatusBadRequest)
				return
			case dto.MType == model.GaugeKind:
				models = append(models, model.MetricKeyWithValue{Kind: dto.MType, Name: dto.ID, Gauge: *dto.Value})
			case dto.MType == model.CounterKind:
				models = append(models, model.MetricKeyWithValue{Kind: dto.MType, Name: dto.ID, Counter: *dto.Delta})
			default:
				w.WriteHeader(http.StatusBadRequest)
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
