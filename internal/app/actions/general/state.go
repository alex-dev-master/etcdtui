package general

import (
	"context"

	"github.com/alexandr/etcdtui/internal/app/connection/etcd"
	"github.com/alexandr/etcdtui/internal/config"
	"github.com/alexandr/etcdtui/internal/ui/panels/debug"
	"github.com/alexandr/etcdtui/internal/ui/panels/details"
	"github.com/alexandr/etcdtui/internal/ui/panels/keys"
	"github.com/alexandr/etcdtui/internal/ui/panels/statusbar"
	client "github.com/alexandr/etcdtui/pkg/etcd"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// State holds all state for the general layout including panels,
// connection manager, and UI state flags.
type State struct {
	// UI panels
	keysPanel      *keys.Panel
	detailsPanel   *details.Panel
	statusBarPanel *statusbar.Panel
	debugPanel     *debug.Panel

	// Connection
	connManager *etcd.Manager

	// Profile and config
	profile       *config.Profile
	configManager *config.Manager

	// Current state
	currentKey  *client.KeyValue
	inEditMode  bool
	watchCancel context.CancelFunc // Cancel function for active watch

	// App reference for UI operations
	app          *tview.Application
	rootFlex     *tview.Flex
	inputCapture func(event *tcell.EventKey) *tcell.EventKey
}

// NewState creates a new State with initialized panels and connection manager.
func NewState() *State {
	return &State{
		keysPanel:      keys.New(),
		detailsPanel:   details.New(),
		statusBarPanel: statusbar.New(),
		debugPanel:     debug.New(),
		connManager:    etcd.NewManager(),
	}
}

// SetProfile sets the profile to use for connection
func (s *State) SetProfile(profile *config.Profile) {
	s.profile = profile
}

// GetProfile returns the current profile
func (s *State) GetProfile() *config.Profile {
	return s.profile
}

// SetConfigManager sets the config manager
func (s *State) SetConfigManager(cm *config.Manager) {
	s.configManager = cm
}

// GetConfigManager returns the config manager
func (s *State) GetConfigManager() *config.Manager {
	return s.configManager
}

// SetApp sets the tview application reference.
func (s *State) SetApp(app *tview.Application) {
	s.app = app
}

// GetApp returns the tview application.
func (s *State) GetApp() *tview.Application {
	return s.app
}

// SetRootFlex sets the root flex layout for returning from modals.
func (s *State) SetRootFlex(flex *tview.Flex) {
	s.rootFlex = flex
}

// GetRootFlex returns the root flex layout.
func (s *State) GetRootFlex() *tview.Flex {
	return s.rootFlex
}

// SetCurrentKey sets the currently selected key.
func (s *State) SetCurrentKey(kv *client.KeyValue) {
	s.currentKey = kv
}

// GetCurrentKey returns the currently selected key.
func (s *State) GetCurrentKey() *client.KeyValue {
	return s.currentKey
}

// SetEditMode sets the edit mode flag and manages InputCapture.
// When entering edit mode, InputCapture is disabled to allow modals to work.
// When exiting edit mode, InputCapture is restored.
func (s *State) SetEditMode(mode bool) {
	s.inEditMode = mode
	if mode {
		// Disable global input capture for modals
		s.app.SetInputCapture(nil)
	} else {
		// Restore global input capture
		s.app.SetInputCapture(s.inputCapture)
	}
}

// IsEditMode returns true if currently in edit mode.
func (s *State) IsEditMode() bool {
	return s.inEditMode
}

// SetInputCapture stores the input capture function for later restoration.
func (s *State) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) {
	s.inputCapture = capture
	s.app.SetInputCapture(capture)
}

// GetKeysPanel returns the keys panel.
func (s *State) GetKeysPanel() *keys.Panel {
	return s.keysPanel
}

// GetDetailsPanel returns the details panel.
func (s *State) GetDetailsPanel() *details.Panel {
	return s.detailsPanel
}

// GetStatusBarPanel returns the status bar panel.
func (s *State) GetStatusBarPanel() *statusbar.Panel {
	return s.statusBarPanel
}

// GetDebugPanel returns the debug panel.
func (s *State) GetDebugPanel() *debug.Panel {
	return s.debugPanel
}

// GetConnectionManager returns the connection manager.
func (s *State) GetConnectionManager() *etcd.Manager {
	return s.connManager
}

// SetStatusBarText sets the status bar text.
func (s *State) SetStatusBarText(text string) {
	s.statusBarPanel.SetText(text)
}
