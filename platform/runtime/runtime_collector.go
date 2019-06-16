package runtime

import (
	"runtime"
)

// Stats represents collected statistics
type Stats struct {
	// Go
	Goarch  string `json:"runtime.arch"`
	Goos    string `json:"runtime.os"`
	Version string `json:"runtime.go_version"`
	// CPU
	NumCPU       int64 `json:"runtime.cpu"`
	NumGoRoutine int64 `json:"runtime.goroutines"`
	NumCGOCall   int64 `json:"runtime.cgo_calls"`
	// Memory
	MemAlloc      int64 `json:"mem.alloc"`
	MemTotalAlloc int64 `json:"mem.total"`
	MemSys        int64 `json:"mem.sys"`
	MemLookups    int64 `json:"mem.lookups"`
	MemMallocs    int64 `json:"mem.malloc"`
	MemFrees      int64 `json:"mem.free"`
	// Heap
	HeapAlloc    int64 `json:"mem.heap.alloc"`
	HeapSys      int64 `json:"mem.heap.sys"`
	HeapIdle     int64 `json:"mem.heap.idle"`
	HeapInuse    int64 `json:"mem.heap.inuse"`
	HeapReleased int64 `json:"mem.heap.released"`
	HeapObjects  int64 `json:"mem.heap.objects"`
	// Stack
	StackInuse  int64 `json:"mem.stack.inuse"`
	StackSys    int64 `json:"mem.stack.sys"`
	MSpanInuse  int64 `json:"mem.stack.mspan_inuse"`
	MSpanSys    int64 `json:"mem.stack.mspan_sys"`
	MCacheInuse int64 `json:"mem.stack.mcache_inuse"`
	MCacheSys   int64 `json:"mem.stack.mcache_sys"`
	OtherSys    int64 `json:"mem.othersys"`
	// GC
	GCSys         int64   `json:"mem.gc.sys"`
	NextGC        int64   `json:"mem.gc.next"`
	LastGC        int64   `json:"mem.gc.last"`
	PauseTotalNs  int64   `json:"mem.gc.pause_total"`
	PauseNs       int64   `json:"mem.gc.pause"`
	NumGC         int64   `json:"mem.gc.count"`
	GCCPUFraction float64 `json:"mem.gc.cpu_fraction"`
}

// -----------------------------------------------------------------------------

// Collect is used to fill the struct with current metrics
func (s *Stats) Collect() {
	// Runtime
	s.NumCPU = int64(runtime.NumCPU())
	s.NumGoRoutine = int64(runtime.NumGoroutine())
	s.NumCGOCall = runtime.NumCgoCall()

	// Memory
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)

	s.MemAlloc = int64(m.Alloc)
	s.MemTotalAlloc = int64(m.TotalAlloc)
	s.MemSys = int64(m.Sys)
	s.MemLookups = int64(m.Lookups)
	s.MemMallocs = int64(m.Mallocs)
	s.MemFrees = int64(m.Frees)

	// Heap
	s.HeapAlloc = int64(m.HeapAlloc)
	s.HeapSys = int64(m.HeapSys)
	s.HeapIdle = int64(m.HeapIdle)
	s.HeapInuse = int64(m.HeapInuse)
	s.HeapReleased = int64(m.HeapReleased)
	s.HeapObjects = int64(m.HeapObjects)

	// Stack
	s.StackInuse = int64(m.StackInuse)
	s.StackSys = int64(m.StackSys)
	s.MSpanInuse = int64(m.MSpanInuse)
	s.MSpanSys = int64(m.MSpanSys)
	s.MCacheInuse = int64(m.MCacheInuse)
	s.MCacheSys = int64(m.MCacheSys)
	s.OtherSys = int64(m.OtherSys)

	// GC
	s.GCSys = int64(m.GCSys)
	s.NextGC = int64(m.NextGC)
	s.LastGC = int64(m.LastGC)
	s.PauseTotalNs = int64(m.PauseTotalNs)
	s.PauseNs = int64(m.PauseNs[(m.NumGC+255)%256])
	s.NumGC = int64(m.NumGC)
	s.GCCPUFraction = m.GCCPUFraction

	// Go
	s.Goos = runtime.GOOS
	s.Goarch = runtime.GOARCH
	s.Version = runtime.Version()
}
