package prometheus

// Config holds information for configuring the Prometheus exporter
type Config struct {
	Namespace string `toml:"namespace" comment:"Prometheus namespace"`
}
