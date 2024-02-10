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
				key := model.Metric{Kind: kind, Name: name}
				s.UpdateFunc(key, func(prev any) any {
					if prevInt, ok := prev.(int64); ok {
						return prevInt + int64(value)
					} else if prev == nil {
						return value
					} else {
						return errIllegalArgument
					}
				})
			} else {
				return err
			}
		}
	case model.GaugeKind:
		{
			if value, err := strconv.ParseFloat(strVal, 64); err == nil {
				key := model.Metric{Kind: kind, Name: name}
				s.Put(key, value)
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
