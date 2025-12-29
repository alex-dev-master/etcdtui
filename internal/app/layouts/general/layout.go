package general

import (
	"context"
	"log"

	generalActions "github.com/alexandr/etcdtui/internal/app/actions/general"
	"github.com/alexandr/etcdtui/internal/ui/panels/details"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type General struct {
	generalActions *generalActions.General
	app            *tview.Application
	rootFlex       *tview.Flex
	mainFlex       *tview.Flex
	contentFlex    *tview.Flex // Flex for main content + debug panel
	ctx            context.Context
	inEditMode     bool // Flag to disable global input capture during edit
}

func NewGeneral(app *tview.Application) *General {
	return &General{generalActions: generalActions.NewGeneral(), app: app}
}

func (g *General) Render(ctx context.Context) (err error) {
	g.ctx = ctx

	if err = g.generalActions.Exec(ctx); err != nil {
		return err
	}

	// Setup action callback for details panel
	g.generalActions.GetDetailsPanel().SetActionCallback(g.handleDetailsAction)

	// Setup tab callback for details panel to switch focus back to keys
	g.generalActions.GetDetailsPanel().SetTabCallback(func() {
		g.app.SetFocus(g.generalActions.GetKeysPanel().GetTree())
	})

	// Layout: tree on left, details on right
	g.mainFlex = tview.NewFlex().
		AddItem(g.generalActions.GetKeysPanel().GetTree(), 0, 1, true).
		AddItem(g.generalActions.GetDetailsPanel().GetView(), 0, 2, false)

	// Content flex: main view + optional debug panel (side by side)
	g.contentFlex = tview.NewFlex().
		AddItem(g.mainFlex, 0, 1, true)

	// Main layout with status bar at bottom
	g.rootFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(g.contentFlex, 0, 1, true).
		AddItem(g.generalActions.GetStatusBarPanel().GetView(), 1, 0, false)

	g.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return g.GetInputCapture(ctx, event)
	})

	if err := g.app.SetRoot(g.rootFlex, true).EnableMouse(false).Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (g *General) GetRootFlex() *tview.Flex {
	return g.rootFlex
}

func (g *General) GetInputCapture(ctx context.Context, event *tcell.EventKey) *tcell.EventKey {
	// When in edit mode, only handle Ctrl+C, pass everything else through
	if g.inEditMode {
		if event.Key() == tcell.KeyCtrlC {
			g.app.Stop()
			return nil
		}
		return event
	}

	// Don't intercept keys when we're not in the main view
	// This allows forms, modals, etc. to handle their own input
	if g.app.GetFocus() != g.generalActions.GetKeysPanel().GetTree() &&
		g.app.GetFocus() != g.generalActions.GetDetailsPanel().GetForm() {
		// Only handle Ctrl+C to quit from anywhere
		if event.Key() == tcell.KeyCtrlC {
			g.app.Stop()
			return nil
		}
		return event
	}

	switch event.Key() {
	case tcell.KeyCtrlC:
		g.app.Stop()
		return nil
	case tcell.KeyF1:
		// Toggle debug panel
		g.toggleDebugPanel()
		return nil
	case tcell.KeyTab:
		// Switch focus from tree to details panel
		current := g.app.GetFocus()
		if current == g.generalActions.GetKeysPanel().GetTree() {
			// Switch to details panel buttons - only if buttons are shown
			g.app.SetFocus(g.generalActions.GetDetailsPanel().GetForm())
			return nil
		} else if current == g.generalActions.GetStatusBarPanel().GetView() {
			g.app.SetFocus(g.generalActions.GetKeysPanel().GetTree())
			return nil
		}
		// Let other views handle Tab themselves
		return event
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
	case 'w':
		g.handleWatch(ctx)
		return nil
	case 'c':
		g.handleCopy(ctx)
		return nil
	}

	return event
}

func (g *General) showHelp(app *tview.Application, rootView tview.Primitive) {
	modal := tview.NewModal().
		SetText(`etcdtui - Interactive TUI for etcd

Navigation:
  ↓/↑ or j/k  - Navigate tree
  Enter       - Expand/collapse or select key
  Tab         - Switch panels (Keys ↔ Details)
  ←/→ or h/l  - Navigate buttons (in Details panel)

Edit Form Navigation:
  ↓/↑         - Navigate between fields and buttons
  Enter       - Activate button or new line in text
  ESC         - Cancel and close

Quick Actions:
  e           - Edit key/value
  d           - Delete key
  r           - Refresh keys
  w           - Watch mode (TODO)
  c           - Copy value (TODO)
  n           - New key (TODO)
  /           - Search (TODO)

Debug:
  F1          - Toggle debug panel

Other:
  ?           - Show this help
  q or Ctrl+C - Quit
  ESC         - Cancel/Close modal`).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(rootView, true)
		})

	app.SetRoot(modal, true)
}

