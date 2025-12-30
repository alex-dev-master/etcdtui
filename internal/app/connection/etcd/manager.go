package etcd

import (
	"context"
	"fmt"
	"sync"

	client "github.com/alex-dev-master/etcdtui/pkg/etcd"
)

// Manager manages etcd connection
type Manager struct {
	client *client.Client
	config *client.Config
	mu     sync.RWMutex
}

// NewManager creates a new connection manager
func NewManager() *Manager {
	return &Manager{}
}

// Connect establishes connection to etcd with the given config
func (m *Manager) Connect(cfg *client.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Close existing connection if any
	if m.client != nil {
		m.client.Close()
	}

	// Create new client
	cli, err := client.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to etcd: %w", err)
	}

	// Test connection
	ctx := context.Background()
	if err := cli.HealthCheck(ctx); err != nil {
		cli.Close()
		return fmt.Errorf("etcd health check failed: %w", err)
	}

	m.client = cli
	m.config = cfg
	return nil
}

// ConnectDefault connects using default configuration
func (m *Manager) ConnectDefault() error {
	return m.Connect(client.DefaultConfig())
}

// Disconnect closes the connection
func (m *Manager) Disconnect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.client != nil {
		err := m.client.Close()
		m.client = nil
		return err
	}
	return nil
}

// GetClient returns the etcd client (thread-safe)
func (m *Manager) GetClient() *client.Client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client
}

// IsConnected checks if connection is established
func (m *Manager) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client != nil
}

// GetConfig returns current configuration
func (m *Manager) GetConfig() *client.Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}
