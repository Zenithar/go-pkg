package runtime

import (
	"context"

	"golang.org/x/xerrors"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	tagKeyApplicationName, _     = tag.NewKey("platform.application.name")
	tagKeyApplicationID, _       = tag.NewKey("platform.application.instance_id")
	tagKeyApplicationVersion, _  = tag.NewKey("platform.application.version")
	tagKeyApplicationRevision, _ = tag.NewKey("platform.application.revision")
	tagKeyProcessOS, _           = tag.NewKey("platform.runtime.os")
	tagKeyProcessArch, _         = tag.NewKey("platform.runtime.arch")
	tagKeyProcessGoVersion, _    = tag.NewKey("platform.runtime.go_version")
)

type ocPublisher struct {
	cfg Config

	mBuildInfo     *stats.Int64Measure
	mNumCPU        *stats.Int64Measure
	mNumGoroutine  *stats.Int64Measure
	mNumCGOCalls   *stats.Int64Measure
	mMemAlloc      *stats.Int64Measure
	mMemTotal      *stats.Int64Measure
	mMemSys        *stats.Int64Measure
	mMemLookups    *stats.Int64Measure
	mMemMalloc     *stats.Int64Measure
	mMemFree       *stats.Int64Measure
	mHeapAlloc     *stats.Int64Measure
	mHeapSys       *stats.Int64Measure
	mHeapIdle      *stats.Int64Measure
	mHeapInUse     *stats.Int64Measure
	mHeapReleased  *stats.Int64Measure
	mHeapObjects   *stats.Int64Measure
	mStackInUse    *stats.Int64Measure
	mStackSys      *stats.Int64Measure
	mMSpanInUse    *stats.Int64Measure
	mMSpanSys      *stats.Int64Measure
	mMCacheInUse   *stats.Int64Measure
	mMCacheSys     *stats.Int64Measure
	mOtherSys      *stats.Int64Measure
	mGCSys         *stats.Int64Measure
	mGCNext        *stats.Int64Measure
	mGCLast        *stats.Int64Measure
	mGCPauseTotal  *stats.Int64Measure
	mGCPause       *stats.Int64Measure
	mGCCount       *stats.Int64Measure
	mGCCPUFraction *stats.Float64Measure
}

// OpenCensus initialize an opencensus metric publisher
func OpenCensus(cfg Config) Publisher {
	defaultTags := []tag.Key{
		tagKeyApplicationName,
		tagKeyApplicationID,
		tagKeyApplicationVersion,
		tagKeyApplicationRevision,
		tagKeyProcessOS,
		tagKeyProcessArch,
		tagKeyProcessGoVersion,
	}

	p := &ocPublisher{
		cfg: cfg,

		mBuildInfo:     int64Measure("runtime/build_info", "Build informations", stats.UnitDimensionless, defaultTags),
		mNumCPU:        int64Measure("runtime/cpu_count", "Num of running CPU", stats.UnitDimensionless, defaultTags),
		mNumGoroutine:  int64Measure("runtime/goroutine_count", "Num of goroutines running", stats.UnitDimensionless, defaultTags),
		mNumCGOCalls:   int64Measure("runtime/cgo_call_count", "Num of CGO calls done", stats.UnitDimensionless, defaultTags),
		mMemAlloc:      int64Measure("runtime/memory/alloc", "Number of heap objects allocated", stats.UnitDimensionless, defaultTags),
		mMemTotal:      int64Measure("runtime/memory/total_alloc", "Total bytes allocated for heap objects", stats.UnitBytes, defaultTags),
		mMemSys:        int64Measure("runtime/memory/sys", "Total of memory obtained from the OS", stats.UnitBytes, defaultTags),
		mMemLookups:    int64Measure("runtime/memory/lookups", "Total of memory obtained from the OS", stats.UnitDimensionless, defaultTags),
		mMemFree:       int64Measure("runtime/memory/free", "Number of heap objects freed", stats.UnitDimensionless, defaultTags),
		mMemMalloc:     int64Measure("runtime/memory/malloc", "Number of mallocs done", stats.UnitDimensionless, defaultTags),
		mHeapAlloc:     int64Measure("runtime/memory/heap_alloc", "Allocated heap objects", stats.UnitBytes, defaultTags),
		mHeapSys:       int64Measure("runtime/memory/heap_sys", "Physical memory returned to the OS", stats.UnitBytes, defaultTags),
		mHeapIdle:      int64Measure("runtime/memory/heap_idle", "Physical memory returned to the OS", stats.UnitBytes, defaultTags),
		mHeapInUse:     int64Measure("runtime/memory/heap_inuse", "Physical memory returned to the OS", stats.UnitBytes, defaultTags),
		mHeapReleased:  int64Measure("runtime/memory/heap_released", "Physical memory returned to the OS", stats.UnitBytes, defaultTags),
		mHeapObjects:   int64Measure("runtime/memory/heap_objects", "Number of allocated heap objects", stats.UnitDimensionless, defaultTags),
		mStackInUse:    int64Measure("runtime/memory/stack_inuse", "", stats.UnitBytes, defaultTags),
		mStackSys:      int64Measure("runtime/memory/stack_sys", "", stats.UnitBytes, defaultTags),
		mMSpanInUse:    int64Measure("runtime/memory/stack_mspan_inuse", "", stats.UnitBytes, defaultTags),
		mMSpanSys:      int64Measure("runtime/memory/stack_mspan_sys", "", stats.UnitBytes, defaultTags),
		mMCacheInUse:   int64Measure("runtime/memory/stack_mcache_inuse", "", stats.UnitBytes, defaultTags),
		mMCacheSys:     int64Measure("runtime/memory/stack_mcache_sys", "", stats.UnitBytes, defaultTags),
		mOtherSys:      int64Measure("runtime/memory/other_sys", "", stats.UnitBytes, defaultTags),
		mGCSys:         int64Measure("runtime/memory/gc_sys", "", stats.UnitDimensionless, defaultTags),
		mGCNext:        int64Measure("runtime/memory/gc_next", "", stats.UnitDimensionless, defaultTags),
		mGCLast:        int64Measure("runtime/memory/gc_last", "", stats.UnitDimensionless, defaultTags),
		mGCPauseTotal:  int64Measure("runtime/memory/gc_pause_total_ns", "", stats.UnitDimensionless, defaultTags),
		mGCPause:       int64Measure("runtime/memory/gc_pause_ns", "", stats.UnitDimensionless, defaultTags),
		mGCCount:       int64Measure("runtime/memory/gc_count", "", stats.UnitDimensionless, defaultTags),
		mGCCPUFraction: float64Measure("runtime/memory/gc_cpu_fraction", "", stats.UnitDimensionless, defaultTags),
	}

	// Return publisher
	return p
}

