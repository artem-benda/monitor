package service

import (
	"context"
	"testing"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetric(t *testing.T) {
	type metric struct {
		kind   string
		name   string
		strVal string
	}
	tests := []struct {
		name    string
		metrics []metric
		want    map[model.MetricKey]model.MetricValue
		wantErr bool
	}{
		{
			name:    "Initial counter update",
			metrics: []metric{{kind: "counter", name: "test", strVal: "1"}},
			want: map[model.MetricKey]model.MetricValue{
				{Kind: "counter", Name: "test"}: {Counter: int64(1)},
			},
			wantErr: false,
		},
		{
			name:    "Initial gauge update",
			metrics: []metric{{kind: "gauge", name: "test", strVal: "1.123456"}},
			want: map[model.MetricKey]model.MetricValue{
				{Kind: "gauge", Name: "test"}: {Gauge: float64(1.123456)},
			},
			wantErr: false,
		},
		{
			name: "Double counter update",
			metrics: []metric{
				{kind: "counter", name: "test", strVal: "1"},
				{kind: "counter", name: "test", strVal: "3"},
			},
			want: map[model.MetricKey]model.MetricValue{
				{Kind: "counter", Name: "test"}: {Counter: int64(4)},
			},
			wantErr: false,
		},
		{
			name: "Double gauge update",
			metrics: []metric{
				{kind: "gauge", name: "test", strVal: "1.123456"},
				{kind: "gauge", name: "test", strVal: "1.654321"},
			},
			want: map[model.MetricKey]model.MetricValue{
				{Kind: "gauge", Name: "test"}: {Gauge: float64(1.654321)},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, err := storage.NewMemStorage(10000, "test.txt", false)
			assert.NoError(t, err, "Error creating store for test")
			for _, metric := range tt.metrics {
				err := UpdateMetric(context.Background(), store, metric.kind, metric.name, metric.strVal)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
			for wantMetric, wantValue := range tt.want {
				value, _, _ := store.Get(context.Background(), wantMetric)
				assert.Equal(t, wantValue, *value)
			}
		})
	}
}
