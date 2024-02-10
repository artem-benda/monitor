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
