package keys

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Panel represents the keys tree panel (left side)
type (
	Panel struct {
		tree *tview.TreeView
		once sync.Once
	}

	Items struct {
	}
)

// New creates a new keys panel
func New() *Panel {
	return &Panel{
		tree: tview.NewTreeView(),
	}
}

// Draw initializes the keys tree with demo data
func (p *Panel) Draw() {
	p.once.Do(p.initialize)
}

func (p *Panel) initialize() {
	root := tview.NewTreeNode("etcd").
		SetColor(tcell.ColorYellow).
		SetExpanded(false)
	p.tree.SetRoot(root).SetCurrentNode(root)

	// Add demo data
	services := tview.NewTreeNode("services").
		SetExpanded(true)
	root.AddChild(services)

	api := tview.NewTreeNode("api").
		SetExpanded(true)
	services.AddChild(api)

	v1Config := tview.NewTreeNode("v1/config").
		SetReference("/services/api/v1/config")
	api.AddChild(v1Config)

	v1Endpoints := tview.NewTreeNode("v1/endpoints").
		SetReference("/services/api/v1/endpoints")
	api.AddChild(v1Endpoints)

	auth := tview.NewTreeNode("auth")
	services.AddChild(auth)

	jwtSecret := tview.NewTreeNode("jwt-secret").
		SetReference("/services/auth/jwt-secret")
	auth.AddChild(jwtSecret)

	// Add config section
	config := tview.NewTreeNode("config").
		SetExpanded(true)
	root.AddChild(config)

	dbURL := tview.NewTreeNode("database-url").
		SetReference("/config/database-url")
	config.AddChild(dbURL)

	redisURL := tview.NewTreeNode("redis-url").
		SetReference("/config/redis-url")
	config.AddChild(redisURL)

	// Add locks section
	locks := tview.NewTreeNode("locks")
	root.AddChild(locks)

	paymentLock := tview.NewTreeNode("payment [TTL: 28s]").
		SetReference("/locks/payment").
		SetColor(tcell.ColorRed)
	locks.AddChild(paymentLock)

	p.tree.SetBorder(true).SetTitle(" Keys ")
}

// GetTree returns the underlying TreeView
func (p *Panel) GetTree() *tview.TreeView {
	return p.tree
}
