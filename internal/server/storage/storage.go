package storage

import (
	"context"
	"errors"

	"github.com/artem-benda/monitor/internal/model"
)

type Storage interface {
	Get(ctx context.Context, key model.MetricKey) (model.MetricValue, bool, error)
	UpsertGauge(ctx context.Context, key model.MetricKey, value model.MetricValue) error
	UpsertCounterAndGet(ctx context.Context, key model.MetricKey, incCounter int64) (int64, error)
	GetAll(ctx context.Context) (map[model.MetricKey]model.MetricValue, error)
}

var (
	errInvaligArgument  = errors.New("invalid argument")
	errNullCounterValue = errors.New("null counter value")
	errInvalidData      = errors.New("invalid metric value for type")
)
