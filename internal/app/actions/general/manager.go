package general

import (
	"github.com/alexandr/etcdtui/internal/ui/panels/details"
	"github.com/alexandr/etcdtui/internal/ui/panels/keys"
	"github.com/alexandr/etcdtui/internal/ui/panels/statusbar"
	"github.com/rivo/tview"
)

type General struct {
	keysPanel      *keys.Panel
	detailsPanel   *details.Panel
	statusBarPanel *statusbar.Panel
}

func NewGeneral() *General {
	return &General{
		keysPanel:      keys.New(),
		detailsPanel:   details.New(),
		statusBarPanel: statusbar.New(),
	}
}

func (g *General) Exec() {
	g.keysPanel.Draw()
	g.detailsPanel.Draw()
	g.statusBarPanel.Draw()

	g.keysPanel.GetTree().SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			// Toggle node expansion
			node.SetExpanded(!node.IsExpanded())
			return
		}

		// Show key details
		key := reference.(string)
		var value string
		switch key {
		case "/services/api/v1/config":
			value = `{
  "port": 8080,
  "timeout": 30,
  "debug": false
}`
		case "/services/api/v1/endpoints":
			value = `["http://api1.example.com", "http://api2.example.com"]`
		case "/services/auth/jwt-secret":
			value = "********** [hidden]"
		case "/config/database-url":
			value = "postgres://localhost:5432/mydb"
		case "/config/redis-url":
			value = "redis://localhost:6379"
		case "/locks/payment":
			value = "Locked by: pod-123\nAcquired: 2m ago"
		default:
			value = "No data"
		}

		detailsText := "[yellow]Key:[white] " + key + "\n\n"
		detailsText += "[yellow]Value:[white]\n" + value + "\n\n"
		detailsText += "[yellow]Revision:[white] 12345\n"
		detailsText += "[yellow]Modified:[white] 2024-12-28 10:15:23\n"
		detailsText += "[yellow]TTL:[white] ∞\n\n"
		detailsText += "[green][e][white] Edit  [green][d][white] Delete  [green][w][white] Watch  [green][c][white] Copy"

		g.detailsPanel.SetText(detailsText)
	})
}

func (g *General) GetKeysPanel() *keys.Panel {
	return g.keysPanel
}

func (g *General) GetDetailsPanel() *details.Panel {
	return g.detailsPanel
}

func (g *General) GetStatusBarPanel() *statusbar.Panel {
	return g.statusBarPanel
}

func (g *General) SetStatusBarText(text string) {
	g.statusBarPanel.SetText(text)
}
