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
	ctx            context.Context
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

	// Layout: tree on left, details on right
	g.mainFlex = tview.NewFlex().
		AddItem(g.generalActions.GetKeysPanel().GetTree(), 0, 1, true).
		AddItem(g.generalActions.GetDetailsPanel().GetView(), 0, 2, false)

	// Main layout with status bar at bottom
	g.rootFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(g.mainFlex, 0, 1, true).
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
	case tcell.KeyTab:
		// Switch focus between tree and details panel
		current := g.app.GetFocus()
		if current == g.generalActions.GetKeysPanel().GetTree() {
			// Switch to details panel buttons
			g.app.SetFocus(g.generalActions.GetDetailsPanel().GetForm())
		} else {
			// Switch back to tree
			g.app.SetFocus(g.generalActions.GetKeysPanel().GetTree())
		}
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

Keyboard Shortcuts:
  ↓/↑ or j/k  - Navigate tree
  Enter       - Expand/collapse or select key
  Tab         - Switch between panels
  ←/→ or h/l  - Navigate buttons in details panel

Actions:
  e           - Edit key/value
  d           - Delete key
  w           - Watch mode (TODO)
  c           - Copy value (TODO)
  n           - New key (TODO)
  r           - Refresh keys
  /           - Search (TODO)

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
		return
	}

	// Create form for editing
	form := tview.NewForm()

	// Add Key input field with navigation
	keyField := tview.NewInputField().
		SetLabel("Key").
		SetText(kv.Key).
		SetFieldWidth(50)

	keyField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyDown {
			// Move to Value field
			form.SetFocus(1) // Index 1 = Value TextArea
			return nil
		}
		return event
	})

	form.AddFormItem(keyField)

	// Add Value text area with navigation
	valueField := tview.NewTextArea().
		SetLabel("Value").
		SetSize(5, 50)

	// Set text separately to ensure it's displayed
	if kv.Value != "" {
		valueField.SetText(kv.Value, true)
	}

	valueField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyUp {
			// Check if cursor is on first line
			row, _, _, _ := valueField.GetCursor()
			if row == 0 {
				// Move to Key field
				form.SetFocus(0) // Index 0 = Key InputField
				return nil
			}
		} else if event.Key() == tcell.KeyDown {
			// Check if cursor is on last line
			row, _, _, _ := valueField.GetCursor()
			text := valueField.GetText()

			// Count lines in text
			lines := 1
			for _, ch := range text {
				if ch == '\n' {
					lines++
				}
			}

			// If on last line, move to buttons
			if row >= lines-1 {
				form.SetFocus(2) // Index 2 = First button
				return nil
			}
		}
		return event
	})

	form.AddFormItem(valueField)

	form.AddButton("Save", func() {
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
			g.app.SetRoot(g.rootFlex, true)
			return
		}

		// Refresh details for the updated key
		if err := g.generalActions.RefreshKeyDetails(ctx, newKey); err != nil {
			g.generalActions.SetStatusBarText("[yellow]Saved but failed to refresh details:[white] " + err.Error())
		} else {
			g.generalActions.SetStatusBarText("[green]Saved:[white] " + newKey)
		}

		g.app.SetRoot(g.rootFlex, true)
	})

	form.AddButton("Cancel", func() {
		g.app.SetRoot(g.rootFlex, true)
	})

	// Setup ESC to close the form and arrow navigation for buttons
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			g.app.SetRoot(g.rootFlex, true)
			return nil
		}

		// Check if we're on buttons
		focusedItem, _ := form.GetFocusedItemIndex()
		if focusedItem >= 2 { // Buttons start at index 2
			if event.Key() == tcell.KeyUp {
				// Move from buttons to Value field
				form.SetFocus(1)
				return nil
			}
		}

		// Pass through all other keys to form fields
		return event
	})

	form.SetBorder(true).SetTitle(" Edit Key (ESC to cancel) ").SetTitleAlign(tview.AlignLeft)
	form.SetCancelFunc(func() {
		g.app.SetRoot(g.rootFlex, true)
	})

	// Set root and focus on the Key field
	g.app.SetRoot(form, true)
	g.app.SetFocus(keyField)
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
	// Switch back to tree focus after action
	defer g.app.SetFocus(g.generalActions.GetKeysPanel().GetTree())

	switch action {
	case details.ActionEdit:
		g.handleEdit(g.ctx)
	case details.ActionDelete:
		g.handleDelete(g.ctx)
	case details.ActionWatch:
		g.handleWatch(g.ctx)
	case details.ActionCopy:
		g.handleCopy(g.ctx)
	}
}
