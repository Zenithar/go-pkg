package tlsconfig

import (
	"crypto/x509"
)

// SystemCertPool returns an new empty cert pool,
// accessing system cert pool
func SystemCertPool() (*x509.CertPool, error) {
	return x509.NewCertPool(), nil
}
