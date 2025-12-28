package layouts

import (
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

func (m *Manager) Render() {
	m.generalLayoutRender()
}

func (m *Manager) generalLayoutRender() {
	m.generalLayout.Render()
}
