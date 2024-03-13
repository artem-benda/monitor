package model

import (
	"errors"
	"strconv"
)

const (
	GaugeKind   = "gauge"
	CounterKind = "counter"
)

var ErrInvalidMetricValue = errors.New("InvalidMetricValueError")

type MetricKey struct {
	Kind string
	Name string
}

type MetricValue struct {
	Gauge   float64
	Counter int64
}

func ValidMetricKind(s string) bool {
	return s == GaugeKind || s == CounterKind
}

func NewGaugeMetricKey(name string) MetricKey {
	return MetricKey{Kind: GaugeKind, Name: name}
}

func NewCounterMetricKey(name string) MetricKey {
	return MetricKey{Kind: CounterKind, Name: name}
}

func StringValue(metricKey MetricKey, val MetricValue) (string, error) {
	switch metricKey.Kind {
	case CounterKind:
		{
			strVal := strconv.FormatInt(val.Counter, 10)
			return strVal, nil
		}
	case GaugeKind:
		{
			strVal := strconv.FormatFloat(val.Gauge, 'f', -1, 64)
			return strVal, nil
		}
	}
	return "", ErrInvalidMetricValue
}

type SaveableMetricValue struct {
	Kind    string  `json:"kind"`
	Name    string  `json:"name"`
	Counter int64   `json:"counter,omitempty"`
	Gauge   float64 `json:"gauge,omitempty"`
}

func AsSaveableMetric(key MetricKey, val MetricValue) SaveableMetricValue {
	return SaveableMetricValue{Kind: key.Kind, Name: key.Name, Gauge: val.Gauge, Counter: val.Counter}
}
