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
