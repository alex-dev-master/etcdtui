package client

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// User represents an etcd user
type User struct {
	Name  string
	Roles []string
}

// Role represents an etcd role
type Role struct {
	Name        string
	Permissions []*Permission
}

// Permission represents a key permission
type Permission struct {
	Key      string
	RangeEnd string
	PermType PermissionType
}

// PermissionType represents the type of permission
type PermissionType int

const (
	PermissionRead PermissionType = iota
	PermissionWrite
	PermissionReadWrite
)

// CreateUser creates a new etcd user
func (c *Client) CreateUser(ctx context.Context, username, password string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.UserAdd(ctx, username, password)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// DeleteUser deletes an etcd user
func (c *Client) DeleteUser(ctx context.Context, username string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.UserDelete(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// ChangePassword changes user password
func (c *Client) ChangePassword(ctx context.Context, username, newPassword string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.UserChangePassword(ctx, username, newPassword)
	if err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}
	return nil
}

// GrantRole grants a role to a user
func (c *Client) GrantRole(ctx context.Context, username, role string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.UserGrantRole(ctx, username, role)
	if err != nil {
		return fmt.Errorf("failed to grant role: %w", err)
	}
	return nil
}

// RevokeRole revokes a role from a user
func (c *Client) RevokeRole(ctx context.Context, username, role string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.UserRevokeRole(ctx, username, role)
	if err != nil {
		return fmt.Errorf("failed to revoke role: %w", err)
	}
	return nil
}

// ListUsers returns all users
func (c *Client) ListUsers(ctx context.Context) ([]*User, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.UserList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*User, 0, len(resp.Users))
	for _, username := range resp.Users {
		// Get user details
		userResp, err := c.client.UserGet(ctx, username)
		if err != nil {
			continue
		}

		users = append(users, &User{
			Name:  username,
			Roles: userResp.Roles,
		})
	}

	return users, nil
}

// CreateRole creates a new role
func (c *Client) CreateRole(ctx context.Context, roleName string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.RoleAdd(ctx, roleName)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	return nil
}

// DeleteRole deletes a role
func (c *Client) DeleteRole(ctx context.Context, roleName string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.RoleDelete(ctx, roleName)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}

// GrantPermission grants permission to a role
func (c *Client) GrantPermission(ctx context.Context, roleName, key, rangeEnd string, permType PermissionType) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	var permTypeEtcd clientv3.PermissionType
	switch permType {
	case PermissionRead:
		permTypeEtcd = clientv3.PermissionType(clientv3.PermRead)
	case PermissionWrite:
		permTypeEtcd = clientv3.PermissionType(clientv3.PermWrite)
	case PermissionReadWrite:
		permTypeEtcd = clientv3.PermissionType(clientv3.PermReadWrite)
	}

	_, err := c.client.RoleGrantPermission(ctx, roleName, key, rangeEnd, permTypeEtcd)
	if err != nil {
		return fmt.Errorf("failed to grant permission: %w", err)
	}
	return nil
}

// RevokePermission revokes permission from a role
func (c *Client) RevokePermission(ctx context.Context, roleName, key, rangeEnd string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.RoleRevokePermission(ctx, roleName, key, rangeEnd)
	if err != nil {
		return fmt.Errorf("failed to revoke permission: %w", err)
	}
	return nil
}

// EnableAuth enables authentication
func (c *Client) EnableAuth(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.AuthEnable(ctx)
	if err != nil {
		return fmt.Errorf("failed to enable auth: %w", err)
	}
	return nil
}

// DisableAuth disables authentication
func (c *Client) DisableAuth(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.AuthDisable(ctx)
	if err != nil {
		return fmt.Errorf("failed to disable auth: %w", err)
	}
	return nil
}
