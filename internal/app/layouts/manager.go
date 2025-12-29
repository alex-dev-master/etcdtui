package layouts

import (
	"context"

	"github.com/alexandr/etcdtui/internal/app/layouts/general"
	"github.com/rivo/tview"
)

type Manager struct {
	generalLayout *general.General
	app           *tview.Application
}

func NewManager() *Manager {
	app := tview.NewApplication()
	return &Manager{
		generalLayout: general.NewGeneral(app),
		app:           app,
	}
}

func (m *Manager) Render(ctx context.Context) (err error) {
	return m.generalLayoutRender(ctx)
}

func (m *Manager) generalLayoutRender(ctx context.Context) (err error) {
	return m.generalLayout.Render(ctx)
}
