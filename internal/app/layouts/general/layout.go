package general

import (
	"context"
	"log"

	generalActions "github.com/alexandr/etcdtui/internal/app/actions/general"
	client "github.com/alexandr/etcdtui/pkg/etcd"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type General struct {
	generalActions *generalActions.General
	app            *tview.Application
	rootFlex       *tview.Flex
}

func NewGeneral(app *tview.Application) *General {
	return &General{generalActions: generalActions.NewGeneral(), app: app}
}

func (g *General) Render(ctx context.Context) (err error) {
	if err = g.generalActions.Exec(ctx); err != nil {
		return err
	}
	// Layout: tree on left, details on right
	mainFlex := tview.NewFlex().
		AddItem(g.generalActions.GetKeysPanel().GetTree(), 0, 1, true).
		AddItem(g.generalActions.GetDetailsPanel().GetView(), 0, 2, false)

	// Main layout with status bar at bottom
	g.rootFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 1, true).
		AddItem(g.generalActions.GetStatusBarPanel().GetView(), 1, 0, false)

	g.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return g.GetInputCapture(ctx, event)
	})

	if err := g.app.SetRoot(g.rootFlex, true).EnableMouse(true).Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (g *General) GetRootFlex() *tview.Flex {
	return g.rootFlex
}

func (g *General) GetInputCapture(ctx context.Context, event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlC:
		g.app.Stop()
		return nil
	}

	switch event.Rune() {
	case 'q':
		g.app.Stop()
		return nil
	case '?':
		g.showHelp(g.app, g.rootFlex)
		return nil
	case '/':
		// TODO: Implement search
		g.generalActions.SetStatusBarText("[yellow]Search:[white] [not implemented yet]")
		return nil
	case 'n':
		// TODO: Implement new key
		g.generalActions.SetStatusBarText("[yellow]New key:[white] [not implemented yet]")
		return nil
	case 'r':
		if err := g.generalActions.RefreshKeys(ctx); err != nil {
			g.generalActions.SetStatusBarText("[red]Failed to refresh:[white] " + err.Error())
		}
		return nil
	case 'd':
		g.handleDelete(ctx)
		return nil
	case 'e':
		g.handleEdit(ctx)
		return nil
	}

	return event
}

func (g *General) showHelp(app *tview.Application, rootView tview.Primitive) {
	modal := tview.NewModal().
		SetText(`etcdtui - Interactive TUI for etcd

Keyboard Shortcuts:
  ↓/↑ or j/k  - Navigate tree
  Enter       - Expand/collapse or edit
  n           - New key
  d           - Delete key
  e           - Edit value
  w           - Watch mode
  l           - Locks dashboard
  r           - Refresh
  /           - Search
  ?           - Show this help
  q or Ctrl+C - Quit

Navigation:
  Tab         - Switch between panels
  Mouse       - Click and scroll support`).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(rootView, true)
		})

	app.SetRoot(modal, true)
}

// handleDelete shows confirmation modal and deletes the selected key
func (g *General) handleDelete(ctx context.Context) {
	tree := g.generalActions.GetKeysPanel().GetTree()
	node := tree.GetCurrentNode()
	if node == nil {
		return
	}

	reference := node.GetReference()
	if reference == nil {
		g.generalActions.SetStatusBarText("[yellow]Cannot delete directory node")
		return
	}

	kv, ok := reference.(*client.KeyValue)
	if !ok {
		return
	}

	// Show confirmation modal
	modal := tview.NewModal().
		SetText("Delete key: " + kv.Key + "?").
		AddButtons([]string{"Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			g.app.SetRoot(g.rootFlex, true)
			if buttonLabel == "Delete" {
				if err := g.generalActions.DeleteKey(ctx, kv.Key); err != nil {
					g.generalActions.SetStatusBarText("[red]Failed to delete:[white] " + err.Error())
				} else {
					g.generalActions.SetStatusBarText("[green]Deleted:[white] " + kv.Key)
				}
			}
		})

	g.app.SetRoot(modal, true)
}

// handleEdit shows edit modal for the selected key
func (g *General) handleEdit(ctx context.Context) {
	tree := g.generalActions.GetKeysPanel().GetTree()
	node := tree.GetCurrentNode()
	if node == nil {
		return
	}

	reference := node.GetReference()
	if reference == nil {
		g.generalActions.SetStatusBarText("[yellow]Cannot edit directory node")
		return
	}

	kv, ok := reference.(*client.KeyValue)
	if !ok {
		return
	}

	// Create form for editing
	form := tview.NewForm()
	form.AddInputField("Key", kv.Key, 50, nil, nil)
	form.AddTextArea("Value", kv.Value, 50, 5, 0, nil)

	form.AddButton("Save", func() {
		keyField := form.GetFormItemByLabel("Key").(*tview.InputField)
		valueField := form.GetFormItemByLabel("Value").(*tview.TextArea)

		newKey := keyField.GetText()
		newValue := valueField.GetText()

		// If key changed, delete old and create new
		if newKey != kv.Key {
			if err := g.generalActions.DeleteKey(ctx, kv.Key); err != nil {
				g.generalActions.SetStatusBarText("[red]Failed to delete old key:[white] " + err.Error())
				g.app.SetRoot(g.rootFlex, true)
				return
			}
		}

		if err := g.generalActions.PutKey(ctx, newKey, newValue); err != nil {
			g.generalActions.SetStatusBarText("[red]Failed to save:[white] " + err.Error())
		} else {
			g.generalActions.SetStatusBarText("[green]Saved:[white] " + newKey)
		}
		g.app.SetRoot(g.rootFlex, true)
	})

	form.AddButton("Cancel", func() {
		g.app.SetRoot(g.rootFlex, true)
	})

	form.SetBorder(true).SetTitle(" Edit Key ").SetTitleAlign(tview.AlignLeft)
	form.SetCancelFunc(func() {
		g.app.SetRoot(g.rootFlex, true)
	})

	g.app.SetRoot(form, true)
}
