package textview

import (
	"sync"

	"github.com/rivo/tview"
)

// TextView is a wrapper around tview.TextView with lazy initialization
type TextView struct {
	textView *tview.TextView
	once     sync.Once
	initFunc func(*tview.TextView)
}

// New creates a new TextView with an optional initialization function
func New(initFunc func(*tview.TextView)) *TextView {
	return &TextView{
		textView: tview.NewTextView(),
		initFunc: initFunc,
	}
}

// Draw performs lazy initialization of the TextView
func (t *TextView) Draw() {
	t.once.Do(func() {
		if t.initFunc != nil {
			t.initFunc(t.textView)
		}
	})
}

// Get returns the underlying tview.TextView
func (t *TextView) Get() *tview.TextView {
	return t.textView
}

// SetText sets the text content
func (t *TextView) SetText(text string) *TextView {
	t.textView.SetText(text)
	return t
}

// SetDynamicColors enables dynamic color parsing
func (t *TextView) SetDynamicColors(enable bool) *TextView {
	t.textView.SetDynamicColors(enable)
	return t
}

// SetBorder sets the border visibility
func (t *TextView) SetBorder(show bool) *TextView {
	t.textView.SetBorder(show)
	return t
}

// SetTitle sets the title of the border
func (t *TextView) SetTitle(title string) *TextView {
	t.textView.SetTitle(title)
	return t
}
