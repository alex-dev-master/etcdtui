package profiles

import (
	"github.com/alexandr/etcdtui/internal/app/actions/profiles"
	"github.com/alexandr/etcdtui/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Layout handles the profile selection screen.
type Layout struct {
	state        *profiles.State
	app          *tview.Application
	rootFlex     *tview.Flex
	inputCapture func(event *tcell.EventKey) *tcell.EventKey
}

// NewLayout creates a new profiles layout.
func NewLayout(app *tview.Application, configManager *config.Manager) *Layout {
	state := profiles.NewState(app, configManager)
	return &Layout{
		app:   app,
		state: state,
	}
}

// SetOnConnect sets the callback when a profile is selected for connection.
func (l *Layout) SetOnConnect(fn func(profile *config.Profile)) {
	l.state.SetOnConnect(fn)
}

// SetOnQuit sets the callback when user quits.
func (l *Layout) SetOnQuit(fn func()) {
	l.state.SetOnQuit(fn)
}

// Render displays the profile selection screen.
func (l *Layout) Render() {
	// Create profile list
	profileList := tview.NewList()
	profileList.SetBorder(true).
		SetTitle(" Profiles ").
		SetTitleAlign(tview.AlignLeft)
	profileList.ShowSecondaryText(true)
	l.state.SetProfileList(profileList)

	// Create details view
	detailsView := tview.NewTextView().
		SetDynamicColors(true)
	detailsView.SetBorder(true).
		SetTitle(" Details ").
		SetTitleAlign(tview.AlignLeft)
	l.state.SetDetailsView(detailsView)

	// Create status bar
	statusBar := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[green::b]Enter[-::-] Connect  [green::b]n[-::-] New  [green::b]e[-::-] Edit  [green::b]d[-::-] Delete  [green::b]q[-::-] Quit")
	l.state.SetStatusBar(statusBar)

	// Load profiles into list
	l.state.RefreshProfileList()

	// Main layout: list on left, details on right
	mainFlex := tview.NewFlex().
		AddItem(profileList, 0, 1, true).
		AddItem(detailsView, 0, 2, false)

	// Root layout with status bar
	l.rootFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 1, true).
		AddItem(statusBar, 1, 0, false)

	l.state.SetRootFlex(l.rootFlex)

	// Store input capture for restoration
	l.inputCapture = l.handleInput

	// Set up input handling
	l.app.SetInputCapture(l.inputCapture)

	l.app.SetRoot(l.rootFlex, true)
}

// restoreInputCapture restores the input capture after form/modal closes.
func (l *Layout) restoreInputCapture() {
	l.app.SetInputCapture(l.inputCapture)
}

// handleInput handles keyboard input.
func (l *Layout) handleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEnter:
		if l.state.GetSelectedProfile() != nil {
			l.state.Connect(l.state.GetSelectedProfile())
			return nil
		}
	case tcell.KeyEsc:
		l.state.Quit()
		return nil
	}

	switch event.Rune() {
	case 'q':
		l.state.Quit()
		return nil
	case 'n':
		l.state.ShowNewProfileForm(l.restoreInputCapture)
		return nil
	case 'e':
		if l.state.GetSelectedProfile() != nil {
			l.state.ShowEditProfileForm(l.state.GetSelectedProfile(), l.restoreInputCapture)
		}
		return nil
	case 'd':
		if l.state.GetSelectedProfile() != nil {
			l.state.ShowDeleteConfirmation(l.state.GetSelectedProfile(), l.restoreInputCapture)
		}
		return nil
	}

	return event
}

// GetSelectedProfile returns the currently selected profile.
func (l *Layout) GetSelectedProfile() *config.Profile {
	return l.state.GetSelectedProfile()
}

// GetState returns the state.
func (l *Layout) GetState() *profiles.State {
	return l.state
}
