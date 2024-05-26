package model

import (
	"errors"
	"strconv"

	"github.com/artem-benda/monitor/internal/dto"
)

// Типы метрик
const (
	GaugeKind   = "gauge"
	CounterKind = "counter"
)

var ErrInvalidMetricValue = errors.New("InvalidMetricValueError")

// MetricKey - доменная модель ключ метрики
type MetricKey struct {
	Kind string
	Name string
}

// MetricKey - доменная модель значение метрики
type MetricValue struct {
	Gauge   float64
	Counter int64
}

// MetricKey - доменная модель ключ и значение метрики
type MetricKeyWithValue struct {
	Kind    string
	Name    string
	Gauge   float64
	Counter int64
}

// ValidMetricKind - проверить, является ли строка допустимым значением типа метрики
func ValidMetricKind(s string) bool {
	return s == GaugeKind || s == CounterKind
}

// NewGaugeMetricKey создать новый ключ метрики с типом Gauge
func NewGaugeMetricKey(name string) MetricKey {
	return MetricKey{Kind: GaugeKind, Name: name}
}

// NewCounterMetricKey создать новый ключ метрики с типом Counter
func NewCounterMetricKey(name string) MetricKey {
	return MetricKey{Kind: CounterKind, Name: name}
}

// StringValue - преобразовать значение метрики в строку
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

// AsDto - Преобразовать значение метрики в модель API
func AsDto(key MetricKey, val MetricValue) dto.Metrics {
	return dto.Metrics{MType: key.Kind, ID: key.Name, Value: &val.Gauge, Delta: &val.Counter}
}
