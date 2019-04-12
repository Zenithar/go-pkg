package prometheus

import (
	"github.com/pkg/errors"
	"go.opencensus.io/exporter/prometheus"
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
