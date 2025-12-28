package general

import (
	"github.com/alexandr/etcdtui/internal/ui/deatailspanel"
	"github.com/alexandr/etcdtui/internal/ui/keyspanel"
	"github.com/alexandr/etcdtui/internal/ui/statusbarpanel"
	"github.com/rivo/tview"
)

type General struct {
	keysPanelTree  *keyspanel.TreeNode
	detailsPanel   *deatailspanel.TextView
	statusBarPanel *statusbarpanel.TextView
}

func NewGeneral() *General {
	return &General{
		keyspanel.NewTreeNode(),
		deatailspanel.NewTextView(),
		statusbarpanel.NewTextView(),
	}
}

func (g *General) Exec() {
	g.keysPanelTree.Draw()
	g.detailsPanel.Draw()
	g.statusBarPanel.Draw()

	// Обработка выбора узла в дереве
	g.keysPanelTree.GetTree().SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			// Раскрываем/скрываем узел
			node.SetExpanded(!node.IsExpanded())
			return
		}

		// Показываем детали ключа
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

		g.detailsPanel.GetTextView().SetText(detailsText)
	})
}

func (g *General) GetKeysPanelTree() (keysPanelTree *keyspanel.TreeNode) {
	return g.keysPanelTree
}

func (g *General) GetDetailsPanel() (detailsPanel *deatailspanel.TextView) {
	return g.detailsPanel
}

func (g *General) GetStatusBarPanel() (statusBarPanel *statusbarpanel.TextView) {
	return g.statusBarPanel
}

func (g *General) SetStatusBarText(text string) {
	g.statusBarPanel.GetTextView().SetText(text)
}
