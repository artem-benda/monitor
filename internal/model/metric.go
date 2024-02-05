package model

import "strconv"

const (
	GaugeKind   = "gauge"
	CounterKind = "counter"
)

type Metric struct {
	Kind string
	Name string
}

func ValidMetricKind(s string) bool {
	return s == GaugeKind || s == CounterKind
}

func NewGaugeMetric(name string) Metric {
	return Metric{Kind: GaugeKind, Name: name}
}

func NewCounterMetric(name string) Metric {
	return Metric{Kind: CounterKind, Name: name}
}

func StringValue(metric Metric, val any) (string, bool) {
	switch metric.Kind {
	case CounterKind:
		{
			if intVal, ok := val.(int64); ok {
				strVal := strconv.FormatInt(intVal, 10)
				return strVal, true
			}
		}
	case GaugeKind:
		{
			if floatVal, ok := val.(float64); ok {
				strVal := strconv.FormatFloat(floatVal, 'f', -1, 64)
				return strVal, true
			}
		}
	}
	return "", false
}
