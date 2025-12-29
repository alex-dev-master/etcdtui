package details

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ActionType represents the type of action to perform
type ActionType int

const (
	ActionNone ActionType = iota
	ActionEdit
	ActionDelete
	ActionWatch
)

// ActionCallback is called when a button is pressed
type ActionCallback func(action ActionType)

// TabCallback is called when Tab is pressed in the form
type TabCallback func()

// Panel represents the details panel (right side)
type Panel struct {
	flex          *tview.Flex
	textView      *tview.TextView
	form          *tview.Form
	once          sync.Once
	callback      ActionCallback
	tabCallback   TabCallback
	buttonsShown  bool
	currentButton int
	mu            sync.Mutex
}

// New creates a new details panel
func New() *Panel {
	return &Panel{
		flex:         tview.NewFlex(),
		textView:     tview.NewTextView(),
		form:         tview.NewForm(),
		buttonsShown: false,
	}
}

// Draw initializes the panel
func (p *Panel) Draw() {
	p.once.Do(p.initialize)
}

func (p *Panel) initialize() {
	// Setup TextView
	p.textView.
		SetDynamicColors(true).
		SetText("[yellow]Select a key to view details[white]\n\nKeyboard Navigation:\n[green]↓/↑[white] or [green]j/k[white] - Navigate tree\n[green]Enter[white] - Select key\n[green]Tab[white] - Switch panels (Keys ↔ Details)\n[green]←/→[white] or [green]h/l[white] - Navigate buttons\n[green]Enter[white] - Activate button\n\nQuick Actions:\n[green]e[white] Edit  [green]d[white] Delete  [green]r[white] Refresh  [green]?[white] Help  [green]q[white] Quit").
		SetScrollable(true)

	// Setup Form with buttons
	p.form.
		SetButtonsAlign(tview.AlignLeft).
		SetButtonBackgroundColor(tcell.ColorDarkGreen).
		SetButtonTextColor(tcell.ColorWhite).
		SetLabelColor(tcell.ColorYellow).
		SetFieldBackgroundColor(tcell.ColorBlack)

	p.form.AddButton("Edit [e]", func() {
		if p.callback != nil {
			p.callback(ActionEdit)
		}
	})

	p.form.AddButton("Delete [d]", func() {
		if p.callback != nil {
			p.callback(ActionDelete)
		}
	})

	p.form.AddButton("Watch [w]", func() {
		if p.callback != nil {
			p.callback(ActionWatch)
		}
	})

	// Setup input capture for button navigation
	p.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Handle Tab to switch back to Keys panel
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab {
			if p.tabCallback != nil {
				p.tabCallback()
				return nil
			}
		}

		p.mu.Lock()
		defer p.mu.Unlock()

		buttonCount := p.form.GetButtonCount()

		switch event.Key() {
		case tcell.KeyRight:
			// Move to next button
			if p.currentButton < buttonCount-1 {
				p.currentButton++
				p.form.SetFocus(p.currentButton)
			}
			return nil
		case tcell.KeyLeft:
			// Move to previous button
			if p.currentButton > 0 {
				p.currentButton--
				p.form.SetFocus(p.currentButton)
			}
			return nil
		}

		switch event.Rune() {
		case 'l', 'L':
			// Move to next button (vim-style)
			if p.currentButton < buttonCount-1 {
				p.currentButton++
				p.form.SetFocus(p.currentButton)
			}
			return nil
		case 'h', 'H':
			// Move to previous button (vim-style)
			if p.currentButton > 0 {
				p.currentButton--
				p.form.SetFocus(p.currentButton)
			}
			return nil
		}

		return event
	})

	// Setup Flex layout - initially without buttons
	p.flex.
		SetDirection(tview.FlexRow).
		AddItem(p.textView, 0, 1, false)

	p.flex.SetBorder(true).SetTitle(" Details ")
}

// GetView returns the underlying Flex
func (p *Panel) GetView() tview.Primitive {
	return p.flex
}

// SetText updates the panel text
func (p *Panel) SetText(text string) {
	p.textView.SetText(text)
}

// ShowButtons shows the action buttons
func (p *Panel) ShowButtons() {
	if !p.buttonsShown {
		p.mu.Lock()
		p.currentButton = 0
		p.mu.Unlock()
		p.flex.AddItem(p.form, 3, 0, false)
		p.form.SetFocus(0)
		p.buttonsShown = true
	}
}

// HideButtons hides the action buttons
func (p *Panel) HideButtons() {
	if p.buttonsShown {
		p.flex.RemoveItem(p.form)
		p.buttonsShown = false
	}
}

// SetActionCallback sets the callback for button actions
func (p *Panel) SetActionCallback(callback ActionCallback) {
	p.callback = callback
}

// SetTabCallback sets the callback for Tab key
func (p *Panel) SetTabCallback(callback TabCallback) {
	p.tabCallback = callback
}

// GetForm returns the form (for focus management)
func (p *Panel) GetForm() *tview.Form {
	return p.form
}
