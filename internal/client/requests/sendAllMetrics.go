package requests

import (
	"log"

	"github.com/artem-benda/monitor/internal/dto"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/go-resty/resty/v2"
)

func SendAllMetrics(resty *resty.Client, metrics map[model.MetricKey]model.MetricValue) {
	for metric, val := range metrics {
		dto := dto.Metrics{ID: metric.Name, MType: metric.Kind, Value: &val.Gauge, Delta: &val.Counter}
		err := sendMetric(resty, dto)
		if err != nil {
			log.Printf("error sending metric %s/%s [ %f  %d ], %s", metric.Kind, metric.Name, val.Gauge, val.Counter, err)
		}
	}

}
