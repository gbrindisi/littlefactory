package tui

import "github.com/charmbracelet/bubbles/viewport"

// OutputPanel wraps the viewport component for displaying agent output.
// It provides methods for content management and rendering.
type OutputPanel struct {
	viewport viewport.Model
}

// NewOutputPanel creates a new output panel with the given dimensions.
func NewOutputPanel(width, height int) *OutputPanel {
	vp := viewport.New(width, height)
	return &OutputPanel{
		viewport: vp,
	}
}

// SetContent updates the viewport content.
func (op *OutputPanel) SetContent(content string) {
	op.viewport.SetContent(content)
}

// GotoBottom scrolls the viewport to the bottom.
func (op *OutputPanel) GotoBottom() {
	op.viewport.GotoBottom()
}

// LineUp scrolls up by n lines.
func (op *OutputPanel) LineUp(n int) {
	op.viewport.LineUp(n)
}

// LineDown scrolls down by n lines.
func (op *OutputPanel) LineDown(n int) {
	op.viewport.LineDown(n)
}

// ViewUp scrolls up by one page.
func (op *OutputPanel) ViewUp() {
	op.viewport.ViewUp()
}

// ViewDown scrolls down by one page.
func (op *OutputPanel) ViewDown() {
	op.viewport.ViewDown()
}

// SetSize updates the viewport dimensions.
func (op *OutputPanel) SetSize(width, height int) {
	op.viewport.Width = width
	op.viewport.Height = height
}

// View renders the viewport content.
func (op *OutputPanel) View() string {
	return op.viewport.View()
}
