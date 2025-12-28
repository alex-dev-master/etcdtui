package general

import (
	"log"

	generalActions "github.com/alexandr/etcdtui/internal/app/actions/general"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type General struct {
	generalActions *generalActions.General
	app            *tview.Application
	rootFlex       *tview.Flex
}

func NewGeneral(app *tview.Application) *General {
	return &General{generalActions: generalActions.NewGeneral(), app: app}
}

func (g *General) Render() {
	g.generalActions.Exec()
	// Layout: дерево слева, детали справа
	mainFlex := tview.NewFlex().
		AddItem(g.generalActions.GetKeysPanelTree().GetTree(), 0, 1, true).
		AddItem(g.generalActions.GetDetailsPanel().GetTextView(), 0, 2, false)

	// Общий layout с статус баром внизу
	g.rootFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 1, true).
		AddItem(g.generalActions.GetStatusBarPanel().GetTextView(), 1, 0, false)

	g.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return g.GetInputCapture(event)
	})

	if err := g.app.SetRoot(g.rootFlex, true).EnableMouse(true).Run(); err != nil {
		log.Fatal(err)
	}
}

func (g *General) GetRootFlex() *tview.Flex {
	return g.rootFlex
}

func (g *General) GetInputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlC:
		g.app.Stop()
		return nil
	}

	switch event.Rune() {
	case 'q':
		g.app.Stop()
		return nil
	case '?':
		g.showHelp(g.app, g.rootFlex)
		return nil
	case '/':
		// TODO: Implement search
		g.generalActions.SetStatusBarText("[yellow]Search:[white] [not implemented yet]")
		return nil
	case 'n':
		// TODO: Implement new key
		g.generalActions.SetStatusBarText("[yellow]New key:[white] [not implemented yet]")
		return nil
	case 'r':
		g.generalActions.SetStatusBarText("[green]Refreshed[white] at " + "now")
		return nil
	}

	return event
}

func (g *General) showHelp(app *tview.Application, rootView tview.Primitive) {
	modal := tview.NewModal().
		SetText(`etcdtui - Interactive TUI for etcd

Keyboard Shortcuts:
  ↓/↑ or j/k  - Navigate tree
  Enter       - Expand/collapse or edit
  n           - New key
  d           - Delete key
  e           - Edit value
  w           - Watch mode
  l           - Locks dashboard
  r           - Refresh
  /           - Search
  ?           - Show this help
  q or Ctrl+C - Quit

Navigation:
  Tab         - Switch between panels
  Mouse       - Click and scroll support`).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(rootView, true)
		})

	app.SetRoot(modal, true)
}
