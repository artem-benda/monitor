package service

import (
	"context"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/storage"
)

func GetAllMetrics(ctx context.Context, s storage.Storage) (map[model.MetricKey]string, error) {
	result := make(map[model.MetricKey]string)
	metrics, err := s.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	for key, val := range metrics {
		strVal, err := model.StringValue(key, val)
		if err != nil {
			return nil, err
		}
		result[key] = strVal
	}
	return result, nil
}
