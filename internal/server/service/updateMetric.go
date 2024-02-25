package service

import (
	"errors"
	"strconv"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/storage"
)

var (
	errIllegalArgument        = errors.New("illegal argument")
	errMetricKindNotSupported = errors.New("metric kind not supported")
)

func UpdateMetric(s storage.Storage, kind string, name string, strVal string) error {
	switch kind {
	case model.CounterKind:
		{
			if value, err := strconv.ParseInt(strVal, 10, 64); err == nil {
				UpdateCounterMetric(s, name, value)
			} else {
				return err
			}
		}
	case model.GaugeKind:
		{
			if value, err := strconv.ParseFloat(strVal, 64); err == nil {
				UpdateGaugeMetric(s, name, value)
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

func UpdateGaugeMetric(s storage.Storage, name string, value float64) {
	key := model.Metric{Kind: model.GaugeKind, Name: name}
	s.Put(key, value)
}

func UpdateCounterMetric(s storage.Storage, name string, value int64) {
	key := model.Metric{Kind: model.CounterKind, Name: name}
	s.UpdateFunc(key, func(prev any) any {
		if prevInt, ok := prev.(int64); ok {
			return prevInt + int64(value)
		} else if prev == nil {
			return value
		} else {
			return errIllegalArgument
		}
	})
}
