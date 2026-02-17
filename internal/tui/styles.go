package tui

import "github.com/charmbracelet/lipgloss"

// Style definitions for TUI components

var (
	// Panel styles
	taskPanelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderRight(true).
			Padding(0, 1)

	// Task list item styles
	activeTaskStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")) // Green for active task

	normalTaskStyle = lipgloss.NewStyle()

	// Status bar style
	statusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("250"))
)

// statusIcon returns the appropriate icon for a given task status.
func statusIcon(status string) string {
	switch status {
	case "closed", "done":
		return "[x]"
	case "blocked":
		return "[!]"
	case "in_progress", "active":
		return "[>]"
	default:
		return "[ ]"
	}
}
