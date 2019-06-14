package diagnostic

import (
	"context"
	"runtime"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"golang.org/x/xerrors"
)

var (
	tagKeyProcessOS, _        = tag.NewKey("runtime.os")
	tagKeyProcessArch, _      = tag.NewKey("runtime.arch")
	tagKeyProcessGoVersion, _ = tag.NewKey("runtime.go_version")
)

var (
	metricProcessGoroutine          = stats.Int64("process/goroutine", "Number of running goroutines", stats.UnitDimensionless)
	metricProcessMemoryTotalAlloc   = stats.Int64("process/memory/total_alloc", "Total bytes allocated for heap objects", stats.UnitBytes)
	metricProcessMemorySys          = stats.Int64("process/memory/sys", "Total of memory obtained from the OS", stats.UnitBytes)
	metricProcessMemoryMalloc       = stats.Int64("process/memory/malloc", "Number of heap objects allocated", stats.UnitDimensionless)
	metricProcessMemoryFree         = stats.Int64("process/memory/free", "Number of heap objects freed", stats.UnitDimensionless)
	metricProcessMemoryHeapAlloc    = stats.Int64("process/memory/heap_alloc", "Allocated heap objects", stats.UnitBytes)
	metricProcessMemoryHeapReleased = stats.Int64("process/memory/heap_released", "Physical memory returned to the OS", stats.UnitBytes)
	metricProcessMemoryHeapObject   = stats.Int64("process/memory/heap_object", "Number of allocated heap objects", stats.UnitDimensionless)
)

var (
	processViews = []*view.View{{
		Name: metricProcessGoroutine.Name(), Description: metricProcessGoroutine.Description(),
		Measure:     metricProcessGoroutine,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}, {
		Name: metricProcessMemoryTotalAlloc.Name(), Description: metricProcessMemoryTotalAlloc.Description(),
		Measure:     metricProcessMemoryTotalAlloc,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}, {
		Name: metricProcessMemorySys.Name(), Description: metricProcessMemorySys.Description(),
		Measure:     metricProcessMemorySys,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}, {
		Name: metricProcessMemoryMalloc.Name(), Description: metricProcessMemoryMalloc.Description(),
		Measure:     metricProcessMemoryMalloc,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}, {
		Name: metricProcessMemoryFree.Name(), Description: metricProcessMemoryFree.Description(),
		Measure:     metricProcessMemoryFree,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}, {
		Name: metricProcessMemoryHeapAlloc.Name(), Description: metricProcessMemoryHeapAlloc.Description(),
		Measure:     metricProcessMemoryHeapAlloc,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}, {
		Name: metricProcessMemoryHeapReleased.Name(), Description: metricProcessMemoryHeapReleased.Description(),
		Measure:     metricProcessMemoryHeapReleased,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}, {
		Name: metricProcessMemoryHeapObject.Name(), Description: metricProcessMemoryHeapObject.Description(),
		Measure:     metricProcessMemoryHeapObject,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
	}}
)

// MonitorProcess monitors the resource used by the process (memory, goroutines).
// OS-Related monitoring is out of the scope of this function (cpu, io (network, disk), ...).
func monitorProcess(interval time.Duration) func() {

	// Register process views
	if err := view.Register(processViews...); err != nil {
		panic(xerrors.Errorf("diagnostic: unable to register runtime metrics view: %w", err))
	}

	var stopReport = make(chan struct{}, 1)

	// For monitoring routine
	go func() {
		tick := time.NewTicker(interval)
		defer tick.Stop()

		ctx := context.Background()
		ctx, _ = tag.New(ctx, // nolint: errcheck
			tag.Upsert(tagKeyProcessOS, runtime.GOOS),
			tag.Upsert(tagKeyProcessArch, runtime.GOARCH),
			tag.Upsert(tagKeyProcessGoVersion, runtime.Version()),
		)

		for {
			select {
			case <-stopReport:
				return
			case <-tick.C:
				var memstats runtime.MemStats
				runtime.ReadMemStats(&memstats)

				stats.Record(ctx,
					metricProcessGoroutine.M(int64(runtime.NumGoroutine())),
					metricProcessMemoryTotalAlloc.M(int64(memstats.TotalAlloc)),
					metricProcessMemorySys.M(int64(memstats.Sys)),
					metricProcessMemoryMalloc.M(int64(memstats.Mallocs)),
					metricProcessMemoryFree.M(int64(memstats.Frees)),
					metricProcessMemoryHeapAlloc.M(int64(memstats.HeapAlloc)),
					metricProcessMemoryHeapReleased.M(int64(memstats.HeapReleased)),
					metricProcessMemoryHeapObject.M(int64(memstats.HeapObjects)),
				)
			}
		}
	}()
	return func() {
		close(stopReport)
	}
}