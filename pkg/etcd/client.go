package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Client wraps etcd client with additional functionality
type Client struct {
	client  *clientv3.Client
	config  *Config
	timeout time.Duration
}

// New creates a new etcd client with the given configuration
func New(cfg *Config) (*Client, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	etcdConfig := clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.DialTimeout,
	}

	// Configure authentication
	if cfg.Username != "" {
		etcdConfig.Username = cfg.Username
		etcdConfig.Password = cfg.Password
	}

	// Configure TLS
	if cfg.TLS != nil && cfg.TLS.Enabled {
		tlsConfig, err := loadTLSConfig(cfg.TLS)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS config: %w", err)
		}
		etcdConfig.TLS = tlsConfig
	}

	// Create etcd client
	cli, err := clientv3.New(etcdConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	return &Client{
		client:  cli,
		config:  cfg,
		timeout: cfg.RequestTimeout,
	}, nil
}

// Close closes the etcd client connection
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// HealthCheck checks if etcd cluster is healthy
func (c *Client) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.Get(ctx, "/health-check", clientv3.WithLimit(1))
	return err
}

// loadTLSConfig creates TLS configuration from files
func loadTLSConfig(cfg *TLSConfig) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: cfg.InsecureSkipVerify,
	}

	if cfg.CertFile != "" && cfg.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load client cert/key: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	if cfg.CAFile != "" {
		caCert, err := os.ReadFile(cfg.CAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA file: %w", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append CA cert")
		}
		tlsConfig.RootCAs = caCertPool
	}

	return tlsConfig, nil
}
