package general

import (
	"context"
	"fmt"

	"github.com/alexandr/etcdtui/internal/app/connection/etcd"
	"github.com/alexandr/etcdtui/internal/ui/panels/details"
	"github.com/alexandr/etcdtui/internal/ui/panels/keys"
	"github.com/alexandr/etcdtui/internal/ui/panels/statusbar"
	client "github.com/alexandr/etcdtui/pkg/etcd"
	"github.com/rivo/tview"
)

type General struct {
	keysPanel      *keys.Panel
	detailsPanel   *details.Panel
	statusBarPanel *statusbar.Panel
	connManager    *etcd.Manager
	currentKey     *client.KeyValue // Currently selected key
}

func NewGeneral() *General {
	return &General{
		keysPanel:      keys.New(),
		detailsPanel:   details.New(),
		statusBarPanel: statusbar.New(),
		connManager:    etcd.NewManager(),
	}
}

// Exec initializes the UI and establishes etcd connection
func (g *General) Exec(ctx context.Context) (err error) {
	g.keysPanel.Draw()
	g.detailsPanel.Draw()
	g.statusBarPanel.Draw()

	if err = g.connManager.ConnectDefault(); err != nil {
		g.SetStatusBarText(fmt.Sprintf("[red]Not connected:[white] %v | [yellow]Press [green]c[white] to configure connection", err))
		return nil
	}

	if err = g.seedingKeysData(ctx); err != nil {
		return err
	}

	// Handle tree node selection
	g.keysPanel.GetTree().SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			// Toggle node expansion
			node.SetExpanded(!node.IsExpanded())
			return
		}

		// Show key details
		if kv, ok := reference.(*client.KeyValue); ok {
			g.showKeyDetails(ctx, kv)
		}
	})

	return nil
}

func (g *General) seedingKeysData(ctx context.Context) (err error) {
	var kvs []*client.KeyValue
	kvs, err = g.connManager.GetClient().List(ctx, "")
	if err != nil {
		return err
	}

	// Load keys
	if err = g.keysPanel.LoadKeys(ctx, kvs); err != nil {
		g.SetStatusBarText(fmt.Sprintf("[yellow]Connected but failed to load keys:[white] %v", err))
	} else {
		g.updateStatusBar(ctx)
	}
	return nil
}

// showKeyDetails displays detailed information about a key
func (g *General) showKeyDetails(ctx context.Context, kv *client.KeyValue) {
	// Store current key
	g.currentKey = kv

	detailsText := fmt.Sprintf("[yellow]Key:[white] %s\n\n", kv.Key)
	detailsText += fmt.Sprintf("[yellow]Value:[white]\n%s\n\n", kv.Value)
	detailsText += fmt.Sprintf("[yellow]Create Revision:[white] %d\n", kv.CreateRevision)
	detailsText += fmt.Sprintf("[yellow]Mod Revision:[white] %d\n", kv.ModRevision)
	detailsText += fmt.Sprintf("[yellow]Version:[white] %d\n", kv.Version)

	if kv.Lease > 0 {
		// Get lease info
		if cli := g.connManager.GetClient(); cli != nil {
			if leaseInfo, err := cli.GetLeaseInfo(ctx, kv.Lease); err == nil {
				detailsText += fmt.Sprintf("[yellow]TTL:[white] %d seconds\n", leaseInfo.TTL)
			} else {
				detailsText += fmt.Sprintf("[yellow]Lease:[white] %d\n", kv.Lease)
			}
		}
	} else {
		detailsText += "[yellow]TTL:[white] ∞\n"
	}

	g.detailsPanel.SetText(detailsText)
	g.detailsPanel.ShowButtons()
}

// GetCurrentKey returns the currently selected key
func (g *General) GetCurrentKey() *client.KeyValue {
	return g.currentKey
}

// RefreshKeys reloads keys from etcd
func (g *General) RefreshKeys(ctx context.Context) error {
	if !g.connManager.IsConnected() {
		return fmt.Errorf("not connected to etcd")
	}

	if err := g.seedingKeysData(ctx); err != nil {
		g.SetStatusBarText(fmt.Sprintf("[red]Failed to refresh:[white] %v", err))
		return err
	}

	g.updateStatusBar(ctx)
	return nil
}

// updateStatusBar updates status bar with current stats
func (g *General) updateStatusBar(ctx context.Context) {
	cli := g.connManager.GetClient()
	if cli == nil {
		g.SetStatusBarText("[yellow]Status:[white] Not connected | [green][c][white] Connect")
		return
	}

	// Get key count
	count, err := cli.GetKeyCount(ctx)
	if err != nil {
		count = 0
	}

	// Get cluster status
	status, err := cli.GetClusterStatus(ctx)
	leaderInfo := "unknown"
	if err == nil && status != nil && status.Leader != "" {
		leaderInfo = status.Leader
	}

	statusText := fmt.Sprintf("[green]Connected[white] | Leader: [cyan]%s[white] | Keys: [yellow]%d[white] | [green][/] Search  [n] New  [r] Refresh  [q] Quit  [?] Help",
		leaderInfo, count)

	g.SetStatusBarText(statusText)
}

// DeleteKey deletes a key from etcd
func (g *General) DeleteKey(ctx context.Context, key string) error {
	cli := g.connManager.GetClient()
	if cli == nil {
		return fmt.Errorf("not connected to etcd")
	}

	if err := cli.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	// Refresh tree
	return g.RefreshKeys(ctx)
}

// PutKey creates or updates a key
func (g *General) PutKey(ctx context.Context, key, value string) error {
	cli := g.connManager.GetClient()
	if cli == nil {
		return fmt.Errorf("not connected to etcd")
	}

	if err := cli.Put(ctx, key, value); err != nil {
		return fmt.Errorf("failed to put key: %w", err)
	}

	// Refresh tree
	return g.RefreshKeys(ctx)
}

// RefreshKeyDetails refreshes details for a specific key
func (g *General) RefreshKeyDetails(ctx context.Context, key string) error {
	cli := g.connManager.GetClient()
	if cli == nil {
		return fmt.Errorf("not connected to etcd")
	}

	// Get updated key value
	kv, err := cli.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to get key: %w", err)
	}

	// Update details panel
	g.showKeyDetails(ctx, kv)
	return nil
}

func (g *General) GetKeysPanel() *keys.Panel {
	return g.keysPanel
}

func (g *General) GetDetailsPanel() *details.Panel {
	return g.detailsPanel
}

func (g *General) GetStatusBarPanel() *statusbar.Panel {
	return g.statusBarPanel
}

func (g *General) SetStatusBarText(text string) {
	g.statusBarPanel.SetText(text)
}

// GetConnectionManager returns the connection manager
func (g *General) GetConnectionManager() *etcd.Manager {
	return g.connManager
}
