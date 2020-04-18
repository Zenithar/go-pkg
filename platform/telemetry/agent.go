package telemetry

import (
	"fmt"
	"os"
	"time"

	"contrib.go.opencensus.io/exporter/ocagent"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// Agent initializes an opentelemetry agent.
func Agent(serviceName string) (func(), error) {
	// Retrieve agent url from environment.
	ocAgentAddr, ok := os.LookupEnv("OTEL_AGENT_ENDPOINT")
	if !ok {
		ocAgentAddr = fmt.Sprintf("%s:%d", ocagent.DefaultAgentHost, ocagent.DefaultAgentPort)
	}

	// Create opentelemetry gRPC connection
	oce, err := ocagent.NewExporter(
		ocagent.WithAddress(ocAgentAddr),
		ocagent.WithInsecure(),
		ocagent.WithServiceName(serviceName),
		ocagent.WithReconnectionPeriod(20*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("telemetry: unable to connect agent (%s): %w", ocAgentAddr, err)
	}

	// Register exporter
	trace.RegisterExporter(oce)
	view.RegisterExporter(oce)

	// No error
	return oce.Flush, nil
}
