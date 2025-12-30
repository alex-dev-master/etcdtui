package profiles

import (
	"github.com/alex-dev-master/etcdtui/internal/config"
	"github.com/rivo/tview"
)

// State holds all state for the profiles layout.
type State struct {
	// Config
	configManager *config.Manager

	// UI components
	profileList *tview.List
	detailsView *tview.TextView
	statusBar   *tview.TextView

	// Current state
	selectedProfile *config.Profile

	// App reference
	app      *tview.Application
	rootFlex *tview.Flex

	// Callbacks
	onConnect func(profile *config.Profile)
	onQuit    func()
}

// NewState creates a new profiles state.
func NewState(app *tview.Application, configManager *config.Manager) *State {
	return &State{
		app:           app,
		configManager: configManager,
	}
}

// SetOnConnect sets the callback when a profile is selected for connection.
func (s *State) SetOnConnect(fn func(profile *config.Profile)) {
	s.onConnect = fn
}

// SetOnQuit sets the callback when user quits.
func (s *State) SetOnQuit(fn func()) {
	s.onQuit = fn
}

// SetRootFlex sets the root flex for returning from modals.
func (s *State) SetRootFlex(flex *tview.Flex) {
	s.rootFlex = flex
}

// GetRootFlex returns the root flex.
func (s *State) GetRootFlex() *tview.Flex {
	return s.rootFlex
}

// SetProfileList sets the profile list component.
func (s *State) SetProfileList(list *tview.List) {
	s.profileList = list
}

// GetProfileList returns the profile list.
func (s *State) GetProfileList() *tview.List {
	return s.profileList
}

// SetDetailsView sets the details view component.
func (s *State) SetDetailsView(view *tview.TextView) {
	s.detailsView = view
}

// GetDetailsView returns the details view.
func (s *State) GetDetailsView() *tview.TextView {
	return s.detailsView
}

// SetStatusBar sets the status bar component.
func (s *State) SetStatusBar(bar *tview.TextView) {
	s.statusBar = bar
}

// GetStatusBar returns the status bar.
func (s *State) GetStatusBar() *tview.TextView {
	return s.statusBar
}

// SetSelectedProfile sets the currently selected profile.
func (s *State) SetSelectedProfile(p *config.Profile) {
	s.selectedProfile = p
}

// GetSelectedProfile returns the currently selected profile.
func (s *State) GetSelectedProfile() *config.Profile {
	return s.selectedProfile
}

// GetConfigManager returns the config manager.
func (s *State) GetConfigManager() *config.Manager {
	return s.configManager
}

// GetApp returns the tview application.
func (s *State) GetApp() *tview.Application {
	return s.app
}

// SetStatusText sets the status bar text.
func (s *State) SetStatusText(text string) {
	if s.statusBar != nil {
		s.statusBar.SetText(text)
	}
}

// Connect triggers the connection callback with the selected profile.
func (s *State) Connect(p *config.Profile) {
	if s.onConnect != nil {
		s.onConnect(p)
	}
}

// Quit triggers the quit callback.
func (s *State) Quit() {
	if s.onQuit != nil {
		s.onQuit()
	}
}
