package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestSendAllMetrics(t *testing.T) {
	metrics := map[model.Metric]string{
		{Kind: "gauge", Name: "test1"}:   "1.223",
		{Kind: "gauge", Name: "test2"}:   "1.554",
		{Kind: "counter", Name: "test3"}: "3",
		{Kind: "counter", Name: "test4"}: "5",
	}
	var count int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	client := srv.Client()
	SendAllMetrics(client, srv.URL, metrics)
	assert.Equal(t, len(metrics), count)
}
