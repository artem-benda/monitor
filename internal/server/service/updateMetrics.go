package service

import (
	"context"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/storage"
)

func UpdateMetrics(ctx context.Context, s storage.Storage, metrics []model.MetricKeyWithValue) error {
	return s.UpsertBatch(ctx, metrics)
}
