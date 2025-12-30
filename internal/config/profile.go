package config

import (
	"encoding/base64"
	"strings"
	"time"

	client "github.com/alex-dev-master/etcdtui/pkg/etcd"
)

// Profile represents a connection profile for etcd
type Profile struct {
	// Name is the profile identifier
	Name string `yaml:"name" mapstructure:"name"`

	// Endpoints is a list of etcd server addresses
	Endpoints []string `yaml:"endpoints" mapstructure:"endpoints"`

	// Username for authentication (optional)
	Username string `yaml:"username,omitempty" mapstructure:"username"`

	// Password for authentication (base64 encoded, optional)
	// Format: "base64:ENCODED_PASSWORD" or plain text
	Password string `yaml:"password,omitempty" mapstructure:"password"`

	// TLS configuration (optional)
	TLS *TLSProfile `yaml:"tls,omitempty" mapstructure:"tls"`

	// Default marks this profile as the default connection
	Default bool `yaml:"default,omitempty" mapstructure:"default"`
}

// TLSProfile represents TLS configuration in a profile
type TLSProfile struct {
	// Enabled enables TLS
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`

	// CAFile is the path to CA certificate
	CAFile string `yaml:"ca_file,omitempty" mapstructure:"ca_file"`

	// CertFile is the path to client certificate
	CertFile string `yaml:"cert_file,omitempty" mapstructure:"cert_file"`

	// KeyFile is the path to client key
	KeyFile string `yaml:"key_file,omitempty" mapstructure:"key_file"`

	// InsecureSkipVerify skips TLS verification (not recommended)
	InsecureSkipVerify bool `yaml:"insecure_skip_verify,omitempty" mapstructure:"insecure_skip_verify"`
}

// DecodePassword decodes the password from storage format
// Supports: "base64:ENCODED" or plain text
func (p *Profile) DecodePassword() string {
	if p.Password == "" {
		return ""
	}

	if strings.HasPrefix(p.Password, "base64:") {
		encoded := strings.TrimPrefix(p.Password, "base64:")
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return p.Password // return as-is if decode fails
		}
		return string(decoded)
	}

	return p.Password
}

// EncodePassword encodes password for storage (base64)
func EncodePassword(password string) string {
	if password == "" {
		return ""
	}
	return "base64:" + base64.StdEncoding.EncodeToString([]byte(password))
}

// ToClientConfig converts Profile to client.Config
func (p *Profile) ToClientConfig() *client.Config {
	cfg := &client.Config{
		Endpoints:      p.Endpoints,
		Username:       p.Username,
		Password:       p.DecodePassword(),
		DialTimeout:    5 * time.Second,
		RequestTimeout: 5 * time.Second,
	}

	if p.TLS != nil && p.TLS.Enabled {
		cfg.TLS = &client.TLSConfig{
			Enabled:            true,
			CAFile:             p.TLS.CAFile,
			CertFile:           p.TLS.CertFile,
			KeyFile:            p.TLS.KeyFile,
			InsecureSkipVerify: p.TLS.InsecureSkipVerify,
		}
	}

	return cfg
}

// Validate checks if the profile has required fields
func (p *Profile) Validate() error {
	if p.Name == "" {
		return ErrProfileNameRequired
	}
	if len(p.Endpoints) == 0 {
		return ErrEndpointsRequired
	}
	return nil
}

// HasAuth returns true if profile has authentication configured
func (p *Profile) HasAuth() bool {
	return p.Username != ""
}

// HasTLS returns true if profile has TLS configured
func (p *Profile) HasTLS() bool {
	return p.TLS != nil && p.TLS.Enabled
}

// DisplayString returns a string representation for UI display
func (p *Profile) DisplayString() string {
	var parts []string
	parts = append(parts, p.Name)

	if len(p.Endpoints) > 0 {
		parts = append(parts, "("+p.Endpoints[0]+")")
	}

	var flags []string
	if p.HasAuth() {
		flags = append(flags, "auth")
	}
	if p.HasTLS() {
		flags = append(flags, "tls")
	}
	if p.Default {
		flags = append(flags, "default")
	}

	if len(flags) > 0 {
		parts = append(parts, "["+strings.Join(flags, ", ")+"]")
	}

	return strings.Join(parts, " ")
}
