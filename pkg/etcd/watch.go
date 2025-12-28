package client

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// EventType represents the type of watch event
type EventType int

const (
	EventTypePut EventType = iota
	EventTypeDelete
)

// WatchEvent represents a change event from etcd
type WatchEvent struct {
	Type           EventType
	Key            string
	Value          string
	PrevValue      string
	CreateRevision int64
	ModRevision    int64
	Version        int64
}

// WatchCallback is called when a watch event occurs
type WatchCallback func(*WatchEvent)

// Watch starts watching a key or prefix for changes
func (c *Client) Watch(ctx context.Context, key string, callback WatchCallback) error {
	watchChan := c.client.Watch(ctx, key)
	return c.processWatchEvents(watchChan, callback)
}

// WatchPrefix starts watching all keys with a given prefix
func (c *Client) WatchPrefix(ctx context.Context, prefix string, callback WatchCallback) error {
	watchChan := c.client.Watch(ctx, prefix, clientv3.WithPrefix())
	return c.processWatchEvents(watchChan, callback)
}

// WatchFromRevision starts watching from a specific revision
func (c *Client) WatchFromRevision(ctx context.Context, key string, revision int64, callback WatchCallback) error {
	watchChan := c.client.Watch(ctx, key, clientv3.WithRev(revision))
	return c.processWatchEvents(watchChan, callback)
}

// processWatchEvents processes events from a watch channel
func (c *Client) processWatchEvents(watchChan clientv3.WatchChan, callback WatchCallback) error {
	for watchResp := range watchChan {
		if watchResp.Err() != nil {
			return fmt.Errorf("watch error: %w", watchResp.Err())
		}

		for _, event := range watchResp.Events {
			watchEvent := &WatchEvent{
				Key:            string(event.Kv.Key),
				Value:          string(event.Kv.Value),
				CreateRevision: event.Kv.CreateRevision,
				ModRevision:    event.Kv.ModRevision,
				Version:        event.Kv.Version,
			}

			switch event.Type {
			case clientv3.EventTypePut:
				watchEvent.Type = EventTypePut
				if event.PrevKv != nil {
					watchEvent.PrevValue = string(event.PrevKv.Value)
				}
			case clientv3.EventTypeDelete:
				watchEvent.Type = EventTypeDelete
				if event.PrevKv != nil {
					watchEvent.PrevValue = string(event.PrevKv.Value)
				}
			}

			callback(watchEvent)
		}
	}
	return nil
}
