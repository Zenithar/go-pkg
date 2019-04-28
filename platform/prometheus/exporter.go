package prometheus

import (
	"context"
	"net/http"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/pkg/errors"
	"go.opencensus.io/stats/view"
	"go.zenithar.org/pkg/log"
)

// NewExporter creates a new, configured Prometheus exporter.
func NewExporter(config Config) (*prometheus.Exporter, error) {
	exporter, err := prometheus.NewExporter(prometheus.Options{
		Namespace: config.Namespace,
		OnError: func(err error) {
			log.CheckErr("Error occured in Prometheus exporter", err)
		},
	})

	return exporter, errors.Wrap(err, "failed to create prometheus exporter")
}

// RegisterExporter adds prometheus exporter
func RegisterExporter(ctx context.Context, conf Config, r *http.ServeMux) error {
	// Start prometheus

	exporter, err := NewExporter(conf)
	if err != nil {
		return errors.Wrap(err, "unable to register prometheus handler")
	}

	// Add exporter
	view.RegisterExporter(exporter)

	// Add metrics handler
	r.Handle("/metrics", exporter)

	// No error
	return nil
}
