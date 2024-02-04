package client

import (
	"log"
	"net/http"

	"github.com/artem-benda/monitor/internal/model"
)

func SendAllMetrics(httpClient *http.Client, apiUrl string, metrics map[model.Metric]string) {
	for metric, strValue := range metrics {
		err := SendMetric(httpClient, apiUrl, metric.Kind, metric.Name, strValue)
		if err != nil {
			log.Printf("error sending metric %s/%s/%s", metric.Kind, metric.Name, strValue)
		}
	}
}
