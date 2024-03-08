package requests

import (
	"log"

	"github.com/artem-benda/monitor/internal/dto"
	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/go-resty/resty/v2"
)

func SendAllMetrics(resty *resty.Client, metrics map[model.Metric]any) {
	for metric, rawValue := range metrics {
		dto := dto.Metrics{ID: metric.Name, MType: metric.Kind}
		switch metric.Kind {
		case model.GaugeKind:
			{
				if value, ok := model.AsGaugeMetric(metric, rawValue); ok {
					dto.Value = &value
				} else {
					logger.Log.Error("Bad type of gauge metric value")
				}
			}
		case model.CounterKind:
			{
				if value, ok := model.AsCounterMetric(metric, rawValue); ok {
					dto.Delta = &value
				} else {
					logger.Log.Error("Bad type of counter metric value")
				}
			}
		}
		err := sendMetric(resty, dto)
		if err != nil {
			log.Printf("error sending metric %s/%s/%s, %s", metric.Kind, metric.Name, rawValue, err)
		}
	}

}
