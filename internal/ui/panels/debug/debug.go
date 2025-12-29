package debug

import (
	"fmt"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Panel represents the debug log panel
type Panel struct {
	textView *tview.TextView
	visible  bool
	mu       sync.Mutex
}

// New creates a new debug panel
func New() *Panel {
	p := &Panel{
		textView: tview.NewTextView(),
		visible:  false,
	}

	p.textView.
		SetDynamicColors(true).
		SetScrollable(true)

	p.textView.SetBorder(true).
		SetTitle(" Debug Logs (F1 to toggle) ").
		SetBorderColor(tcell.ColorYellow)

	return p
}

// GetView returns the underlying TextView
func (p *Panel) GetView() *tview.TextView {
	return p.textView
}

// IsVisible returns whether the panel is visible
func (p *Panel) IsVisible() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.visible
}

// SetVisible sets the visibility state
func (p *Panel) SetVisible(visible bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.visible = visible
}

// Log writes a log message to the debug panel
func (p *Panel) Log(format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[gray]%s[white] %s\n", timestamp, message)

	fmt.Fprintf(p.textView, logLine)
}

// LogInfo writes an info message
func (p *Panel) LogInfo(format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[gray]%s[white] [cyan][INFO][white] %s\n", timestamp, message)

	fmt.Fprintf(p.textView, logLine)
}

// LogError writes an error message
func (p *Panel) LogError(format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[gray]%s[white] [red][ERROR][white] %s\n", timestamp, message)

	fmt.Fprintf(p.textView, logLine)
}

// LogWarn writes a warning message
func (p *Panel) LogWarn(format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[gray]%s[white] [yellow][WARN][white] %s\n", timestamp, message)

	fmt.Fprintf(p.textView, logLine)
}

// LogDebug writes a debug message
func (p *Panel) LogDebug(format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[gray]%s[white] [green][DEBUG][white] %s\n", timestamp, message)

	fmt.Fprintf(p.textView, logLine)
}

// Clear clears all logs
func (p *Panel) Clear() {
	p.textView.Clear()
}
