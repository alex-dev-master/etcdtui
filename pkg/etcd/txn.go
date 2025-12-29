package client

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// CompareAndSwap performs an atomic compare-and-swap operation
func (c *Client) CompareAndSwap(ctx context.Context, key, oldValue, newValue string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	txn := c.client.Txn(ctx)

	// If key value equals oldValue, then set to newValue
	resp, err := txn.If(
		clientv3.Compare(clientv3.Value(key), "=", oldValue),
	).Then(
		clientv3.OpPut(key, newValue),
	).Commit()

	if err != nil {
		return false, fmt.Errorf("compare-and-swap failed: %w", err)
	}

	return resp.Succeeded, nil
}

// CreateIfNotExists creates a key only if it doesn't exist
func (c *Client) CreateIfNotExists(ctx context.Context, key, value string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	txn := c.client.Txn(ctx)

	// If key version is 0 (doesn't exist), create it
	resp, err := txn.If(
		clientv3.Compare(clientv3.Version(key), "=", 0),
	).Then(
		clientv3.OpPut(key, value),
	).Commit()

	if err != nil {
		return false, fmt.Errorf("create-if-not-exists failed: %w", err)
	}

	return resp.Succeeded, nil
}

// UpdateIfExists updates a key only if it exists
func (c *Client) UpdateIfExists(ctx context.Context, key, value string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	txn := c.client.Txn(ctx)

	// If key version > 0 (exists), update it
	resp, err := txn.If(
		clientv3.Compare(clientv3.Version(key), ">", 0),
	).Then(
		clientv3.OpPut(key, value),
	).Commit()

	if err != nil {
		return false, fmt.Errorf("update-if-exists failed: %w", err)
	}

	return resp.Succeeded, nil
}
