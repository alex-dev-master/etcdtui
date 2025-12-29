package client

import (
	"context"
	"testing"
	"time"
)

// TestDefaultConfig verifies default configuration
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if len(cfg.Endpoints) != 1 {
		t.Errorf("Expected 1 endpoint, got %d", len(cfg.Endpoints))
	}

	if cfg.Endpoints[0] != "localhost:2379" {
		t.Errorf("Expected localhost:2379, got %s", cfg.Endpoints[0])
	}

	if cfg.DialTimeout != 5*time.Second {
		t.Errorf("Expected 5s dial timeout, got %v", cfg.DialTimeout)
	}

	if cfg.RequestTimeout != 10*time.Second {
		t.Errorf("Expected 10s request timeout, got %v", cfg.RequestTimeout)
	}
}

// TestKeyValueStruct verifies KeyValue structure
func TestKeyValueStruct(t *testing.T) {
	kv := &KeyValue{
		Key:            "/test/key",
		Value:          "test value",
		CreateRevision: 1,
		ModRevision:    2,
		Version:        1,
		Lease:          0,
	}

	if kv.Key != "/test/key" {
		t.Errorf("Expected key '/test/key', got '%s'", kv.Key)
	}

	if kv.Value != "test value" {
		t.Errorf("Expected value 'test value', got '%s'", kv.Value)
	}
}

// TestBuildTree verifies tree building from flat keys
func TestBuildTree(t *testing.T) {
	keys := []*KeyValue{
		{Key: "/services/api/config"},
		{Key: "/services/api/version"},
		{Key: "/services/auth/secret"},
		{Key: "/config/db"},
	}

	tree := BuildTree(keys)

	if tree == nil {
		t.Fatal("Tree should not be nil")
	}

	// Check if services node exists
	if _, exists := tree["services"]; !exists {
		t.Error("Expected 'services' node in tree")
	}

	// Check if config node exists
	if _, exists := tree["config"]; !exists {
		t.Error("Expected 'config' node in tree")
	}
}

// TestEventType verifies event type constants
func TestEventType(t *testing.T) {
	if EventTypePut != 0 {
		t.Error("EventTypePut should be 0")
	}

	if EventTypeDelete != 1 {
		t.Error("EventTypeDelete should be 1")
	}
}

// TestPermissionType verifies permission type constants
func TestPermissionType(t *testing.T) {
	if PermissionRead != 0 {
		t.Error("PermissionRead should be 0")
	}

	if PermissionWrite != 1 {
		t.Error("PermissionWrite should be 1")
	}

	if PermissionReadWrite != 2 {
		t.Error("PermissionReadWrite should be 2")
	}
}

// Example test demonstrating client usage
func ExampleNew() {
	cfg := DefaultConfig()
	client, err := New(cfg)
	if err != nil {
		// Handle error
		return
	}
	defer client.Close()

	ctx := context.Background()

	// Put a key
	_ = client.Put(ctx, "/example/key", "value")

	// Get a key
	kv, err := client.Get(ctx, "/example/key")
	if err == nil {
		_ = kv.Value
	}
}

// Benchmark for BuildTree function
func BenchmarkBuildTree(b *testing.B) {
	keys := make([]*KeyValue, 100)
	for i := 0; i < 100; i++ {
		keys[i] = &KeyValue{
			Key: "/services/api/endpoint" + string(rune(i)),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildTree(keys)
	}
}
