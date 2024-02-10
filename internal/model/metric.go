package model

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
