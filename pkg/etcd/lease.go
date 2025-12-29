package client

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// LeaseInfo represents information about a lease
type LeaseInfo struct {
	ID  int64
	TTL int64
}

// PutWithTTL stores a key-value pair with TTL
func (c *Client) PutWithTTL(ctx context.Context, key, value string, ttl time.Duration) (*LeaseInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Create lease
	lease, err := c.client.Grant(ctx, int64(ttl.Seconds()))
	if err != nil {
		return nil, fmt.Errorf("failed to grant lease: %w", err)
	}

	// Put with lease
	_, err = c.client.Put(ctx, key, value, clientv3.WithLease(lease.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to put key with lease: %w", err)
	}

	return &LeaseInfo{
		ID:  int64(lease.ID),
		TTL: lease.TTL,
	}, nil
}

// KeepAlive keeps a lease alive by renewing it periodically
func (c *Client) KeepAlive(ctx context.Context, leaseID int64) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	ch, err := c.client.KeepAlive(ctx, clientv3.LeaseID(leaseID))
	if err != nil {
		return nil, fmt.Errorf("failed to keep alive lease %d: %w", leaseID, err)
	}
	return ch, nil
}

// RevokeLease revokes a lease, deleting all associated keys
func (c *Client) RevokeLease(ctx context.Context, leaseID int64) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.Revoke(ctx, clientv3.LeaseID(leaseID))
	if err != nil {
		return fmt.Errorf("failed to revoke lease %d: %w", leaseID, err)
	}
	return nil
}

// GetLeaseInfo retrieves information about a lease
func (c *Client) GetLeaseInfo(ctx context.Context, leaseID int64) (*LeaseInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.TimeToLive(ctx, clientv3.LeaseID(leaseID))
	if err != nil {
		return nil, fmt.Errorf("failed to get lease info %d: %w", leaseID, err)
	}

	return &LeaseInfo{
		ID:  int64(resp.ID),
		TTL: resp.TTL,
	}, nil
}

// ListLeases returns all active leases
func (c *Client) ListLeases(ctx context.Context) ([]int64, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.Leases(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list leases: %w", err)
	}

	leases := make([]int64, 0, len(resp.Leases))
	for _, lease := range resp.Leases {
		leases = append(leases, int64(lease.ID))
	}

	return leases, nil
}
