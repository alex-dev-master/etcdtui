package general

import (
	"context"
	"log"

	"github.com/alex-dev-master/etcdtui/internal/app/actions/general"
	"github.com/alex-dev-master/etcdtui/internal/config"
	"github.com/alex-dev-master/etcdtui/internal/ui/panels/details"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Layout handles the visual layout and input routing for the general view.
type Layout struct {
	state           *general.State
	app             *tview.Application
	rootFlex        *tview.Flex
	mainFlex        *tview.Flex
	contentFlex     *tview.Flex
	profile         *config.Profile
	configManager   *config.Manager
	onSwitchProfile func()
}

// NewLayout creates a new Layout with the given tview application.
func NewLayout(app *tview.Application) *Layout {
	state := general.NewState()
	state.SetApp(app)
	return &Layout{
		state: state,
		app:   app,
	}
}

// SetProfile sets the profile to use for connection
func (l *Layout) SetProfile(profile *config.Profile) {
	l.profile = profile
}

// SetConfigManager sets the config manager
func (l *Layout) SetConfigManager(cm *config.Manager) {
	l.configManager = cm
}

// SetOnSwitchProfile sets the callback for switching profiles
func (l *Layout) SetOnSwitchProfile(fn func()) {
	l.onSwitchProfile = fn
}

// Render initializes and displays the layout.
func (l *Layout) Render(ctx context.Context) error {
	// Pass profile to state for connection
	if l.profile != nil {
		l.state.SetProfile(l.profile)
	}
	if l.configManager != nil {
		l.state.SetConfigManager(l.configManager)
	}

	if err := l.state.InitConnection(ctx); err != nil {
		return err
	}

	// Setup action callback for details panel
	l.state.GetDetailsPanel().SetActionCallback(func(action details.ActionType) {
		l.state.HandleDetailsAction(ctx, action)
	})

	// Setup tab callback for details panel to switch focus back to keys
	l.state.GetDetailsPanel().SetTabCallback(func() {
		l.app.SetFocus(l.state.GetKeysPanel().GetTree())
	})

	// Layout: tree on left, details on right
	l.mainFlex = tview.NewFlex().
		AddItem(l.state.GetKeysPanel().GetTree(), 0, 1, true).
		AddItem(l.state.GetDetailsPanel().GetView(), 0, 2, false)

	// Content flex: main view + optional debug panel (side by side)
	l.contentFlex = tview.NewFlex().
		AddItem(l.mainFlex, 0, 1, true)

	// Main layout with status bar at bottom
	l.rootFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(l.contentFlex, 0, 1, true).
		AddItem(l.state.GetStatusBarPanel().GetView(), 1, 0, false)

	// Set root flex in state for modals to return
	l.state.SetRootFlex(l.rootFlex)

	// Set input capture through state so it can be disabled/restored for modals
	l.state.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return l.handleInput(ctx, event)
	})

	if err := l.app.SetRoot(l.rootFlex, true).EnableMouse(false).Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}

// handleInput routes keyboard input to appropriate handlers.
func (l *Layout) handleInput(ctx context.Context, event *tcell.EventKey) *tcell.EventKey {
	// When in edit mode, only handle Ctrl+C, pass everything else through
	if l.state.IsEditMode() {
		if event.Key() == tcell.KeyCtrlC {
			l.app.Stop()
			return nil
		}
		return event
	}

	// Don't intercept keys when we're not in the main view
	if l.app.GetFocus() != l.state.GetKeysPanel().GetTree() &&
		l.app.GetFocus() != l.state.GetDetailsPanel().GetForm() {
		if event.Key() == tcell.KeyCtrlC {
			l.app.Stop()
			return nil
		}
		return event
	}

	// Handle special keys
	switch event.Key() {
	case tcell.KeyCtrlC:
		l.app.Stop()
		return nil
	case tcell.KeyF1:
		l.state.ToggleDebugPanel(l.contentFlex)
		return nil
	case tcell.KeyTab:
		return l.handleTab()
	}

	// Handle rune keys
	switch event.Rune() {
	case 'q':
		l.app.Stop()
		return nil
	case '?':
		l.state.ShowHelp()
		return nil
	case '/':
		l.state.HandleSearch(ctx)
		return nil
	case 'n':
		l.state.HandleCreateNewKey(ctx)
		return nil
	case 'r':
		if err := l.state.RefreshKeys(ctx); err != nil {
			l.state.SetStatusBarText("[red]Failed to refresh:[white] " + err.Error())
		}
		return nil
	case 'd':
		l.state.HandleDelete(ctx)
		return nil
	case 'e':
		l.state.HandleEdit(ctx)
		return nil
	case 'w':
		l.state.HandleWatch(ctx)
		return nil
	case 'p':
		// Switch profiles
		if l.onSwitchProfile != nil {
			l.onSwitchProfile()
		}
		return nil
	}

	return event
}

// handleTab switches focus between panels.
func (l *Layout) handleTab() *tcell.EventKey {
	current := l.app.GetFocus()
	if current == l.state.GetKeysPanel().GetTree() && l.state.GetCurrentKey() != nil {
		l.app.SetFocus(l.state.GetDetailsPanel().GetForm())
		return nil
	}
	return nil
}

// GetRootFlex returns the root flex layout.
func (l *Layout) GetRootFlex() *tview.Flex {
	return l.rootFlex
}

// GetState returns the state.
func (l *Layout) GetState() *general.State {
	return l.state
}
