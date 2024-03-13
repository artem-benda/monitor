package requests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestSendAllMetrics(t *testing.T) {
	metrics := map[model.MetricKey]model.MetricValue{
		{Kind: "gauge", Name: "test1"}:   {Gauge: float64(453.223)},
		{Kind: "gauge", Name: "test2"}:   {Gauge: float64(1.554)},
		{Kind: "counter", Name: "test3"}: {Counter: int64(3)},
		{Kind: "counter", Name: "test4"}: {Counter: int64(5)},
	}
	var count int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	client := srv.Client()
	resty := resty.NewWithClient(client)
	resty.SetBaseURL(srv.URL)
	SendAllMetrics(resty, metrics)
	assert.Equal(t, len(metrics), count)
}
