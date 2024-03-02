package service

import (
	"errors"
	"strconv"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/storage"
)

var (
	errMetricKindNotSupported = errors.New("metric kind not supported")
)

func UpdateMetric(s storage.Storage, kind string, name string, strVal string) error {
	switch kind {
	case model.CounterKind:
		{
			if value, err := strconv.ParseInt(strVal, 10, 64); err == nil {
				_, err := UpdateAndGetCounterMetric(s, name, value)
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
				_, err := UpdateAndGetGaugeMetric(s, name, value)
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

func UpdateAndGetGaugeMetric(s storage.Storage, name string, value float64) (float64, error) {
	key := model.Metric{Kind: model.GaugeKind, Name: name}
	err := s.Put(key, value)
	return value, err
}

func UpdateAndGetCounterMetric(s storage.Storage, name string, value int64) (int64, error) {
	key := model.Metric{Kind: model.CounterKind, Name: name}
	next, err := s.UpdateAndGetFunc(key, func(prev any) any {
		if prevInt, ok := prev.(int64); ok {
			return prevInt + int64(value)
		}
		return value
	})
	return next.(int64), err
}
