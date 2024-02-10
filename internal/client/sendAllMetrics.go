package client

import (
	"log"
	"net/http"

	"github.com/artem-benda/monitor/internal/model"
)

func SendAllMetrics(httpClient *http.Client, apiURL string, metrics map[model.Metric]string) {
	for metric, strValue := range metrics {
		err := sendMetric(httpClient, apiURL, metric.Kind, metric.Name, strValue)
		if err != nil {
			log.Printf("error sending metric %s/%s/%s, %s", metric.Kind, metric.Name, strValue, err)
		}
	}
}
