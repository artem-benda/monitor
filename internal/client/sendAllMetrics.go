package client

import (
	"log"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/go-resty/resty/v2"
)

func SendAllMetrics(resty *resty.Client, apiURL string, metrics map[model.Metric]string) {
	for metric, strValue := range metrics {
		err := sendMetric(resty, apiURL, metric.Kind, metric.Name, strValue)
		if err != nil {
			log.Printf("error sending metric %s/%s/%s, %s", metric.Kind, metric.Name, strValue, err)
		}
	}
}
