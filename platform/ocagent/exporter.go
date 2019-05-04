package ocagent

import (
	"context"
	"crypto/tls"

	"contrib.go.opencensus.io/exporter/ocagent"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.zenithar.org/pkg/tlsconfig"
	"golang.org/x/xerrors"
	"google.golang.org/grpc/credentials"
)

// newExporter creates a new, configured Prometheus exporter.
func newExporter(config Config) (*ocagent.Exporter, error) {

	sopts := []ocagent.ExporterOption{
		ocagent.WithServiceName(config.ServiceName),
		ocagent.WithAddress(config.Address),
	}

	// Enable TLS if requested
	if config.UseTLS {
		// Client authentication enabled but not required
		clientAuth := tls.VerifyClientCertIfGiven
		if config.TLS.ClientAuthenticationRequired {
			clientAuth = tls.RequireAndVerifyClientCert
		}

		// Generate TLS configuration
		tlsConfig, err := tlsconfig.Server(tlsconfig.Options{
			KeyFile:    config.TLS.PrivateKeyPath,
			CertFile:   config.TLS.CertificatePath,
			CAFile:     config.TLS.CACertificatePath,
			ClientAuth: clientAuth,
		})
		if err != nil {
			return nil, xerrors.Errorf("platform: unbale to initialize TLC settings : %w", err)
		}

		// Create the TLS credentials
		sopts = append(sopts, ocagent.WithTLSCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		sopts = append(sopts, ocagent.WithInsecure())
	}

	// Initialize exporter
	exporter, err := ocagent.NewExporter(sopts...)
	if err != nil {
		return nil, xerrors.Errorf("platform: unable to initialize ocagent exporter: %w", err)
	}

	return exporter, err
}

// RegisterExporter add jaeger as trace exporter
func RegisterExporter(ctx context.Context, conf Config) (func() error, error) {
	// Start tracing
	exporter, err := newExporter(conf)
	if err != nil {
		return nil, xerrors.Errorf("platform: failed to create ocagent exporter: %w", err)
	}

	// Register this exporter
	trace.RegisterExporter(exporter)
	view.RegisterExporter(exporter)

	// No error
	return func() error {
		err := exporter.Stop()
		if err != nil {
			return xerrors.Errorf("platform: unable to stop ocagent exporter : %w", err)
		}
		return nil
	}, nil
}
