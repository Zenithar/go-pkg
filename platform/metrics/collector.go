package metrics

// Collector defines metrics collector contract
type Collector interface {
	Collect(ctx context.Context) error
}