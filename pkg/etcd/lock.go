package client

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// Lock represents a distributed lock
type Lock struct {
	session *concurrency.Session
	mutex   *concurrency.Mutex
}

// AcquireLock acquires a distributed lock with the given key and TTL
func (c *Client) AcquireLock(ctx context.Context, key string, ttl time.Duration) (*Lock, error) {
	// Create a session with TTL
	session, err := concurrency.NewSession(c.client, concurrency.WithTTL(int(ttl.Seconds())))
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Create a mutex on the key
	mutex := concurrency.NewMutex(session, key)

	// Try to acquire the lock
	if err := mutex.Lock(ctx); err != nil {
		_ = session.Close()
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return &Lock{
		session: session,
		mutex:   mutex,
	}, nil
}

// ReleaseLock releases a distributed lock
func (c *Client) ReleaseLock(ctx context.Context, lock *Lock) error {
	if lock == nil {
		return nil
	}

	// Unlock the mutex
	if err := lock.mutex.Unlock(ctx); err != nil {
		return fmt.Errorf("failed to unlock: %w", err)
	}

	// Close the session
	if err := lock.session.Close(); err != nil {
		return fmt.Errorf("failed to close session: %w", err)
	}

	return nil
}

// TryLock attempts to acquire a lock without blocking
func (c *Client) TryLock(ctx context.Context, key string, ttl time.Duration) (*Lock, error) {
	session, err := concurrency.NewSession(c.client, concurrency.WithTTL(int(ttl.Seconds())))
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	mutex := concurrency.NewMutex(session, key)

	// Try to acquire with timeout
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if err := mutex.TryLock(ctx); err != nil {
		_ = session.Close()
		return nil, fmt.Errorf("failed to try lock: %w", err)
	}

	return &Lock{
		session: session,
		mutex:   mutex,
	}, nil
}

// IsLocked checks if a key is currently locked
func (c *Client) IsLocked(ctx context.Context, key string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Check if there are any keys with the lock prefix
	resp, err := c.client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithLimit(1))
	if err != nil {
		return false, fmt.Errorf("failed to check lock: %w", err)
	}

	return len(resp.Kvs) > 0, nil
}
