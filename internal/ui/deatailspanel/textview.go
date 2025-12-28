package deatailspanel

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
	details := d.textView.
		SetDynamicColors(true).
		SetText("[yellow]Select a key to view details[white]\n\nNavigation:\n[green]↓/↑[white] or [green]j/k[white] - Navigate\n[green]Enter[white] - Expand/collapse\n[green]?[white] - Help\n[green]q[white] - Quit")
	details.SetBorder(true).SetTitle(" TextView ")
}

func (d *TextView) GetTextView() *tview.TextView {
	return d.textView
}
