package general

import (
	"context"
	"fmt"

	client "github.com/alexandr/etcdtui/pkg/etcd"
	"github.com/rivo/tview"
)

// InitConnection establishes connection to etcd and loads initial data.
func (s *State) InitConnection(ctx context.Context) error {
	s.keysPanel.Draw()
	s.detailsPanel.Draw()
	s.statusBarPanel.Draw()

	if err := s.connManager.ConnectDefault(); err != nil {
		s.SetStatusBarText(fmt.Sprintf("[red]Not connected:[white] %v | [yellow]Press [green]c[white] to configure connection", err))
		return nil
	}

	if err := s.seedingKeysData(ctx); err != nil {
		return err
	}

	// Enter key - toggle expand/collapse
	s.keysPanel.GetTree().SetSelectedFunc(func(node *tview.TreeNode) {
		children := node.GetChildren()
		if len(children) > 0 {
			node.SetExpanded(!node.IsExpanded())
		}
	})

	// Navigation (arrow keys) - show details when moving to a node
	s.keysPanel.GetTree().SetChangedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference != nil {
			if kv, ok := reference.(*client.KeyValue); ok {
				s.showKeyDetails(ctx, kv)
			}
		} else {
			// Clear details for directory-only nodes
			s.detailsPanel.SetText("[yellow]Directory[white]\n\nSelect a key to view details")
			s.detailsPanel.HideButtons()
			s.currentKey = nil
		}
	})

	return nil
}

// seedingKeysData loads keys from etcd into the keys panel.
func (s *State) seedingKeysData(ctx context.Context) error {
	kvs, err := s.connManager.GetClient().List(ctx, "")
	if err != nil {
		return err
	}

	if err = s.keysPanel.LoadKeys(ctx, kvs); err != nil {
		s.SetStatusBarText(fmt.Sprintf("[yellow]Connected but failed to load keys:[white] %v", err))
	} else {
		s.updateStatusBar(ctx)
	}
	return nil
}

// showKeyDetails displays detailed information about a key.
func (s *State) showKeyDetails(ctx context.Context, kv *client.KeyValue) {
	s.currentKey = kv

	detailsText := fmt.Sprintf("[yellow]Key:[white] %s\n\n", kv.Key)
	detailsText += fmt.Sprintf("[yellow]Value:[white]\n%s\n\n", kv.Value)
	detailsText += fmt.Sprintf("[yellow]Create Revision:[white] %d\n", kv.CreateRevision)
	detailsText += fmt.Sprintf("[yellow]Mod Revision:[white] %d\n", kv.ModRevision)
	detailsText += fmt.Sprintf("[yellow]Version:[white] %d\n", kv.Version)

	if kv.Lease > 0 {
		if cli := s.connManager.GetClient(); cli != nil {
			if leaseInfo, err := cli.GetLeaseInfo(ctx, kv.Lease); err == nil {
				detailsText += fmt.Sprintf("[yellow]TTL:[white] %d seconds\n", leaseInfo.TTL)
			} else {
				detailsText += fmt.Sprintf("[yellow]Lease:[white] %d\n", kv.Lease)
			}
		}
	} else {
		detailsText += "[yellow]TTL:[white] ∞\n"
	}

	s.detailsPanel.SetText(detailsText)
	s.detailsPanel.ShowButtons()
}

// updateStatusBar updates status bar with current stats.
func (s *State) updateStatusBar(ctx context.Context) {
	cli := s.connManager.GetClient()
	if cli == nil {
		s.SetStatusBarText("[yellow]Status:[white] Not connected | [green][c][white] Connect")
		return
	}

	count, err := cli.GetKeyCount(ctx)
	if err != nil {
		count = 0
	}

	status, err := cli.GetClusterStatus(ctx)
	leaderInfo := "unknown"
	if err == nil && status != nil && status.Leader != "" {
		leaderInfo = status.Leader
	}

	statusText := fmt.Sprintf("[green]Connected[-] | Leader: [cyan]%s[-] | Keys: [yellow]%d[-] | [green::b]/[-::-] Search  [green::b]n[-::-] New  [green::b]r[-::-] Refresh  [green::b]q[-::-] Quit  [green::b]?[-::-] Help",
		leaderInfo, count)

	s.SetStatusBarText(statusText)
}

// RefreshKeys reloads keys from etcd.
func (s *State) RefreshKeys(ctx context.Context) error {
	if !s.connManager.IsConnected() {
		return fmt.Errorf("not connected to etcd")
	}

	if err := s.seedingKeysData(ctx); err != nil {
		s.SetStatusBarText(fmt.Sprintf("[red]Failed to refresh:[white] %v", err))
		return err
	}

	s.updateStatusBar(ctx)
	return nil
}

// DeleteKey deletes a key from etcd.
func (s *State) DeleteKey(ctx context.Context, key string) error {
	cli := s.connManager.GetClient()
	if cli == nil {
		return fmt.Errorf("not connected to etcd")
	}

	if err := cli.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	return s.RefreshKeys(ctx)
}

// PutKey creates or updates a key.
func (s *State) PutKey(ctx context.Context, key, value string) error {
	cli := s.connManager.GetClient()
	if cli == nil {
		return fmt.Errorf("not connected to etcd")
	}

	if err := cli.Put(ctx, key, value); err != nil {
		return fmt.Errorf("failed to put key: %w", err)
	}

	return s.RefreshKeys(ctx)
}

// RefreshKeyDetails refreshes details for a specific key.
func (s *State) RefreshKeyDetails(ctx context.Context, key string) error {
	cli := s.connManager.GetClient()
	if cli == nil {
		return fmt.Errorf("not connected to etcd")
	}

	kv, err := cli.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to get key: %w", err)
	}

	s.showKeyDetails(ctx, kv)
	return nil
}

// SearchByPrefix searches keys by prefix and updates the tree.
func (s *State) SearchByPrefix(ctx context.Context, prefix string) error {
	cli := s.connManager.GetClient()
	if cli == nil {
		return fmt.Errorf("not connected to etcd")
	}

	kvs, err := cli.List(ctx, prefix)
	if err != nil {
		return fmt.Errorf("failed to search keys: %w", err)
	}

	if err = s.keysPanel.LoadKeys(ctx, kvs); err != nil {
		return fmt.Errorf("failed to load search results: %w", err)
	}

	// Clear current key selection and details
	s.currentKey = nil
	s.detailsPanel.SetText(fmt.Sprintf("[yellow]Search results for:[white] %s\n\n[cyan]%d keys found[-]", prefix, len(kvs)))
	s.detailsPanel.HideButtons()

	return nil
}
