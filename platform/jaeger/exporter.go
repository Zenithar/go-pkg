package jaeger

import (
	"github.com/pkg/errors"
	"go.opencensus.io/exporter/jaeger"
	"go.zenithar.org/pkg/log"
)

// NewExporter creates a new, configured Jaeger exporter.
func NewExporter(config Config) (*jaeger.Exporter, error) {
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: config.CollectorEndpoint,
		AgentEndpoint:     config.AgentEndpoint,
		Username:          config.Username,
		Password:          config.Password,
		OnError: func(err error) {
			log.CheckErr("Error occured in Jaeger exporter", err)
		},
		Process: jaeger.Process{
			ServiceName: config.ServiceName,
		},
	})

	return exporter, errors.Wrap(err, "failed to create jaeger exporter")
}