// -----------------------------------------------------------------------------

func (a *ocPublisher) Publish(rstats *Stats) {
	// Set tag values
	ctx := context.Background()
	ctx, _ = tag.New(ctx, // nolint: errcheck
		tag.Upsert(tagKeyProcessOS, rstats.Goos),
		tag.Upsert(tagKeyProcessArch, rstats.Goarch),
		tag.Upsert(tagKeyProcessGoVersion, rstats.Version),
		tag.Upsert(tagKeyApplicationID, a.cfg.ID),
		tag.Upsert(tagKeyApplicationName, a.cfg.Name),
		tag.Upsert(tagKeyApplicationRevision, a.cfg.Revision),
		tag.Upsert(tagKeyApplicationVersion, a.cfg.Version),
	)

	// Set metric values
	stats.Record(ctx,
		a.mBuildInfo.M(1),
		a.mNumCPU.M(rstats.NumCPU),
		a.mNumGoroutine.M(rstats.NumGoRoutine),
		a.mNumCGOCalls.M(rstats.NumCGOCall),
		a.mMemAlloc.M(rstats.MemAlloc),
		a.mMemTotal.M(rstats.MemTotalAlloc),
		a.mMemSys.M(rstats.MemSys),
		a.mMemLookups.M(rstats.MemLookups),
		a.mMemFree.M(rstats.MemFrees),
		a.mMemMalloc.M(rstats.MemMallocs),
		a.mHeapAlloc.M(rstats.HeapAlloc),
		a.mHeapSys.M(rstats.HeapSys),
		a.mHeapIdle.M(rstats.HeapIdle),
		a.mHeapInUse.M(rstats.HeapInuse),
		a.mHeapReleased.M(rstats.HeapReleased),
		a.mHeapObjects.M(rstats.HeapObjects),
		a.mStackInUse.M(rstats.StackInuse),
		a.mStackSys.M(rstats.StackSys),
		a.mMSpanInUse.M(rstats.MSpanInuse),
		a.mMSpanSys.M(rstats.MSpanSys),
		a.mMCacheInUse.M(rstats.MCacheInuse),
		a.mMCacheSys.M(rstats.MCacheSys),
		a.mOtherSys.M(rstats.OtherSys),
		a.mGCSys.M(rstats.GCSys),
		a.mGCNext.M(rstats.NextGC),
		a.mGCLast.M(rstats.LastGC),
		a.mGCPauseTotal.M(rstats.PauseTotalNs),
		a.mGCPause.M(rstats.PauseNs),
		a.mGCCount.M(rstats.NumGC),
		a.mGCCPUFraction.M(rstats.GCCPUFraction),
	)
}

// -----------------------------------------------------------------------------

func int64Measure(name, description, unit string, tags []tag.Key) *stats.Int64Measure {
	// Initialize a measure
	m := stats.Int64(name, description, unit)

	// Initialize a view
	v := &view.View{
		Name: m.Name(), Description: m.Description(),
		Measure:     m,
		Aggregation: view.LastValue(),
		TagKeys:     tags,
	}

	// Register the view
	if err := view.Register(v); err != nil {
		panic(xerrors.Errorf("opencensus: unable to register runtime metrics view: %w", err))
	}

	// Return the measure
	return m
}

func float64Measure(name, description, unit string, tags []tag.Key) *stats.Float64Measure {
	// Initialize a measure
	m := stats.Float64(name, description, unit)

	// Initialize a view
	v := &view.View{
		Name: m.Name(), Description: m.Description(),
		Measure:     m,
		Aggregation: view.LastValue(),
		TagKeys:     tags,
	}

	// Register the view
	if err := view.Register(v); err != nil {
		panic(xerrors.Errorf("opencensus: unable to register runtime metrics view: %w", err))
	}

	// Return the measure
	return m
}
