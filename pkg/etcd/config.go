package client

import "time"

// Config represents etcd connection configuration
type Config struct {
	// Endpoints is a list of etcd server addresses
	Endpoints []string

	// Username for authentication (optional)
	Username string

	// Password for authentication (optional)
	Password string

	// TLS configuration
	TLS *TLSConfig

	// Timeouts
	DialTimeout    time.Duration
	RequestTimeout time.Duration
	KeepAlive      time.Duration
}

// TLSConfig represents TLS/SSL configuration
type TLSConfig struct {
	// Enable TLS
	Enabled bool

	// Path to client certificate file
	CertFile string

	// Path to client key file
	KeyFile string

	// Path to CA certificate file
	CAFile string

	// Skip TLS verification (not recommended for production)
	InsecureSkipVerify bool
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Endpoints:      []string{"localhost:2379"},
		DialTimeout:    5 * time.Second,
		RequestTimeout: 10 * time.Second,
		KeepAlive:      30 * time.Second,
	}
}
