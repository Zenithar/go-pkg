package ocagent

import "golang.org/x/xerrors"

// Config holds information necessary for sending trace to OpenCensus agent.
type Config struct {
	// Address of the OC Agent
	Address string `toml:"address" default:"localhost:55678" comment:"OpenCensus agent address"`
	// ServiceName is the name of the process.
	ServiceName string `toml:"serviceName" comment:"Service name"`
	UseTLS      bool   `toml:"useTLS" default:"false" comment:"Enable TLS listener"`
	TLS         struct {
		CertificatePath              string `toml:"certificatePath" default:"" comment:"Certificate path"`
		PrivateKeyPath               string `toml:"privateKeyPath" default:"" comment:"Private Key path"`
		CACertificatePath            string `toml:"caCertificatePath" default:"" comment:"CA Certificate Path"`
		ClientAuthenticationRequired bool   `toml:"clientAuthenticationRequired" default:"false" comment:"Force client authentication"`
	} `toml:"TLS" comment:"TLS Socket settings"`
}

// Validate checks that the configuration is valid.
func (c Config) Validate() error {
	if c.ServiceName == "" {
		return xerrors.New("ocagent: service name must not be blank")
	}
	return nil
}
