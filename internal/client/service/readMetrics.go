package service

import (
	"math/rand"
	"runtime"

	"github.com/artem-benda/monitor/internal/client/storage"
	"github.com/artem-benda/monitor/internal/model"
)

var (
	allocMetric         = model.NewGaugeMetricKey("Alloc")
	buckHashSysMetric   = model.NewGaugeMetricKey("BuckHashSys")
	freesMetric         = model.NewGaugeMetricKey("Frees")
	gccpuFractionMetric = model.NewGaugeMetricKey("GCCPUFraction")
	gcSysMetric         = model.NewGaugeMetricKey("GCSys")
	heapAllocMetric     = model.NewGaugeMetricKey("HeapAlloc")
	heapIdleMetric      = model.NewGaugeMetricKey("HeapIdle")
	heapInuseMetric     = model.NewGaugeMetricKey("HeapInuse")
	heapObjectsMetric   = model.NewGaugeMetricKey("HeapObjects")
	heapReleasedMetric  = model.NewGaugeMetricKey("HeapReleased")
	heapSysMetric       = model.NewGaugeMetricKey("HeapSys")
	lastGCMetric        = model.NewGaugeMetricKey("LastGC")
	lookupsMetric       = model.NewGaugeMetricKey("Lookups")
	mCacheInuseMetric   = model.NewGaugeMetricKey("MCacheInuse")
	mCacheSysMetric     = model.NewGaugeMetricKey("MCacheSys")
	mSpanInuseMetric    = model.NewGaugeMetricKey("MSpanInuse")
	mSpanSysMetric      = model.NewGaugeMetricKey("MSpanSys")
	mallocsMetric       = model.NewGaugeMetricKey("Mallocs")
	nextGCMetric        = model.NewGaugeMetricKey("NextGC")
	numForcedGCMetric   = model.NewGaugeMetricKey("NumForcedGC")
	numGCMetric         = model.NewGaugeMetricKey("NumGC")
	otherSysMetric      = model.NewGaugeMetricKey("OtherSys")
	pauseTotalNsMetric  = model.NewGaugeMetricKey("PauseTotalNs")
	stackInuseMetric    = model.NewGaugeMetricKey("StackInuse")
	stackSysMetric      = model.NewGaugeMetricKey("StackSys")
	sysMetric           = model.NewGaugeMetricKey("Sys")
	totalAllocMetric    = model.NewGaugeMetricKey("TotalAlloc")
	pollCountMetric     = model.NewCounterMetricKey("PollCount")
	randomValueMetric   = model.NewGaugeMetricKey("RandomValue")
)

// ReadMetrics - Получить актуальные значения метрик в виде мапы
func ReadMetrics(counter storage.Counter) map[model.MetricKey]model.MetricValue {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	m := make(map[model.MetricKey]model.MetricValue)

	m[allocMetric] = model.MetricValue{Gauge: float64(stats.Alloc)}
	m[buckHashSysMetric] = model.MetricValue{Gauge: float64(stats.BuckHashSys)}
	m[freesMetric] = model.MetricValue{Gauge: float64(stats.Frees)}
	m[gccpuFractionMetric] = model.MetricValue{Gauge: stats.GCCPUFraction}
	m[gcSysMetric] = model.MetricValue{Gauge: float64(stats.GCSys)}
	m[heapAllocMetric] = model.MetricValue{Gauge: float64(stats.HeapAlloc)}
	m[heapIdleMetric] = model.MetricValue{Gauge: float64(stats.HeapIdle)}
	m[heapInuseMetric] = model.MetricValue{Gauge: float64(stats.HeapInuse)}
	m[heapObjectsMetric] = model.MetricValue{Gauge: float64(stats.HeapObjects)}
	m[heapReleasedMetric] = model.MetricValue{Gauge: float64(stats.HeapReleased)}
	m[heapSysMetric] = model.MetricValue{Gauge: float64(stats.HeapSys)}
	m[lastGCMetric] = model.MetricValue{Gauge: float64(stats.LastGC)}
	m[lookupsMetric] = model.MetricValue{Gauge: float64(stats.Lookups)}
	m[mCacheInuseMetric] = model.MetricValue{Gauge: float64(stats.MCacheInuse)}
	m[mCacheSysMetric] = model.MetricValue{Gauge: float64(stats.MCacheSys)}
	m[mSpanInuseMetric] = model.MetricValue{Gauge: float64(stats.MSpanInuse)}
	m[mSpanSysMetric] = model.MetricValue{Gauge: float64(stats.MSpanSys)}
	m[mallocsMetric] = model.MetricValue{Gauge: float64(stats.Mallocs)}
	m[nextGCMetric] = model.MetricValue{Gauge: float64(stats.NextGC)}
	m[numForcedGCMetric] = model.MetricValue{Gauge: float64(stats.NumForcedGC)}
	m[numGCMetric] = model.MetricValue{Gauge: float64(stats.NumGC)}
	m[otherSysMetric] = model.MetricValue{Gauge: float64(stats.OtherSys)}
	m[pauseTotalNsMetric] = model.MetricValue{Gauge: float64(stats.PauseTotalNs)}
	m[stackInuseMetric] = model.MetricValue{Gauge: float64(stats.StackInuse)}
	m[stackSysMetric] = model.MetricValue{Gauge: float64(stats.StackSys)}
	m[sysMetric] = model.MetricValue{Gauge: float64(stats.Sys)}
	m[totalAllocMetric] = model.MetricValue{Gauge: float64(stats.TotalAlloc)}

	m[pollCountMetric] = model.MetricValue{Counter: int64(counter.IncrementAndGet())}
	m[randomValueMetric] = model.MetricValue{Gauge: float64(rand.Int63())}

	return m
}
