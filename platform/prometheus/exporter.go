package prometheus

import (
	"context"
	"net/http"

	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
	"go.zenithar.org/pkg/log"
	"golang.org/x/xerrors"
)

// newExporter creates a new, configured Prometheus exporter.
func newExporter(config Config) (*prometheus.Exporter, error) {
	exporter, err := prometheus.NewExporter(prometheus.Options{
		Namespace: config.Namespace,
		OnError: func(err error) {
			log.CheckErr("Error occured in Prometheus exporter", err)
		},
	})

	return exporter, err
}

// RegisterExporter adds prometheus exporter
func RegisterExporter(ctx context.Context, conf Config, r *http.ServeMux) (func() error, error) {
	// Start prometheus

	exporter, err := newExporter(conf)
	if err != nil {
		return nil, xerrors.Errorf("platform: unable to register prometheus exporter: %w", err)
	}

	// Add exporter
	view.RegisterExporter(exporter)

	// Add metrics handler
	r.Handle("/metrics", exporter)

	// No error
	return nil, nil
}
