package application

import (
	"context"

	"go.zenithar.org/pkg/platform/metrics"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"golang.org/x/xerrors"
)

// -----------------------------------------------------------------------------
// Tags
// -----------------------------------------------------------------------------

var (
	tagKeyProcessOS, _        = tag.NewKey("runtime.os")
	tagKeyProcessArch, _      = tag.NewKey("runtime.arch")
	tagKeyProcessGoVersion, _ = tag.NewKey("runtime.go_version")
)

// -----------------------------------------------------------------------------
// Measures
// -----------------------------------------------------------------------------

var (
	metricProcessGoroutine          = stats.Int64("go/runtime/goroutine", "Number of running goroutines", stats.UnitDimensionless)
	metricProcessThread 			= stats.Int64("go/runtime/thread", "Number of OS threads created", stats.UnitDimensionless)
	metricProcessGoinfo             = stats.Int64("go/runtime/go_info", "Information about the Go environment", stats.UnitDimensionless)
	metricProcessMemoryTotalAlloc   = stats.Int64("go/runtime/memory/total_alloc", "Total bytes allocated for heap objects", stats.UnitBytes)
	metricProcessMemorySys          = stats.Int64("go/runtime/memory/sys", "Total of memory obtained from the OS", stats.UnitBytes)
	metricProcessMemoryMalloc       = stats.Int64("go/runtime/memory/malloc", "Number of heap objects allocated", stats.UnitDimensionless)
	metricProcessMemoryFree         = stats.Int64("go/runtime/memory/free", "Number of heap objects freed", stats.UnitDimensionless)
	metricProcessMemoryHeapAlloc    = stats.Int64("go/runtime/memory/heap_alloc", "Allocated heap objects", stats.UnitBytes)
	metricProcessMemoryHeapReleased = stats.Int64("go/runtime/memory/heap_released", "Physical memory returned to the OS", stats.UnitBytes)
	metricProcessMemoryHeapObject   = stats.Int64("go/runtime/memory/heap_object", "Number of allocated heap objects", stats.UnitDimensionless)
)

// -----------------------------------------------------------------------------
// Views
// -----------------------------------------------------------------------------

var (
	processViews = []*view.View{{
		Name: metricProcessGoroutine.Name(), Description: metricProcessGoroutine.Description(),
		Measure:     metricProcessGoroutine,
		Aggregation: view.LastValue(),
		TagKeys: []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}, {
		Name: metricProcessGoinfo.Name(), Description: metricProcessGoinfo.Description(),
		Measure:     metricProcessGoinfo,
		Aggregation: view.LastValue(),
		TagKeys: []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}, {
		Name: metricProcessMemoryTotalAlloc.Name(), Description: metricProcessMemoryTotalAlloc.Description(),
		Measure:     metricProcessMemoryTotalAlloc,
		Aggregation: view.LastValue(),
		TagKeys: []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}, {
		Name: metricProcessMemorySys.Name(), Description: metricProcessMemorySys.Description(),
		Measure:     metricProcessMemorySys,
		Aggregation: view.LastValue(),
		TagKeys: []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}, {
		Name: metricProcessMemoryMalloc.Name(), Description: metricProcessMemoryMalloc.Description(),
		Measure:     metricProcessMemoryMalloc,
		Aggregation: view.LastValue(),
		TagKeys: []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}, {
		Name: metricProcessMemoryFree.Name(), Description: metricProcessMemoryFree.Description(),
		Measure:     metricProcessMemoryFree,
		Aggregation: view.LastValue(),
		TagKeys: []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}, {
		Name: metricProcessMemoryHeapAlloc.Name(), Description: metricProcessMemoryHeapAlloc.Description(),
		Measure:     metricProcessMemoryHeapAlloc,
		Aggregation: view.LastValue(),
		TagKeys: []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}, {
		Name: metricProcessMemoryHeapReleased.Name(), Description: metricProcessMemoryHeapReleased.Description(),
		Measure:     metricProcessMemoryHeapReleased,
		Aggregation: view.LastValue(),
		TagKeys: []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}, {
		Name: metricProcessMemoryHeapObject.Name(), Description: metricProcessMemoryHeapObject.Description(),
		Measure:     metricProcessMemoryHeapObject,
		Aggregation: view.LastValue(),
		TagKeys: []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}}
)

// -----------------------------------------------------------------------------
type runtimeCollector struct {
	name     string
}

// New returns an application metric collector instance
func New(name string) (metrics.Collector, error) {
	// Return collector instance
	return &runtimeCollector{
		name:     name,
	}, nil
}

// -----------------------------------------------------------------------------

func (c *runtimeCollector) Collect(ctx context.Context) error {
	// Update tag values
	ctx, err := tag.New(ctx,
		tag.Upsert(tagKeyProcessOS, runtime.GOOS),
		tag.Upsert(tagKeyProcessArch, runtime.GOARCH),
		tag.Upsert(tagKeyProcessGoVersion, runtime.Version()),
	)
	if err != nil {
		return xerrors.Errorf("metrics: unable to assign tag values : %w", err)
	}

	// Get memory information from Go runtime
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)

	// Update measures
	stats.Record(ctx,
		metricProcessGoroutine.M(int64(runtime.NumGoroutine())),
		metricProcessGoinfo.M(1),
		metricProcessMemoryTotalAlloc.M(int64(memstats.TotalAlloc)),
		metricProcessMemorySys.M(int64(memstats.Sys)),
		metricProcessMemoryMalloc.M(int64(memstats.Mallocs)),
		metricProcessMemoryFree.M(int64(memstats.Frees)),
		metricProcessMemoryHeapAlloc.M(int64(memstats.HeapAlloc)),
		metricProcessMemoryHeapReleased.M(int64(memstats.HeapReleased)),
		metricProcessMemoryHeapObject.M(int64(memstats.HeapObjects)),
	)

	// No error
	return nil
}
