package layouts

import (
	"context"
	"fmt"

	"github.com/alexandr/etcdtui/internal/app/layouts/general"
	"github.com/alexandr/etcdtui/internal/app/layouts/profiles"
	"github.com/alexandr/etcdtui/internal/config"
	"github.com/rivo/tview"
)

// Manager manages application layouts.
type Manager struct {
	app            *tview.Application
	configManager  *config.Manager
	profileName    string // profile to use (from CLI flag)
	generalLayout  *general.Layout
	profilesLayout *profiles.Layout
}

// NewManager creates a new layout manager.
func NewManager() *Manager {
	app := tview.NewApplication()
	configManager := config.NewManager()

	return &Manager{
		app:           app,
		configManager: configManager,
	}
}

// SetProfileName sets the profile name to use (from CLI flag)
func (m *Manager) SetProfileName(name string) {
	m.profileName = name
}

// Render renders the application.
func (m *Manager) Render(ctx context.Context) error {
	// Load config
	if err := m.configManager.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// If profile specified via CLI flag, connect directly
	if m.profileName != "" {
		profile, err := m.configManager.GetProfile(m.profileName)
		if err != nil {
			return fmt.Errorf("profile '%s' not found", m.profileName)
		}
		return m.renderGeneralLayout(ctx, profile)
	}

	// Show profiles layout for selection
	return m.renderProfilesLayout(ctx)
}

// renderProfilesLayout shows the profile selection screen.
func (m *Manager) renderProfilesLayout(ctx context.Context) error {
	m.profilesLayout = profiles.NewLayout(m.app, m.configManager)

	// Set callback for when user selects a profile
	m.profilesLayout.SetOnConnect(func(profile *config.Profile) {
		// Stop the current app run to restart with general layout
		m.app.Stop()

		// Create new app instance for general layout
		m.app = tview.NewApplication()
		if err := m.renderGeneralLayout(ctx, profile); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	})

	// Set callback for quit
	m.profilesLayout.SetOnQuit(func() {
		m.app.Stop()
	})

	m.profilesLayout.Render()

	return m.app.Run()
}

// renderGeneralLayout shows the main etcd browser.
func (m *Manager) renderGeneralLayout(ctx context.Context, profile *config.Profile) error {
	m.generalLayout = general.NewLayout(m.app)
	m.generalLayout.SetProfile(profile)
	m.generalLayout.SetConfigManager(m.configManager)

	// Set callback to switch back to profiles
	m.generalLayout.SetOnSwitchProfile(func() {
		m.app.Stop()

		// Create new app instance for profiles layout
		m.app = tview.NewApplication()
		if err := m.renderProfilesLayout(ctx); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	})

	return m.generalLayout.Render(ctx)
}

// GetConfigManager returns the config manager
func (m *Manager) GetConfigManager() *config.Manager {
	return m.configManager
}
