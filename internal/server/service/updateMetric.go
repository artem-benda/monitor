package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/storage"
)

var (
	errMetricKindNotSupported = errors.New("metric kind not supported")
)

func UpdateMetric(ctx context.Context, s storage.Storage, kind string, name string, strVal string) error {
	switch kind {
	case model.CounterKind:
		{
			if value, err := strconv.ParseInt(strVal, 10, 64); err == nil {
				_, err := UpdateAndGetCounterMetric(ctx, s, name, value)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	case model.GaugeKind:
		{
			if value, err := strconv.ParseFloat(strVal, 64); err == nil {
				_, err := UpdateAndGetGaugeMetric(ctx, s, name, value)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	default:
		{
			return errMetricKindNotSupported
		}
	}
	return nil
}

func UpdateAndGetGaugeMetric(ctx context.Context, s storage.Storage, name string, value float64) (float64, error) {
	key := model.MetricKey{Kind: model.GaugeKind, Name: name}
	err := s.UpsertGauge(ctx, key, model.MetricValue{Gauge: value})
	return value, err
}

func UpdateAndGetCounterMetric(ctx context.Context, s storage.Storage, name string, value int64) (int64, error) {
	key := model.MetricKey{Kind: model.CounterKind, Name: name}
	return s.UpsertCounterAndGet(ctx, key, value)
}
