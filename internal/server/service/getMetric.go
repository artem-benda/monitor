package service

import (
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/storage"
)

func GetMetric(storage storage.Storage, kind string, name string) (string, bool) {
	key := model.Metric{Kind: kind, Name: name}
	if val, ok := storage.Get(key); ok {
		return model.StringValue(key, val)
	}
	return "", false
}

func GetGaugeMetric(storage storage.Storage, name string) (float64, bool) {
	key := model.Metric{Kind: model.GaugeKind, Name: name}
	if val, ok := storage.Get(key); ok {
		return model.AsGaugeMetric(key, val)
	}
	return 0, false
}

func GetCounterMetric(storage storage.Storage, name string) (int64, bool) {
	key := model.Metric{Kind: model.CounterKind, Name: name}
	if val, ok := storage.Get(key); ok {
		return model.AsCounterMetric(key, val)
	}
	return 0, false
}
