package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/littlefactory/internal/driver"
	"github.com/yourusername/littlefactory/internal/tasks"
)

// renderStatusBar renders the bottom status bar showing run progress and keyboard hints.
func renderStatusBar(
	taskList []tasks.Task,
	iteration int,
	maxIterations int,
	autoFollow bool,
	runComplete bool,
	finalStatus driver.RunStatus,
	width int,
) string {
	// Count task statuses
	var done, pending, blocked int
	for _, task := range taskList {
		switch task.Status {
		case "closed", "done":
			done++
		case "blocked":
			blocked++
		default:
			pending++
		}
	}

	// Build status text
	var statusText string
	if runComplete {
		statusText = fmt.Sprintf("Run complete: %s", finalStatus)
	} else {
		statusText = fmt.Sprintf("Iteration %d/%d | Tasks: %d done, %d pending, %d blocked",
			iteration, maxIterations, done, pending, blocked)
	}

	// Add keyboard hints
	followText := "off"
	if autoFollow {
		followText = "on"
	}
	hints := fmt.Sprintf("q:quit | f:follow(%s) | up/dn:scroll", followText)

	// Combine status and hints
	leftText := statusText
	rightText := hints

	// Calculate spacing
	spacing := width - lipgloss.Width(leftText) - lipgloss.Width(rightText)
	if spacing < 0 {
		spacing = 0
	}

	statusLine := leftText + lipgloss.NewStyle().Width(spacing).Render("") + rightText

	return statusBarStyle.
		Width(width).
		Render(statusLine)
}
