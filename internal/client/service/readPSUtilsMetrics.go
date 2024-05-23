package service

import (
	"github.com/artem-benda/monitor/internal/model"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

var (
	totalMemoryMetric     = model.NewGaugeMetricKey("TotalMemory")
	freeMemoryMetric      = model.NewGaugeMetricKey("FreeMemory")
	cpuUtilization1Metric = model.NewGaugeMetricKey("CPUutilization1")
)

func ReadPSUtilsMetrics() map[model.MetricKey]model.MetricValue {
	m := make(map[model.MetricKey]model.MetricValue)

	v, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(0, false)

	m[totalMemoryMetric] = model.MetricValue{Gauge: float64(v.Total)}
	m[freeMemoryMetric] = model.MetricValue{Gauge: float64(v.Free)}
	m[cpuUtilization1Metric] = model.MetricValue{Gauge: float64(c[0])}

	return m
}
