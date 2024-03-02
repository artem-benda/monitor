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

type Metric struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
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

type SaveableMetricValue struct {
	Kind         string  `json:"kind"`
	Name         string  `json:"name"`
	Int64Value   int64   `json:"intVal,omitempty"`
	Float64Value float64 `json:"floatVal,omitempty"`
}

func AsSaveableMetric(metric Metric, val any) (SaveableMetricValue, error) {
	if intVal, ok := val.(int64); ok {
		return SaveableMetricValue{metric.Kind, metric.Name, intVal, 0}, nil
	}
	if floatVal, ok := val.(float64); ok {
		return SaveableMetricValue{metric.Kind, metric.Name, 0, floatVal}, nil
	}
	return SaveableMetricValue{}, ErrInvalidMetricValue
}
