package service

import (
	"testing"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/storage"
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
		want    map[model.Metric]any
		wantErr bool
	}{
		{
			name:    "Initial counter update",
			metrics: []metric{{kind: "counter", name: "test", strVal: "1"}},
			want: map[model.Metric]any{
				{Kind: "counter", Name: "test"}: int64(1),
			},
			wantErr: false,
		},
		{
			name:    "Initial gauge update",
			metrics: []metric{{kind: "gauge", name: "test", strVal: "1.123456"}},
			want: map[model.Metric]any{
				{Kind: "gauge", Name: "test"}: float64(1.123456),
			},
			wantErr: false,
		},
		{
			name: "Double counter update",
			metrics: []metric{
				{kind: "counter", name: "test", strVal: "1"},
				{kind: "counter", name: "test", strVal: "3"},
			},
			want: map[model.Metric]any{
				{Kind: "counter", Name: "test"}: int64(4),
			},
			wantErr: false,
		},
		{
			name: "Double gauge update",
			metrics: []metric{
				{kind: "gauge", name: "test", strVal: "1.123456"},
				{kind: "gauge", name: "test", strVal: "1.654321"},
			},
			want: map[model.Metric]any{
				{Kind: "gauge", Name: "test"}: float64(1.654321),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.NewStorage()
			for _, metric := range tt.metrics {
				err := UpdateMetric(store, metric.kind, metric.name, metric.strVal)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
			for wantMetric, wantValue := range tt.want {
				value, _ := store.Get(wantMetric)
				assert.Equal(t, wantValue, value)
			}
		})
	}
}
