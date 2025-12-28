package details

import (
	"github.com/alexandr/etcdtui/internal/ui/components/textview"
	"github.com/rivo/tview"
)

// Panel represents the details panel (right side)
type Panel struct {
	view *textview.TextView
}

// New creates a new details panel
func New() *Panel {
	return &Panel{
		view: textview.New(func(tv *tview.TextView) {
			tv.SetDynamicColors(true).
				SetText("[yellow]Select a key to view details[white]\n\nNavigation:\n[green]↓/↑[white] or [green]j/k[white] - Navigate\n[green]Enter[white] - Expand/collapse\n[green]?[white] - Help\n[green]q[white] - Quit")
			tv.SetBorder(true).SetTitle(" Details ")
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

// SetText updates the panel text
func (p *Panel) SetText(text string) {
	p.view.SetText(text)
}
