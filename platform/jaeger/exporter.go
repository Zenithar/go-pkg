package jaeger

import (
	"context"

	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
	"go.zenithar.org/pkg/log"
	"golang.org/x/xerrors"
)

// newExporter creates a new, configured Jaeger exporter.
func newExporter(config Config) (*jaeger.Exporter, error) {
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

	return exporter, err
}

// RegisterExporter add jaeger as trace exporter
func RegisterExporter(ctx context.Context, conf Config) (func() error, error) {
	// Start tracing

	exporter, err := newExporter(conf)
	if err != nil {
		return nil, xerrors.Errorf("platform: failed to create jaeger exporter: %w", err)
	}

	trace.RegisterExporter(exporter)

	// No error
	return func() error {
		exporter.Flush()
		return nil
	}, nil
}
