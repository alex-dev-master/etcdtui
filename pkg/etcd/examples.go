package client

// This file contains usage examples for the etcd client

/*
Example 1: Basic CRUD operations

	// Create client
	cfg := DefaultConfig()
	cfg.Endpoints = []string{"localhost:2379"}
	client, err := New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	// Put a key
	err = client.Put(ctx, "/config/app", "value1")

	// Get a key
	kv, err := client.Get(ctx, "/config/app")
	fmt.Println(kv.Value) // "value1"

	// List keys with prefix
	kvs, err := client.List(ctx, "/config/")
	for _, kv := range kvs {
		fmt.Printf("%s = %s\n", kv.Key, kv.Value)
	}

	// Delete a key
	err = client.Delete(ctx, "/config/app")

Example 2: Watch for changes

	client, _ := New(DefaultConfig())
	defer client.Close()

	ctx := context.Background()

	// Watch a single key
	err := client.Watch(ctx, "/config/app", func(event *WatchEvent) {
		if event.Type == EventTypePut {
			fmt.Printf("Key updated: %s = %s\n", event.Key, event.Value)
		} else {
			fmt.Printf("Key deleted: %s\n", event.Key)
		}
	})

	// Watch all keys with prefix
	err = client.WatchPrefix(ctx, "/config/", func(event *WatchEvent) {
		fmt.Printf("Event: %v\n", event)
	})

Example 3: Lease and TTL

	client, _ := New(DefaultConfig())
	defer client.Close()

	ctx := context.Background()

	// Put with TTL
	lease, err := client.PutWithTTL(ctx, "/session/user1", "active", 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created lease %d with TTL %d seconds\n", lease.ID, lease.TTL)

	// Keep lease alive
	keepAliveCh, err := client.KeepAlive(ctx, lease.ID)
	go func() {
		for range keepAliveCh {
			// Lease renewed
		}
	}()

	// Revoke lease (deletes associated keys)
	err = client.RevokeLease(ctx, lease.ID)

Example 4: Distributed Lock

	client, _ := New(DefaultConfig())
	defer client.Close()

	ctx := context.Background()

	// Acquire lock
	lock, err := client.AcquireLock(ctx, "/locks/resource1", 30*time.Second)
	if err != nil {
		log.Fatal("Failed to acquire lock:", err)
	}

	// Do work with exclusive access
	fmt.Println("Lock acquired, doing work...")
	time.Sleep(5 * time.Second)

	// Release lock
	err = client.ReleaseLock(ctx, lock)
	if err != nil {
		log.Fatal("Failed to release lock:", err)
	}

Example 5: Transactions

	client, _ := New(DefaultConfig())
	defer client.Close()

	ctx := context.Background()

	// Compare and swap
	success, err := client.CompareAndSwap(ctx, "/config/version", "1.0", "1.1")
	if success {
		fmt.Println("Version updated successfully")
	} else {
		fmt.Println("Version mismatch, update failed")
	}

	// Create only if doesn't exist
	created, err := client.CreateIfNotExists(ctx, "/config/init", "initialized")
	if created {
		fmt.Println("Config initialized")
	}

Example 6: Cluster Status

	client, _ := New(DefaultConfig())
	defer client.Close()

	ctx := context.Background()

	status, err := client.GetClusterStatus(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Leader: %s\n", status.Leader)
	fmt.Printf("Members: %d\n", len(status.Members))
	for _, member := range status.Members {
		fmt.Printf("  - %s (%s)\n", member.Name, member.Endpoint)
	}

	// Get key count
	count, err := client.GetKeyCount(ctx)
	fmt.Printf("Total keys: %d\n", count)

Example 7: Authentication with TLS

	cfg := &Config{
		Endpoints: []string{"https://etcd.example.com:2379"},
		Username:  "admin",
		Password:  "secret",
		TLS: &TLSConfig{
			Enabled:  true,
			CertFile: "/path/to/client.crt",
			KeyFile:  "/path/to/client.key",
			CAFile:   "/path/to/ca.crt",
		},
		DialTimeout:    5 * time.Second,
		RequestTimeout: 10 * time.Second,
	}

	client, err := New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Check health
	err = client.HealthCheck(context.Background())
	if err != nil {
		log.Fatal("Cluster unhealthy:", err)
	}
*/
