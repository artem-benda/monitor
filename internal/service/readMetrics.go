package service

import (
	"math/rand"
	"runtime"
	"strconv"

	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/storage"
)

func ReadMetrics(counter storage.Counter) map[model.Metric]string {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	m := make(map[model.Metric]string)

	m[model.NewGaugeMetric("Alloc")] = strconv.FormatUint(stats.Alloc, 10)
	m[model.NewGaugeMetric("BuckHashSys")] = strconv.FormatUint(stats.BuckHashSys, 10)
	m[model.NewGaugeMetric("Frees")] = strconv.FormatUint(stats.Frees, 10)
	m[model.NewGaugeMetric("GCCPUFraction")] = strconv.FormatFloat(stats.GCCPUFraction, 'f', 5, 64)
	m[model.NewGaugeMetric("GCSys")] = strconv.FormatUint(stats.GCSys, 10)
	m[model.NewGaugeMetric("HeapAlloc")] = strconv.FormatUint(stats.HeapAlloc, 10)
	m[model.NewGaugeMetric("HeapIdle")] = strconv.FormatUint(stats.HeapIdle, 10)
	m[model.NewGaugeMetric("HeapInuse")] = strconv.FormatUint(stats.HeapInuse, 10)
	m[model.NewGaugeMetric("HeapObjects")] = strconv.FormatUint(stats.HeapObjects, 10)
	m[model.NewGaugeMetric("HeapReleased")] = strconv.FormatUint(stats.HeapReleased, 10)
	m[model.NewGaugeMetric("HeapSys")] = strconv.FormatUint(stats.HeapSys, 10)
	m[model.NewGaugeMetric("LastGC")] = strconv.FormatUint(stats.LastGC, 10)
	m[model.NewGaugeMetric("Lookups")] = strconv.FormatUint(stats.Lookups, 10)
	m[model.NewGaugeMetric("MCacheInuse")] = strconv.FormatUint(stats.MCacheInuse, 10)
	m[model.NewGaugeMetric("MCacheSys")] = strconv.FormatUint(stats.MCacheSys, 10)
	m[model.NewGaugeMetric("MSpanInuse")] = strconv.FormatUint(stats.MSpanInuse, 10)
	m[model.NewGaugeMetric("MSpanSys")] = strconv.FormatUint(stats.MSpanSys, 10)
	m[model.NewGaugeMetric("Mallocs")] = strconv.FormatUint(stats.Mallocs, 10)
	m[model.NewGaugeMetric("NextGC")] = strconv.FormatUint(stats.NextGC, 10)
	m[model.NewGaugeMetric("NumForcedGC")] = strconv.FormatUint(uint64(stats.NumForcedGC), 10)
	m[model.NewGaugeMetric("NumGC")] = strconv.FormatUint(uint64(stats.NumGC), 10)
	m[model.NewGaugeMetric("OtherSys")] = strconv.FormatUint(stats.OtherSys, 10)
	m[model.NewGaugeMetric("PauseTotalNs")] = strconv.FormatUint(stats.PauseTotalNs, 10)
	m[model.NewGaugeMetric("StackInuse")] = strconv.FormatUint(stats.StackInuse, 10)
	m[model.NewGaugeMetric("StackSys")] = strconv.FormatUint(stats.StackSys, 10)
	m[model.NewGaugeMetric("Sys")] = strconv.FormatUint(stats.Sys, 10)
	m[model.NewGaugeMetric("TotalAlloc")] = strconv.FormatUint(stats.TotalAlloc, 10)

	m[model.NewCounterMetric("PollCount")] = strconv.FormatUint(counter.IncrementAndGet(), 10)
	m[model.NewGaugeMetric("RandomValue")] = strconv.FormatUint(rand.Uint64(), 10)

	return m
}
