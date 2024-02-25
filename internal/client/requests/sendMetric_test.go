package requests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artem-benda/monitor/internal/dto"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func Test_sendGaugeMetric(t *testing.T) {
	type metric struct {
		name  string
		value float64
	}
	type want struct {
		method      string
		contentType string
	}
	tests := []struct {
		name   string
		metric metric
		want   want
	}{
		{
			name: "Sent counter metric",
			metric: metric{
				name:  "testcounter",
				value: 123.2215654,
			},
			want: want{
				method:      http.MethodPost,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, tt.want.contentType, r.Header.Get("Content-Type"))
				w.Header().Add("Content-type", "application/json")
				w.WriteHeader(http.StatusOK)
			}))
			defer srv.Close()
			client := srv.Client()
			resty := resty.NewWithClient(client)
			resty.SetBaseURL(srv.URL)
			dto := dto.Metrics{ID: tt.metric.name, MType: model.GaugeKind, Value: &tt.metric.value}
			err := sendMetric(resty, dto)
			assert.NoError(t, err)
		})
	}
}

func Test_sendCounterMetric(t *testing.T) {
	type metric struct {
		name  string
		delta int64
	}
	type want struct {
		method      string
		contentType string
	}
	tests := []struct {
		name   string
		metric metric
		want   want
	}{
		{
			name: "Sent counter metric",
			metric: metric{
				name:  "testcounter",
				delta: 123445,
			},
			want: want{
				method:      http.MethodPost,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, tt.want.contentType, r.Header.Get("Content-Type"))
				w.Header().Add("Content-type", "text/plain")
				w.WriteHeader(http.StatusOK)
			}))
			defer srv.Close()
			client := srv.Client()
			resty := resty.NewWithClient(client)
			resty.SetBaseURL(srv.URL)
			dto := dto.Metrics{ID: tt.metric.name, MType: model.CounterKind, Delta: &tt.metric.delta}
			err := sendMetric(resty, dto)
			assert.NoError(t, err)
		})
	}
}
