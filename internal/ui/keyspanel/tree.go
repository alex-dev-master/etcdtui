package keyspanel

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TreeNode struct {
	tree *tview.TreeView
	once sync.Once
}

func NewTreeNode() *TreeNode {
	return &TreeNode{
		tree: tview.NewTreeView(),
	}
}

// Draw Создаём дерево ключей (левая панель)
func (t *TreeNode) Draw() {
	t.once.Do(t.draw)
}

func (t *TreeNode) draw() {
	root := tview.NewTreeNode("etcd").
		SetColor(tcell.ColorYellow).
		SetExpanded(true)
	t.tree.SetRoot(root).SetCurrentNode(root)

	// Добавляем примеры данных для демонстрации
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

	// Добавляем секцию config
	config := tview.NewTreeNode("config").
		SetExpanded(true)
	root.AddChild(config)

	dbURL := tview.NewTreeNode("database-url").
		SetReference("/config/database-url")
	config.AddChild(dbURL)

	redisURL := tview.NewTreeNode("redis-url").
		SetReference("/config/redis-url")
	config.AddChild(redisURL)

	// Добавляем секцию locks
	locks := tview.NewTreeNode("locks")
	root.AddChild(locks)

	paymentLock := tview.NewTreeNode("payment [TTL: 28s]").
		SetReference("/locks/payment").
		SetColor(tcell.ColorRed)
	locks.AddChild(paymentLock)

	t.tree.SetBorder(true).SetTitle(" Keys ")
}

func (t *TreeNode) GetTree() *tview.TreeView {
	return t.tree
}
