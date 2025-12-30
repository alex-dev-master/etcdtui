package keys

import (
	"context"
	"sort"
	"strings"
	"sync"

	client "github.com/alex-dev-master/etcdtui/pkg/etcd"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// treeNode represents a node that can be both a key and a directory
type treeNode struct {
	kv       *client.KeyValue // nil if this is just a directory
	children map[string]*treeNode
}

func newTreeNode() *treeNode {
	return &treeNode{
		children: make(map[string]*treeNode),
	}
}

// Panel represents the keys tree panel (left side)
type Panel struct {
	tree *tview.TreeView
	once sync.Once
}

// New creates a new keys panel
func New() *Panel {
	return &Panel{
		tree: tview.NewTreeView(),
	}
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
func (p *Panel) LoadKeys(ctx context.Context, kvs []*client.KeyValue) error {
	// Clear existing tree
	root := p.tree.GetRoot()
	root.ClearChildren()

	// Build hierarchical tree from flat keys
	tree := buildHierarchy(kvs)

	// Add nodes to tview tree
	p.addNodes(root, tree)

	return nil
}

// buildHierarchy converts flat key list to hierarchical structure
func buildHierarchy(kvs []*client.KeyValue) *treeNode {
	root := newTreeNode()

	for _, kv := range kvs {
		parts := strings.Split(strings.Trim(kv.Key, "/"), "/")
		current := root

		for i, part := range parts {
			if part == "" {
				continue
			}

			// Create child node if doesn't exist
			if _, exists := current.children[part]; !exists {
				current.children[part] = newTreeNode()
			}

			if i == len(parts)-1 {
				// This is the actual key - set the kv
				current.children[part].kv = kv
			}

			current = current.children[part]
		}
	}

	return root
}

// addNodes recursively adds nodes to the tview tree
func (p *Panel) addNodes(parent *tview.TreeNode, node *treeNode) {
	// Sort children keys for consistent display
	keys := make([]string, 0, len(node.children))
	for k := range node.children {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		child := node.children[key]
		hasChildren := len(child.children) > 0
		isKey := child.kv != nil

		// Build display text with indicator
		displayText := key
		if hasChildren {
			displayText = "▶ " + key
		}

		var treeNode *tview.TreeNode

		if isKey {
			// This is an actual key (may also have children)
			treeNode = tview.NewTreeNode(displayText).
				SetReference(child.kv).
				SetColor(tcell.ColorGreen).
				SetExpanded(false)
		} else {
			// This is just a directory (no value)
			treeNode = tview.NewTreeNode(displayText).
				SetColor(tcell.ColorAqua).
				SetExpanded(false)
		}

		parent.AddChild(treeNode)

		// Recursively add children
		if hasChildren {
			p.addNodes(treeNode, child)
		}
	}
}

// GetTree returns the underlying TreeView
func (p *Panel) GetTree() *tview.TreeView {
	return p.tree
}
