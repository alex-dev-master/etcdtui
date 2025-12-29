package general

import (
	"context"

	"github.com/alexandr/etcdtui/internal/ui/panels/details"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// HandleEdit shows edit form for the selected key.
func (s *State) HandleEdit(ctx context.Context) {
	kv := s.GetCurrentKey()
	if kv == nil {
		s.SetStatusBarText("[yellow]No key selected")
		s.debugPanel.LogWarn("Edit attempted with no key selected")
		return
	}

	s.debugPanel.LogInfo("Opening edit form for key: %s", kv.Key)

	// Enable edit mode to bypass global input capture
	s.SetEditMode(true)

	// Helper to close form and restore main view
	closeForm := func() {
		s.SetEditMode(false)
		s.app.SetRoot(s.rootFlex, true)
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

		s.debugPanel.LogDebug("Save button clicked - Key: %s, Value length: %d", newKey, len(newValue))

		// If key changed, delete old and create new
		if newKey != kv.Key {
			s.debugPanel.LogInfo("Key renamed from '%s' to '%s'", kv.Key, newKey)
			if err := s.DeleteKey(ctx, kv.Key); err != nil {
				s.SetStatusBarText("[red]Failed to delete old key:[white] " + err.Error())
				s.debugPanel.LogError("Failed to delete old key '%s': %v", kv.Key, err)
				closeForm()
				return
			}
		}

		if err := s.PutKey(ctx, newKey, newValue); err != nil {
			s.SetStatusBarText("[red]Failed to save:[white] " + err.Error())
			s.debugPanel.LogError("Failed to save key '%s': %v", newKey, err)
			closeForm()
			return
		}

		s.debugPanel.LogInfo("Successfully saved key: %s", newKey)

		// Refresh details for the updated key
		if err := s.RefreshKeyDetails(ctx, newKey); err != nil {
			s.SetStatusBarText("[yellow]Saved but failed to refresh details:[white] " + err.Error())
			s.debugPanel.LogWarn("Saved but failed to refresh details: %v", err)
		} else {
			s.SetStatusBarText("[green]Saved:[white] " + newKey)
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
	s.app.SetRoot(form, true)
	form.SetFocus(0)
}

// HandleDelete shows confirmation modal and deletes the selected key.
func (s *State) HandleDelete(ctx context.Context) {
	kv := s.GetCurrentKey()
	if kv == nil {
		s.SetStatusBarText("[yellow]No key selected")
		return
	}

	modal := tview.NewModal().
		SetText("Delete key: " + kv.Key + "?").
		AddButtons([]string{"Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			s.app.SetRoot(s.rootFlex, true)
			if buttonLabel == "Delete" {
				if err := s.DeleteKey(ctx, kv.Key); err != nil {
					s.SetStatusBarText("[red]Failed to delete:[white] " + err.Error())
				} else {
					s.SetStatusBarText("[green]Deleted:[white] " + kv.Key)
				}
			}
		})

	s.app.SetRoot(modal, true)
}

// HandleWatch shows watch mode for the selected key.
func (s *State) HandleWatch(ctx context.Context) {
	kv := s.GetCurrentKey()
	if kv == nil {
		s.SetStatusBarText("[yellow]No key selected")
		return
	}

	// TODO: Implement watch functionality
	s.SetStatusBarText("[yellow]Watch mode for " + kv.Key + " [not implemented yet]")
}

// HandleCopy copies the selected key value to clipboard.
func (s *State) HandleCopy(ctx context.Context) {
	kv := s.GetCurrentKey()
	if kv == nil {
		s.SetStatusBarText("[yellow]No key selected")
		return
	}

	// TODO: Implement clipboard copy functionality
	s.SetStatusBarText("[yellow]Copy to clipboard [not implemented yet]")
}

// HandleDetailsAction handles actions from the details panel buttons.
func (s *State) HandleDetailsAction(ctx context.Context, action details.ActionType) {
	switch action {
	case details.ActionEdit:
		// Don't restore focus - edit form handles its own focus
		s.HandleEdit(ctx)
	case details.ActionDelete:
		s.HandleDelete(ctx)
		s.app.SetFocus(s.keysPanel.GetTree())
	case details.ActionWatch:
		s.HandleWatch(ctx)
		s.app.SetFocus(s.keysPanel.GetTree())
	case details.ActionCopy:
		s.HandleCopy(ctx)
		s.app.SetFocus(s.keysPanel.GetTree())
	}
}

// ToggleDebugPanel shows/hides the debug panel.
func (s *State) ToggleDebugPanel(contentFlex *tview.Flex) {
	if s.debugPanel.IsVisible() {
		s.debugPanel.LogInfo("Debug panel hidden")
		contentFlex.RemoveItem(s.debugPanel.GetView())
		s.debugPanel.SetVisible(false)
		s.SetStatusBarText("[yellow]Debug panel hidden (F1 to show)")
	} else {
		contentFlex.AddItem(s.debugPanel.GetView(), 0, 1, false)
		s.debugPanel.SetVisible(true)
		s.SetStatusBarText("[yellow]Debug panel visible (F1 to hide)")
		s.debugPanel.LogInfo("Debug panel shown")
		s.debugPanel.LogInfo("Application started in debug mode")
	}
}

// ShowHelp displays the help modal.
func (s *State) ShowHelp() {
	modal := tview.NewModal().
		SetText(`etcdtui - Interactive TUI for etcd

Navigation:
  ↓/↑ or j/k  - Navigate tree
  Enter       - Expand/collapse or select key
  Tab         - Switch panels (Keys ↔ Details)
  ←/→ or h/l  - Navigate buttons (in Details panel)

Edit Form Navigation:
  Tab         - Navigate between fields
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
			s.app.SetRoot(s.rootFlex, true)
		})

	s.app.SetRoot(modal, true)
}
