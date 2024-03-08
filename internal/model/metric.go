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

func AsGaugeMetric(metric Metric, val any) (float64, bool) {
	if floatVal, ok := val.(float64); ok {
		return floatVal, true
	}
	return 0, false
}

func AsCounterMetric(metric Metric, val any) (int64, bool) {
	if intVal, ok := val.(int64); ok {
		return intVal, true
	}
	return 0, false
}
