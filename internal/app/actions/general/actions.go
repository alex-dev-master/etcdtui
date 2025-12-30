package general

import (
	"context"
	"fmt"

	"github.com/alex-dev-master/etcdtui/internal/ui/panels/details"
	client "github.com/alex-dev-master/etcdtui/pkg/etcd"
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

	form.AddTextView("Key", kv.Key, 50, 5, true, false)

	// Add Value text area
	form.AddTextArea("Value", kv.Value, 50, 0, 0, nil)

	form.AddButton("Save", func() {
		newValue := form.GetFormItemByLabel("Value").(*tview.TextArea).GetText()
		s.debugPanel.LogDebug("Save button clicked - Key: %s, Value length: %d", kv.Key, len(newValue))

		if err := s.PutKey(ctx, kv.Key, newValue); err != nil {
			s.SetStatusBarText("[red]Failed to save:[white] " + err.Error())
			s.debugPanel.LogError("Failed to save key '%s': %v", kv.Key, err)
			closeForm()
			return
		}

		s.debugPanel.LogInfo("Successfully saved key: %s", kv.Key)

		// Refresh details for the updated key
		if err := s.RefreshKeyDetails(ctx, kv.Key); err != nil {
			s.SetStatusBarText("[yellow]Saved but failed to refresh details:[white] " + err.Error())
			s.debugPanel.LogWarn("Saved but failed to refresh details: %v", err)
		} else {
			s.SetStatusBarText("[green]Saved:[white] " + kv.Key)
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

	// Enable edit mode to bypass global input capture
	s.SetEditMode(true)

	modal := tview.NewModal().
		SetText("Delete key: " + kv.Key + "?").
		AddButtons([]string{"Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			s.SetEditMode(false)
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

	// Cancel any existing watch
	if s.watchCancel != nil {
		s.watchCancel()
	}

	s.SetEditMode(true)

	// Create watch context with cancel
	watchCtx, cancel := context.WithCancel(ctx)
	s.watchCancel = cancel

	// Handle close
	closeWatch := func() {
		cancel()
		s.watchCancel = nil
		s.SetEditMode(false)
		s.app.SetRoot(s.rootFlex, true)
		s.SetStatusBarText("[yellow]Watch stopped for " + kv.Key)
	}

	// Create watch log view
	logView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)

	logView.SetBorder(true).
		SetTitle(" Watch: " + kv.Key + " (Press ESC to stop) ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(tcell.ColorYellow)

	// Add initial value
	_, _ = logView.Write([]byte("[cyan]Started watching key: " + kv.Key + "[-]\n\n"))
	_, _ = logView.Write([]byte("[yellow]Current value:[-]\n" + kv.Value + "\n\n"))
	_, _ = logView.Write([]byte("[gray]Waiting for changes...[-]\n"))

	// Center the watch window
	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(logView, 0, 3, true).
			AddItem(nil, 0, 1, false), 0, 2, true).
		AddItem(nil, 0, 1, false)

	// Set global InputCapture for ESC (app level, not primitive level)
	s.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			closeWatch()
			return nil
		}
		return event
	})

	s.app.SetRoot(flex, true)

	// Start watch in background
	go func() {
		cli := s.connManager.GetClient()
		if cli == nil {
			s.app.QueueUpdateDraw(func() {
				_, _ = logView.Write([]byte("[red]Error: Not connected to etcd[-]\n"))
			})
			return
		}

		err := cli.Watch(watchCtx, kv.Key, func(event *client.WatchEvent) {
			s.app.QueueUpdateDraw(func() {
				revision := event.ModRevision
				switch event.Type {
				case client.EventTypePut:
					_, _ = logView.Write([]byte(fmt.Sprintf("\n[green]► PUT[-] [gray](rev %d)[-]\n", revision)))
					_, _ = logView.Write([]byte("[yellow]New value:[-]\n" + event.Value + "\n"))
				case client.EventTypeDelete:
					_, _ = logView.Write([]byte(fmt.Sprintf("\n[red]► DELETE[-] [gray](rev %d)[-]\n", revision)))
					_, _ = logView.Write([]byte("[gray]Key was deleted[-]\n"))
				}
				logView.ScrollToEnd()
			})
		})

		if err != nil && watchCtx.Err() == nil {
			s.app.QueueUpdateDraw(func() {
				_, _ = logView.Write([]byte("[red]Watch error: " + err.Error() + "[-]\n"))
			})
		}
	}()
}

// HandleDetailsAction handles actions from the details panel buttons.
func (s *State) HandleDetailsAction(ctx context.Context, action details.ActionType) {
	switch action {
	case details.ActionEdit:
		// Don't restore focus - edit form handles its own focus
		s.HandleEdit(ctx)
	case details.ActionDelete:
		// Don't restore focus - delete modal handles its own focus
		s.HandleDelete(ctx)
	case details.ActionWatch:
		s.HandleWatch(ctx)
		s.app.SetFocus(s.keysPanel.GetTree())
	}
}

// HandleSearch shows search input for prefix search.
func (s *State) HandleSearch(ctx context.Context) {
	s.debugPanel.LogInfo("Opening search form")

	// Enable edit mode to bypass global input capture
	s.SetEditMode(true)

	// Helper to close form and restore main view
	closeForm := func() {
		s.SetEditMode(false)
		s.app.SetRoot(s.rootFlex, true)
	}

	// Create form for search
	form := tview.NewForm()

	// Add prefix input field
	form.AddInputField("Prefix", "/", 50, nil, nil)

	form.AddButton("Search", func() {
		prefix := form.GetFormItemByLabel("Prefix").(*tview.InputField).GetText()
		s.debugPanel.LogDebug("Search button clicked - Prefix: %s", prefix)

		if err := s.SearchByPrefix(ctx, prefix); err != nil {
			s.SetStatusBarText("[red]Search failed:[white] " + err.Error())
			s.debugPanel.LogError("Search failed: %v", err)
		} else {
			s.SetStatusBarText("[green]Search:[white] " + prefix)
			s.debugPanel.LogInfo("Search completed for prefix: %s", prefix)
		}

		closeForm()
	})

	form.AddButton("Clear", func() {
		// Reload all keys
		if err := s.RefreshKeys(ctx); err != nil {
			s.SetStatusBarText("[red]Failed to reload keys:[white] " + err.Error())
		} else {
			s.SetStatusBarText("[green]All keys loaded")
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

	form.SetBorder(true).SetTitle(" Search by Prefix (Tab to navigate, ESC cancel) ").SetTitleAlign(tview.AlignLeft)
	form.SetCancelFunc(closeForm)

	// Set root and focus on first form field
	s.app.SetRoot(form, true)
	form.SetFocus(0)
}

// HandleCreateNewKey create new key.
func (s *State) HandleCreateNewKey(ctx context.Context) {
	s.debugPanel.LogInfo("Opening form for new key")

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
	form.AddInputField("Key", "", 50, nil, nil)

	// Add Value text area
	form.AddTextArea("Value", "", 50, 5, 0, nil)

	form.AddButton("Save", func() {
		newKey := form.GetFormItemByLabel("Key").(*tview.InputField).GetText()
		newValue := form.GetFormItemByLabel("Value").(*tview.TextArea).GetText()

		s.debugPanel.LogDebug("Save button clicked - Key: %s, Value length: %d", newKey, len(newValue))

		if err := s.PutKey(ctx, newKey, newValue); err != nil {
			s.SetStatusBarText("[red]Failed to save:[white] " + err.Error())
			s.debugPanel.LogError("Failed to save key '%s': %v", newKey, err)
			closeForm()
			return
		}

		s.debugPanel.LogInfo("Successfully saved new key: %s", newKey)

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

	form.SetBorder(true).SetTitle(" Create Key (Tab to navigate, ESC cancel) ").SetTitleAlign(tview.AlignLeft)
	form.SetCancelFunc(closeForm)

	// Set root and focus on first form field
	s.app.SetRoot(form, true)
	form.SetFocus(0)
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

// ShowHelp displays the help window.
func (s *State) ShowHelp() {
	// Enable edit mode to bypass global input capture
	s.SetEditMode(true)

	helpText := `[yellow::b]etcdtui - Interactive TUI for etcd[-:-:-]

[cyan::b]Navigation[-:-:-]
  [green]↑/↓[-]         Navigate tree
  [green]Enter[-]       Expand/collapse node
  [green]Tab[-]         Switch panels (Keys ↔ Details)
  [green]←/→[-]         Navigate buttons

[cyan::b]Keys[-:-:-]
  [green]e[-]           Edit key/value
  [green]d[-]           Delete key
  [green]n[-]           New key
  [green]r[-]           Refresh keys
  [green]w[-]           Watch mode
  [green]/[-]           Search by prefix

[cyan::b]Other[-:-:-]
  [green]p[-]           Switch profile
  [green]F1[-]          Toggle debug panel
  [green]?[-]           Show this help
  [green]q[-]           Quit

[gray::d]Press ESC or Enter to close[-:-:-]`

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(helpText).
		SetTextAlign(tview.AlignLeft)

	textView.SetBorder(true).
		SetTitle(" Help ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(tcell.ColorYellow)

	// Handle ESC and Enter to close
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyEnter {
			s.SetEditMode(false)
			s.app.SetRoot(s.rootFlex, true)
			return nil
		}
		return event
	})

	// Center the help window
	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(textView, 22, 1, true).
			AddItem(nil, 0, 1, false), 50, 1, true).
		AddItem(nil, 0, 1, false)

	s.app.SetRoot(flex, true)
}
