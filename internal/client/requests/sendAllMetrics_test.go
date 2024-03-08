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
	metrics := map[model.Metric]any{
		{Kind: "gauge", Name: "test1"}:   float64(453.223),
		{Kind: "gauge", Name: "test2"}:   float64(1.554),
		{Kind: "counter", Name: "test3"}: int64(3),
		{Kind: "counter", Name: "test4"}: int64(5),
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
