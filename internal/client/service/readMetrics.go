package service

import (
	"math/rand"
	"runtime"

	"github.com/artem-benda/monitor/internal/client/storage"
	"github.com/artem-benda/monitor/internal/model"
)

var (
	allocMetric         = model.NewGaugeMetric("Alloc")
	buckHashSysMetric   = model.NewGaugeMetric("BuckHashSys")
	freesMetric         = model.NewGaugeMetric("Frees")
	gccpuFractionMetric = model.NewGaugeMetric("GCCPUFraction")
	gcSysMetric         = model.NewGaugeMetric("GCSys")
	heapAllocMetric     = model.NewGaugeMetric("HeapAlloc")
	heapIdleMetric      = model.NewGaugeMetric("HeapIdle")
	heapInuseMetric     = model.NewGaugeMetric("HeapInuse")
	heapObjectsMetric   = model.NewGaugeMetric("HeapObjects")
	heapReleasedMetric  = model.NewGaugeMetric("HeapReleased")
	heapSysMetric       = model.NewGaugeMetric("HeapSys")
	lastGCMetric        = model.NewGaugeMetric("LastGC")
	lookupsMetric       = model.NewGaugeMetric("Lookups")
	mCacheInuseMetric   = model.NewGaugeMetric("MCacheInuse")
	mCacheSysMetric     = model.NewGaugeMetric("MCacheSys")
	mSpanInuseMetric    = model.NewGaugeMetric("MSpanInuse")
	mSpanSysMetric      = model.NewGaugeMetric("MSpanSys")
	mallocsMetric       = model.NewGaugeMetric("Mallocs")
	nextGCMetric        = model.NewGaugeMetric("NextGC")
	numForcedGCMetric   = model.NewGaugeMetric("NumForcedGC")
	numGCMetric         = model.NewGaugeMetric("NumGC")
	otherSysMetric      = model.NewGaugeMetric("OtherSys")
	pauseTotalNsMetric  = model.NewGaugeMetric("PauseTotalNs")
	stackInuseMetric    = model.NewGaugeMetric("StackInuse")
	stackSysMetric      = model.NewGaugeMetric("StackSys")
	sysMetric           = model.NewGaugeMetric("Sys")
	totalAllocMetric    = model.NewGaugeMetric("TotalAlloc")
	pollCountMetric     = model.NewCounterMetric("PollCount")
	randomValueMetric   = model.NewGaugeMetric("RandomValue")
)

func ReadMetrics(counter storage.Counter) map[model.Metric]any {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	m := make(map[model.Metric]any)

	m[allocMetric] = float64(stats.Alloc)
	m[buckHashSysMetric] = float64(stats.BuckHashSys)
	m[freesMetric] = float64(stats.Frees)
	m[gccpuFractionMetric] = stats.GCCPUFraction
	m[gcSysMetric] = float64(stats.GCSys)
	m[heapAllocMetric] = float64(stats.HeapAlloc)
	m[heapIdleMetric] = float64(stats.HeapIdle)
	m[heapInuseMetric] = float64(stats.HeapInuse)
	m[heapObjectsMetric] = float64(stats.HeapObjects)
	m[heapReleasedMetric] = float64(stats.HeapReleased)
	m[heapSysMetric] = float64(stats.HeapSys)
	m[lastGCMetric] = float64(stats.LastGC)
	m[lookupsMetric] = float64(stats.Lookups)
	m[mCacheInuseMetric] = float64(stats.MCacheInuse)
	m[mCacheSysMetric] = float64(stats.MCacheSys)
	m[mSpanInuseMetric] = float64(stats.MSpanInuse)
	m[mSpanSysMetric] = float64(stats.MSpanSys)
	m[mallocsMetric] = float64(stats.Mallocs)
	m[nextGCMetric] = float64(stats.NextGC)
	m[numForcedGCMetric] = float64(stats.NumForcedGC)
	m[numGCMetric] = float64(stats.NumGC)
	m[otherSysMetric] = float64(stats.OtherSys)
	m[pauseTotalNsMetric] = float64(stats.PauseTotalNs)
	m[stackInuseMetric] = float64(stats.StackInuse)
	m[stackSysMetric] = float64(stats.StackSys)
	m[sysMetric] = float64(stats.Sys)
	m[totalAllocMetric] = float64(stats.TotalAlloc)

	m[pollCountMetric] = int64(counter.IncrementAndGet())
	m[randomValueMetric] = float64(rand.Int63())

	return m
}
