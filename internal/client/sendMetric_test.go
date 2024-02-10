package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sendMetric(t *testing.T) {
	type metric struct {
		kind   string
		name   string
		strVal string
	}
	type want struct {
		method      string
		urlPath     string
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
				kind:   "counter",
				name:   "testcounter",
				strVal: "111",
			},
			want: want{
				method:      http.MethodPost,
				urlPath:     "/update/counter/testcounter/111",
				contentType: "text/plain",
			},
		},
		{
			name: "Sent gauge metric",
			metric: metric{
				kind:   "gauge",
				name:   "testgauge",
				strVal: "222",
			},
			want: want{
				method:      http.MethodPost,
				urlPath:     "/update/gauge/testgauge/222",
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, tt.want.urlPath, r.URL.Path)
				assert.Equal(t, r.Header.Get("Content-Type"), tt.want.contentType)
				w.Header().Add("Content-type", "text/plain")
				w.WriteHeader(http.StatusOK)
			}))
			defer srv.Close()
			client := srv.Client()
			err := sendMetric(client, srv.URL, tt.metric.kind, tt.metric.name, tt.metric.strVal)
			assert.NoError(t, err)
		})
	}
}
