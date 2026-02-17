package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/gbrindisi/littlefactory/internal/tasks"
)

// renderTasksPanel renders the left panel containing the task list.
// It displays each task with a status icon and highlights the active task.
func renderTasksPanel(taskList []tasks.Task, activeTaskID string, width, height int) string {
	if len(taskList) == 0 {
		return taskPanelStyle.
			Width(width).
			Height(height).
			Render("No tasks")
	}

	var lines []string
	for _, task := range taskList {
		// Get status icon
		icon := statusIcon(task.Status)

		// Determine if this is the active task
		style := normalTaskStyle
		if task.ID == activeTaskID {
			icon = "[>]"
			style = activeTaskStyle
		}

		// Truncate title to fit panel width (minus icon and padding)
		maxTitleLen := width - 6 // "[x] " = 4 chars + some padding
		title := task.Title
		if len(title) > maxTitleLen {
			title = title[:maxTitleLen-3] + "..."
		}

		line := fmt.Sprintf("%s %s", icon, title)
		lines = append(lines, style.Render(line))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	// Wrap in styled box
	return taskPanelStyle.
		Width(width).
		Height(height).
		Render(content)
}