// handleDelete shows confirmation modal and deletes the selected key
func (g *General) handleDelete(ctx context.Context) {
	kv := g.generalActions.GetCurrentKey()
	if kv == nil {
		g.generalActions.SetStatusBarText("[yellow]No key selected")
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
	kv := g.generalActions.GetCurrentKey()
	if kv == nil {
		g.generalActions.SetStatusBarText("[yellow]No key selected")
		g.generalActions.GetDebugPanel().LogWarn("Edit attempted with no key selected")
		return
	}

	g.generalActions.GetDebugPanel().LogInfo("Opening edit form for key: %s", kv.Key)

	// Enable edit mode to bypass global input capture
	g.inEditMode = true

	// Helper to close form and restore main view
	closeForm := func() {
		g.inEditMode = false
		g.app.SetRoot(g.rootFlex, true)
	}

	// Create form for editing
	form := tview.NewForm()

	// Add Key input field
	form.AddInputField("Key", kv.Key, 50, nil, nil)

	// Add Value text area
	form.AddTextArea("Value", kv.Value, 50, 5, 0, nil)

	form.AddButton("Save", func() {
		newKey := form.GetFormItemByLabel("Key").(*tview.InputField).GetText()
		newValue := form.GetFormItemByLabel("Value").(*tview.TextArea).GetText()

		g.generalActions.GetDebugPanel().LogDebug("Save button clicked - Key: %s, Value length: %d", newKey, len(newValue))

		// If key changed, delete old and create new
		if newKey != kv.Key {
			g.generalActions.GetDebugPanel().LogInfo("Key renamed from '%s' to '%s'", kv.Key, newKey)
			if err := g.generalActions.DeleteKey(ctx, kv.Key); err != nil {
				g.generalActions.SetStatusBarText("[red]Failed to delete old key:[white] " + err.Error())
				g.generalActions.GetDebugPanel().LogError("Failed to delete old key '%s': %v", kv.Key, err)
				closeForm()
				return
			}
		}

		if err := g.generalActions.PutKey(ctx, newKey, newValue); err != nil {
			g.generalActions.SetStatusBarText("[red]Failed to save:[white] " + err.Error())
			g.generalActions.GetDebugPanel().LogError("Failed to save key '%s': %v", newKey, err)
			closeForm()
			return
		}

		g.generalActions.GetDebugPanel().LogInfo("Successfully saved key: %s", newKey)

		// Refresh details for the updated key
		if err := g.generalActions.RefreshKeyDetails(ctx, newKey); err != nil {
			g.generalActions.SetStatusBarText("[yellow]Saved but failed to refresh details:[white] " + err.Error())
			g.generalActions.GetDebugPanel().LogWarn("Saved but failed to refresh details: %v", err)
		} else {
			g.generalActions.SetStatusBarText("[green]Saved:[white] " + newKey)
		}

		closeForm()
	})

	form.AddButton("Cancel", func() {
		closeForm()
	})

	// Setup ESC to close the form
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			closeForm()
			return nil
		}
		return event
	})

	form.SetBorder(true).SetTitle(" Edit Key (Tab to navigate, ESC cancel) ").SetTitleAlign(tview.AlignLeft)
	form.SetCancelFunc(closeForm)

	// Set root and focus on first form field
	g.app.SetRoot(form, true)
	form.SetFocus(0)
}

// handleWatch shows watch mode for the selected key
func (g *General) handleWatch(ctx context.Context) {
	kv := g.generalActions.GetCurrentKey()
	if kv == nil {
		g.generalActions.SetStatusBarText("[yellow]No key selected")
		return
	}

	// TODO: Implement watch functionality
	g.generalActions.SetStatusBarText("[yellow]Watch mode for " + kv.Key + " [not implemented yet]")
}

// handleCopy copies the selected key value to clipboard
func (g *General) handleCopy(ctx context.Context) {
	kv := g.generalActions.GetCurrentKey()
	if kv == nil {
		g.generalActions.SetStatusBarText("[yellow]No key selected")
		return
	}

	// TODO: Implement clipboard copy functionality
	g.generalActions.SetStatusBarText("[yellow]Copy to clipboard [not implemented yet]")
}

// handleDetailsAction handles actions from the details panel buttons
func (g *General) handleDetailsAction(action details.ActionType) {
	switch action {
	case details.ActionEdit:
		// Don't restore focus - edit form handles its own focus
		g.handleEdit(g.ctx)
	case details.ActionDelete:
		g.handleDelete(g.ctx)
		g.app.SetFocus(g.generalActions.GetKeysPanel().GetTree())
	case details.ActionWatch:
		g.handleWatch(g.ctx)
		g.app.SetFocus(g.generalActions.GetKeysPanel().GetTree())
	case details.ActionCopy:
		g.handleCopy(g.ctx)
		g.app.SetFocus(g.generalActions.GetKeysPanel().GetTree())
	}
}

// toggleDebugPanel shows/hides the debug panel
func (g *General) toggleDebugPanel() {
	debugPanel := g.generalActions.GetDebugPanel()

	if debugPanel.IsVisible() {
		// Log before hiding
		debugPanel.LogInfo("Debug panel hidden")
		// Hide debug panel
		g.contentFlex.RemoveItem(debugPanel.GetView())
		debugPanel.SetVisible(false)
		g.generalActions.SetStatusBarText("[yellow]Debug panel hidden (F1 to show)")
	} else {
		// Show debug panel on the right side
		g.contentFlex.AddItem(debugPanel.GetView(), 0, 1, false)
		debugPanel.SetVisible(true)
		g.generalActions.SetStatusBarText("[yellow]Debug panel visible (F1 to hide)")
		debugPanel.LogInfo("Debug panel shown")
		debugPanel.LogInfo("Application started in debug mode")
	}
}
