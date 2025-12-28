package keys

import (
	"context"
	"strings"
	"sync"

	client "github.com/alexandr/etcdtui/pkg/etcd"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Panel represents the keys tree panel (left side)
type Panel struct {
	tree   *tview.TreeView
	once   sync.Once
	client *client.Client
}

// New creates a new keys panel
func New() *Panel {
	return &Panel{
		tree: tview.NewTreeView(),
	}
}

// SetClient sets the etcd client for the panel
func (p *Panel) SetClient(cli *client.Client) {
	p.client = cli
}

// Draw initializes the keys tree
func (p *Panel) Draw() {
	p.once.Do(p.initialize)
}

func (p *Panel) initialize() {
	root := tview.NewTreeNode("etcd").
		SetColor(tcell.ColorYellow).
		SetExpanded(true)
	p.tree.SetRoot(root).SetCurrentNode(root)

	p.tree.SetBorder(true).SetTitle(" Keys ")
}

// LoadKeys loads all keys from etcd and builds the tree
func (p *Panel) LoadKeys(ctx context.Context) error {
	if p.client == nil {
		return nil // No client, show empty tree
	}

	// Get all keys from etcd
	kvs, err := p.client.List(ctx, "")
	if err != nil {
		return err
	}

	// Clear existing tree
	root := p.tree.GetRoot()
	root.ClearChildren()

	// Build hierarchical tree from flat keys
	tree := buildHierarchy(kvs)

	// Add nodes to tview tree
	for key, value := range tree {
		p.addNode(root, key, value)
	}

	return nil
}

// buildHierarchy converts flat key list to hierarchical structure
func buildHierarchy(kvs []*client.KeyValue) map[string]interface{} {
	tree := make(map[string]interface{})

	for _, kv := range kvs {
		parts := strings.Split(strings.Trim(kv.Key, "/"), "/")
		current := tree

		for i, part := range parts {
			if part == "" {
				continue
			}

			if i == len(parts)-1 {
				// Leaf node - store the KeyValue
				current[part] = kv
			} else {
				// Branch node - create nested map
				if _, exists := current[part]; !exists {
					current[part] = make(map[string]interface{})
				}
				if nested, ok := current[part].(map[string]interface{}); ok {
					current = nested
				}
			}
		}
	}

	return tree
}

// addNode recursively adds nodes to the tree
func (p *Panel) addNode(parent *tview.TreeNode, key string, value interface{}) {
	switch v := value.(type) {
	case *client.KeyValue:
		// Leaf node - actual key
		node := tview.NewTreeNode(key).
			SetReference(v).
			SetColor(tcell.ColorGreen)
		parent.AddChild(node)

	case map[string]interface{}:
		// Branch node - directory
		node := tview.NewTreeNode(key).
			SetColor(tcell.NewRGBColor(0, 255, 255)). // Cyan
			SetExpanded(false)
		parent.AddChild(node)

		// Recursively add children
		for childKey, childValue := range v {
			p.addNode(node, childKey, childValue)
		}
	}
}

// GetTree returns the underlying TreeView
func (p *Panel) GetTree() *tview.TreeView {
	return p.tree
}

// Refresh reloads keys from etcd
func (p *Panel) Refresh(ctx context.Context) error {
	return p.LoadKeys(ctx)
}
