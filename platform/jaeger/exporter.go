package jaeger

import (
	"context"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
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

// RegisterExporter add jaeger as trace exporter
func RegisterExporter(ctx context.Context, debug bool, conf Config) error {
	// Start tracing

	exporter, err := NewExporter(conf)
	if err != nil {
		return errors.Wrap(err, "unable to register jaeger exporter")
	}

	trace.RegisterExporter(exporter)

	// Trace everything when debugging is enabled
	if debug {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	// No error
	return nil
}
