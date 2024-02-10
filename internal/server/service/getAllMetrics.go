package service

import (
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/storage"
)

func GetAllMetrics(storage storage.Storage) map[model.Metric]string {
	result := make(map[model.Metric]string)
	for key, val := range storage.GetAll() {
		if strVal, ok := model.StringValue(key, val); ok {
			result[key] = strVal
		}
	}
	return result
}
