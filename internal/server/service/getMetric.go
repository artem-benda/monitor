package service

import (
	"context"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/storage"
)

func GetMetric(ctx context.Context, s storage.Storage, kind string, name string) (string, bool, error) {
	key := model.MetricKey{Kind: kind, Name: name}
	val, ok, err := s.Get(ctx, key)

	if err != nil {
		return "", false, err
	}

	if !ok {
		return "", false, nil
	}

	stringValue, err := model.StringValue(key, *val)
	if err != nil {
		return "", false, err
	}

	return stringValue, true, nil
}

func GetGaugeMetric(ctx context.Context, storage storage.Storage, name string) (float64, bool, error) {
	key := model.MetricKey{Kind: model.GaugeKind, Name: name}
	val, ok, err := storage.Get(ctx, key)
	if err != nil {
		return 0, false, err
	}
	if ok {
		return val.Gauge, true, nil
	}
	return 0, false, nil
}

func GetCounterMetric(ctx context.Context, storage storage.Storage, name string) (int64, bool, error) {
	key := model.MetricKey{Kind: model.CounterKind, Name: name}
	val, ok, err := storage.Get(ctx, key)
	if err != nil {
		return 0, false, err
	}
	if ok {
		return val.Counter, true, nil
	}
	return 0, false, nil
}
