package general

import (
	"context"
	"fmt"
	"time"

	"github.com/alexandr/etcdtui/internal/app/connection"
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
	connManager    *connection.Manager
	ctx            context.Context
}

func NewGeneral() *General {
	return &General{
		keysPanel:      keys.New(),
		detailsPanel:   details.New(),
		statusBarPanel: statusbar.New(),
		connManager:    connection.NewManager(),
		ctx:            context.Background(),
	}
}

// Exec initializes the UI and establishes etcd connection
func (g *General) Exec() {
	g.keysPanel.Draw()
	g.detailsPanel.Draw()
	g.statusBarPanel.Draw()

	// Try to connect to etcd
	if err := g.connManager.ConnectDefault(); err != nil {
		g.SetStatusBarText(fmt.Sprintf("[red]Not connected:[white] %v | [yellow]Press [green]c[white] to configure connection", err))
	} else {
		// Set client for keys panel
		g.keysPanel.SetClient(g.connManager.GetClient())

		// Load keys
		if err := g.keysPanel.LoadKeys(g.ctx); err != nil {
			g.SetStatusBarText(fmt.Sprintf("[yellow]Connected but failed to load keys:[white] %v", err))
		} else {
			g.updateStatusBar()
		}
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
			g.showKeyDetails(kv)
		}
	})
}

// showKeyDetails displays detailed information about a key
func (g *General) showKeyDetails(kv *client.KeyValue) {
	detailsText := fmt.Sprintf("[yellow]Key:[white] %s\n\n", kv.Key)
	detailsText += fmt.Sprintf("[yellow]Value:[white]\n%s\n\n", kv.Value)
	detailsText += fmt.Sprintf("[yellow]Create Revision:[white] %d\n", kv.CreateRevision)
	detailsText += fmt.Sprintf("[yellow]Mod Revision:[white] %d\n", kv.ModRevision)
	detailsText += fmt.Sprintf("[yellow]Version:[white] %d\n", kv.Version)

	if kv.Lease > 0 {
		// Get lease info
		if cli := g.connManager.GetClient(); cli != nil {
			if leaseInfo, err := cli.GetLeaseInfo(g.ctx, kv.Lease); err == nil {
				detailsText += fmt.Sprintf("[yellow]TTL:[white] %d seconds\n", leaseInfo.TTL)
			} else {
				detailsText += fmt.Sprintf("[yellow]Lease:[white] %d\n", kv.Lease)
			}
		}
	} else {
		detailsText += "[yellow]TTL:[white] ∞\n"
	}

	detailsText += "\n[green][e][white] Edit  [green][d][white] Delete  [green][w][white] Watch  [green][c][white] Copy"

	g.detailsPanel.SetText(detailsText)
}

// RefreshKeys reloads keys from etcd
func (g *General) RefreshKeys() error {
	if !g.connManager.IsConnected() {
		return fmt.Errorf("not connected to etcd")
	}

	if err := g.keysPanel.Refresh(g.ctx); err != nil {
		g.SetStatusBarText(fmt.Sprintf("[red]Failed to refresh:[white] %v", err))
		return err
	}

	g.updateStatusBar()
	return nil
}

// updateStatusBar updates status bar with current stats
func (g *General) updateStatusBar() {
	cli := g.connManager.GetClient()
	if cli == nil {
		g.SetStatusBarText("[yellow]Status:[white] Not connected | [green][c][white] Connect")
		return
	}

	// Get key count
	count, err := cli.GetKeyCount(g.ctx)
	if err != nil {
		count = 0
	}

	// Get cluster status
	status, err := cli.GetClusterStatus(g.ctx)
	leaderInfo := "unknown"
	if err == nil && status != nil && status.Leader != "" {
		leaderInfo = status.Leader
	}

	statusText := fmt.Sprintf("[green]Connected[white] | Leader: [cyan]%s[white] | Keys: [yellow]%d[white] | [green][/] Search  [n] New  [r] Refresh  [q] Quit  [?] Help",
		leaderInfo, count)

	g.SetStatusBarText(statusText)
}

// DeleteKey deletes a key from etcd
func (g *General) DeleteKey(key string) error {
	cli := g.connManager.GetClient()
	if cli == nil {
		return fmt.Errorf("not connected to etcd")
	}

	if err := cli.Delete(g.ctx, key); err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	// Refresh tree
	return g.RefreshKeys()
}

// PutKey creates or updates a key
func (g *General) PutKey(key, value string) error {
	cli := g.connManager.GetClient()
	if cli == nil {
		return fmt.Errorf("not connected to etcd")
	}

	if err := cli.Put(g.ctx, key, value); err != nil {
		return fmt.Errorf("failed to put key: %w", err)
	}

	// Refresh tree
	return g.RefreshKeys()
}

// PutKeyWithTTL creates or updates a key with TTL
func (g *General) PutKeyWithTTL(key, value string, ttl time.Duration) error {
	cli := g.connManager.GetClient()
	if cli == nil {
		return fmt.Errorf("not connected to etcd")
	}

	if _, err := cli.PutWithTTL(g.ctx, key, value, ttl); err != nil {
		return fmt.Errorf("failed to put key with TTL: %w", err)
	}

	// Refresh tree
	return g.RefreshKeys()
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
func (g *General) GetConnectionManager() *connection.Manager {
	return g.connManager
}
