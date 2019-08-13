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
	tagApplicationName, _     = tag.New("application.name")
	tagApplicationID, _       = tag.New("application.id")
	tagApplicationVersion, _  = tag.New("application.version")
	tagApplicationRevision, _ = tag.New("application.revision")
)

// -----------------------------------------------------------------------------
// Measures
// -----------------------------------------------------------------------------

var (
	mApplicationInfo = stats.Int64("application/info", "Informations about running application", stats.UnitDimensionless)
)

// -----------------------------------------------------------------------------
// Views
// -----------------------------------------------------------------------------

var (
	// InfoView display application information
	InfoView = &view.View{
		Name: mApplicationInfo.Name(), Description: mApplicationInfo.Description(),
		Measure:     mApplicationInfo,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{tagApplicationName, tagApplicationID, tagApplicationVersion, tagApplicationRevision},
	}
)

// -----------------------------------------------------------------------------
type applicationCollector struct {
	name     string
	id       string
	version  string
	revision string
}

// New returns an application metric collector instance
func New(cfg Config) (metrics.Collector, error) {

	// Validate config
	if err := cfg.Validate(); err != nil {
		return xerrors.Errorf("metrics: unable to validate application collector configuration : %w", err)
	}

	// Return collector instance
	return &applicationCollector{
		name:     cfg.Name,
		id:       cfg.ID,
		version:  cfg.Version,
		revision: cfg.Revision,
	}, nil
}

// -----------------------------------------------------------------------------

func (c *applicationCollector) Collect(ctx context.Context) error {
	// Update tag values
	ctx, err := tag.New(ctx,
		tag.Upsert(tagApplicationName, c.name),
		tag.Upsert(tagApplicationID, c.id),
		tag.Upsert(tagApplicationVersion, c.version),
		tag.Upsert(tagApplicationRevision, c.revision),
	)
	if err != nil {
		return xerrors.Errorf("metrics: unable to assign tag values : %w", err)
	}

	// Update measures
	stats.Record(ctx,
		mApplicationInfo.M(1),
	)

	// No error
	return nil
}
