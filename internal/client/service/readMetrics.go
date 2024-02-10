package service

import (
	"math/rand"
	"runtime"
	"strconv"

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

func ReadMetrics(counter storage.Counter) map[model.Metric]string {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	m := make(map[model.Metric]string)

	m[allocMetric] = strconv.FormatUint(stats.Alloc, 10)
	m[buckHashSysMetric] = strconv.FormatUint(stats.BuckHashSys, 10)
	m[freesMetric] = strconv.FormatUint(stats.Frees, 10)
	m[gccpuFractionMetric] = strconv.FormatFloat(stats.GCCPUFraction, 'f', -1, 64)
	m[gcSysMetric] = strconv.FormatUint(stats.GCSys, 10)
	m[heapAllocMetric] = strconv.FormatUint(stats.HeapAlloc, 10)
	m[heapIdleMetric] = strconv.FormatUint(stats.HeapIdle, 10)
	m[heapInuseMetric] = strconv.FormatUint(stats.HeapInuse, 10)
	m[heapObjectsMetric] = strconv.FormatUint(stats.HeapObjects, 10)
	m[heapReleasedMetric] = strconv.FormatUint(stats.HeapReleased, 10)
	m[heapSysMetric] = strconv.FormatUint(stats.HeapSys, 10)
	m[lastGCMetric] = strconv.FormatUint(stats.LastGC, 10)
	m[lookupsMetric] = strconv.FormatUint(stats.Lookups, 10)
	m[mCacheInuseMetric] = strconv.FormatUint(stats.MCacheInuse, 10)
	m[mCacheSysMetric] = strconv.FormatUint(stats.MCacheSys, 10)
	m[mSpanInuseMetric] = strconv.FormatUint(stats.MSpanInuse, 10)
	m[mSpanSysMetric] = strconv.FormatUint(stats.MSpanSys, 10)
	m[mallocsMetric] = strconv.FormatUint(stats.Mallocs, 10)
	m[nextGCMetric] = strconv.FormatUint(stats.NextGC, 10)
	m[numForcedGCMetric] = strconv.FormatUint(uint64(stats.NumForcedGC), 10)
	m[numGCMetric] = strconv.FormatUint(uint64(stats.NumGC), 10)
	m[otherSysMetric] = strconv.FormatUint(stats.OtherSys, 10)
	m[pauseTotalNsMetric] = strconv.FormatUint(stats.PauseTotalNs, 10)
	m[stackInuseMetric] = strconv.FormatUint(stats.StackInuse, 10)
	m[stackSysMetric] = strconv.FormatUint(stats.StackSys, 10)
	m[sysMetric] = strconv.FormatUint(stats.Sys, 10)
	m[totalAllocMetric] = strconv.FormatUint(stats.TotalAlloc, 10)

	m[pollCountMetric] = strconv.FormatUint(counter.IncrementAndGet(), 10)
	m[randomValueMetric] = strconv.FormatUint(rand.Uint64(), 10)

	return m
}
