package statusbarpanel

import (
	"sync"

	"github.com/rivo/tview"
)

type TextView struct {
	textView *tview.TextView
	once     sync.Once
}

func NewTextView() *TextView {
	return &TextView{textView: tview.NewTextView()}
}

// Draw Панель деталей (правая панель)
func (d *TextView) Draw() {
	d.once.Do(d.draw)
}

func (d *TextView) draw() {
	// Статус бар внизу
	statusBar := d.textView.
		SetDynamicColors(true).
		SetText("[yellow]Status:[white] Not connected | [yellow]Keys:[white] 0 | [green][/] Search  [n] New  [r] Refresh  [q] Quit  [?] Help")
	statusBar.SetBorder(false)
}

func (d *TextView) GetTextView() *tview.TextView {
	return d.textView
}
