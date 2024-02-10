package storage

import "github.com/artem-benda/monitor/internal/model"

type Storage interface {
	Get(key model.Metric) (any, bool)
	Put(key model.Metric, value any)
	UpdateFunc(key model.Metric, fn func(prev any) any)
}
