package client

import (
	"context"
	"fmt"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// ClusterStatus represents the status of etcd cluster
type ClusterStatus struct {
	Members   []*MemberStatus
	Leader    string
	IsHealthy bool
}

// MemberStatus represents status of a single etcd member
type MemberStatus struct {
	ID       uint64
	Name     string
	Endpoint string
	IsLeader bool
}

// GetClusterStatus retrieves the status of the etcd cluster
func (c *Client) GetClusterStatus(ctx context.Context) (*ClusterStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Get member list
	memberList, err := c.client.MemberList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get member list: %w", err)
	}

	// Get cluster status
	status := &ClusterStatus{
		Members:   make([]*MemberStatus, 0, len(memberList.Members)),
		IsHealthy: true,
	}

	for _, member := range memberList.Members {
		memberStatus := &MemberStatus{
			ID:       member.ID,
			Name:     member.Name,
			IsLeader: member.ID == memberList.Header.MemberId,
		}

		if len(member.ClientURLs) > 0 {
			memberStatus.Endpoint = member.ClientURLs[0]
		}

		if memberStatus.IsLeader {
			status.Leader = memberStatus.Name
		}

		status.Members = append(status.Members, memberStatus)
	}

	return status, nil
}

// GetKeyCount returns the total number of keys in etcd
func (c *Client) GetKeyCount(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.Get(ctx, "", clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return 0, fmt.Errorf("failed to count keys: %w", err)
	}

	return resp.Count, nil
}

// GetKeyCountWithPrefix returns the number of keys with a specific prefix
func (c *Client) GetKeyCountWithPrefix(ctx context.Context, prefix string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return 0, fmt.Errorf("failed to count keys with prefix: %w", err)
	}

	return resp.Count, nil
}

// BuildTree builds a hierarchical tree structure from etcd keys
func BuildTree(keys []*KeyValue) map[string]interface{} {
	tree := make(map[string]interface{})

	for _, kv := range keys {
		parts := strings.Split(strings.Trim(kv.Key, "/"), "/")
		current := tree

		for i, part := range parts {
			if i == len(parts)-1 {
				// Leaf node - store the value
				current[part] = kv
			} else {
				// Branch node - create nested map
				if _, exists := current[part]; !exists {
					current[part] = make(map[string]interface{})
				}
				current = current[part].(map[string]interface{})
			}
		}
	}

	return tree
}

// CompactHistory compacts etcd history up to a given revision
func (c *Client) CompactHistory(ctx context.Context, revision int64) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.Compact(ctx, revision)
	if err != nil {
		return fmt.Errorf("failed to compact history: %w", err)
	}

	return nil
}

// Snapshot creates a snapshot of the etcd data
// Note: This is a placeholder. In real usage, you'd want to stream this to a file
func (c *Client) Snapshot(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout*5) // Longer timeout for snapshot
	defer cancel()

	// Get snapshot reader
	_, err := c.client.Snapshot(ctx)
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	// Note: In real usage, you would read from the snapshot and write to a file
	// For now, just return an informative error
	return fmt.Errorf("snapshot functionality requires streaming to file - use etcd client Snapshot() directly")
}
