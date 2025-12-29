package layouts

import (
	"context"

	"github.com/alexandr/etcdtui/internal/app/layouts/general"
	"github.com/rivo/tview"
)

// Manager manages application layouts.
type Manager struct {
	generalLayout *general.Layout
	app           *tview.Application
}

// NewManager creates a new layout manager.
func NewManager() *Manager {
	app := tview.NewApplication()
	return &Manager{
		generalLayout: general.NewLayout(app),
		app:           app,
	}
}

// Render renders the current layout.
func (m *Manager) Render(ctx context.Context) error {
	return m.generalLayout.Render(ctx)
}
