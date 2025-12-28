package statusbar

import (
	"github.com/alexandr/etcdtui/internal/ui/components/textview"
	"github.com/rivo/tview"
)

// Panel represents the status bar panel (bottom)
type Panel struct {
	view *textview.TextView
}

// New creates a new status bar panel
func New() *Panel {
	return &Panel{
		view: textview.New(func(tv *tview.TextView) {
			tv.SetDynamicColors(true).
				SetText("[yellow]Status:[white] Not connected | [yellow]Keys:[white] 0 | [green][/] Search  [n] New  [r] Refresh  [q] Quit  [?] Help")
			tv.SetBorder(false)
		}),
	}
}

// Draw initializes the panel
func (p *Panel) Draw() {
	p.view.Draw()
}

// GetView returns the underlying TextView
func (p *Panel) GetView() *tview.TextView {
	return p.view.Get()
}

// SetText updates the status bar text
func (p *Panel) SetText(text string) {
	p.view.SetText(text)
}
