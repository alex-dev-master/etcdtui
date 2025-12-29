package client

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// KeyValue represents a key-value pair with metadata
type KeyValue struct {
	Key            string
	Value          string
	CreateRevision int64
	ModRevision    int64
	Version        int64
	Lease          int64
}

// Get retrieves a single key from etcd
func (c *Client) Get(ctx context.Context, key string) (*KeyValue, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	kv := resp.Kvs[0]
	return &KeyValue{
		Key:            string(kv.Key),
		Value:          string(kv.Value),
		CreateRevision: kv.CreateRevision,
		ModRevision:    kv.ModRevision,
		Version:        kv.Version,
		Lease:          int64(kv.Lease),
	}, nil
}

// Put stores a key-value pair in etcd
func (c *Client) Put(ctx context.Context, key, value string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.Put(ctx, key, value)
	if err != nil {
		return fmt.Errorf("failed to put key %s: %w", key, err)
	}
	return nil
}

// Delete removes a key from etcd
func (c *Client) Delete(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// List retrieves all keys with the given prefix
func (c *Client) List(ctx context.Context, prefix string) ([]*KeyValue, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to list keys with prefix %s: %w", prefix, err)
	}

	kvs := make([]*KeyValue, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		kvs = append(kvs, &KeyValue{
			Key:            string(kv.Key),
			Value:          string(kv.Value),
			CreateRevision: kv.CreateRevision,
			ModRevision:    kv.ModRevision,
			Version:        kv.Version,
			Lease:          int64(kv.Lease),
		})
	}

	return kvs, nil
}

// GetWithRevision retrieves a key at a specific revision
func (c *Client) GetWithRevision(ctx context.Context, key string, revision int64) (*KeyValue, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.Get(ctx, key, clientv3.WithRev(revision))
	if err != nil {
		return nil, fmt.Errorf("failed to get key %s at revision %d: %w", key, revision, err)
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("key not found at revision %d: %s", revision, key)
	}

	kv := resp.Kvs[0]
	return &KeyValue{
		Key:            string(kv.Key),
		Value:          string(kv.Value),
		CreateRevision: kv.CreateRevision,
		ModRevision:    kv.ModRevision,
		Version:        kv.Version,
		Lease:          int64(kv.Lease),
	}, nil
}

// DeletePrefix removes all keys with the given prefix
func (c *Client) DeletePrefix(ctx context.Context, prefix string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.Delete(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return 0, fmt.Errorf("failed to delete prefix %s: %w", prefix, err)
	}

	return resp.Deleted, nil
}
