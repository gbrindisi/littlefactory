package tui

import (
	"strings"
	"testing"

	"github.com/gbrindisi/littlefactory/internal/tasks"
)

func TestRenderTasksPanel_EmptyList(t *testing.T) {
	result := renderTasksPanel([]tasks.Task{}, "", 30, 10)

	if !strings.Contains(result, "No tasks") {
		t.Errorf("Expected 'No tasks' message for empty list, got: %s", result)
	}
}

func TestRenderTasksPanel_WithTasks(t *testing.T) {
	taskList := []tasks.Task{
		{
			ID:     "task-1",
			Title:  "First task",
			Status: "pending",
		},
		{
			ID:     "task-2",
			Title:  "Second task",
			Status: "done",
		},
		{
			ID:     "task-3",
			Title:  "Blocked task",
			Status: "blocked",
		},
	}

	result := renderTasksPanel(taskList, "", 30, 10)

	// Check that all tasks are rendered
	if !strings.Contains(result, "First task") {
		t.Errorf("Expected 'First task' in output")
	}
	if !strings.Contains(result, "Second task") {
		t.Errorf("Expected 'Second task' in output")
	}
	if !strings.Contains(result, "Blocked task") {
		t.Errorf("Expected 'Blocked task' in output")
	}

	// Check that status icons are present
	if !strings.Contains(result, "[ ]") {
		t.Errorf("Expected '[ ]' icon for pending task")
	}
	if !strings.Contains(result, "[x]") {
		t.Errorf("Expected '[x]' icon for done task")
	}
	if !strings.Contains(result, "[!]") {
		t.Errorf("Expected '[!]' icon for blocked task")
	}
}

func TestRenderTasksPanel_ActiveTask(t *testing.T) {
	taskList := []tasks.Task{
		{
			ID:     "task-1",
			Title:  "First task",
			Status: "pending",
		},
		{
			ID:     "task-2",
			Title:  "Active task",
			Status: "in_progress",
		},
	}

	result := renderTasksPanel(taskList, "task-2", 30, 10)

	// Active task should have [>] icon (overrides status icon)
	if !strings.Contains(result, "[>]") {
		t.Errorf("Expected '[>]' icon for active task")
	}
	if !strings.Contains(result, "Active task") {
		t.Errorf("Expected 'Active task' in output")
	}
}

func TestRenderTasksPanel_CursorPosition(t *testing.T) {
	taskList := []tasks.Task{
		{
			ID:     "task-1",
			Title:  "First task",
			Status: "pending",
		},
		{
			ID:     "task-2",
			Title:  "Second task",
			Status: "pending",
		},
	}

	// Test without cursor (cursor removed from implementation)
	result := renderTasksPanel(taskList, "", 30, 10)
	if !strings.Contains(result, "First task") {
		t.Errorf("Expected 'First task' in output")
	}
	if !strings.Contains(result, "Second task") {
		t.Errorf("Expected 'Second task' in output")
	}
}

func TestRenderTasksPanel_TitleTruncation(t *testing.T) {
	taskList := []tasks.Task{
		{
			ID:     "task-1",
			Title:  "This is a very long task title that should be truncated to fit within the panel width",
			Status: "pending",
		},
	}

	result := renderTasksPanel(taskList, "", 30, 10)

	// The result should contain the truncated title with "..."
	if !strings.Contains(result, "...") {
		t.Errorf("Expected truncation indicator '...' for long title")
	}

	// The full title should not be present (it's too long)
	fullTitle := "This is a very long task title that should be truncated to fit within the panel width"
	if strings.Contains(result, fullTitle) {
		t.Errorf("Long title should have been truncated")
	}
}

func TestRenderTasksPanel_Dimensions(t *testing.T) {
	taskList := []tasks.Task{
		{
			ID:     "task-1",
			Title:  "Task",
			Status: "pending",
		},
	}

	// Test with different dimensions to ensure no panic
	testCases := []struct {
		width  int
		height int
	}{
		{30, 10},
		{20, 5},
		{50, 20},
		{10, 2},
	}

	for _, tc := range testCases {
		result := renderTasksPanel(taskList, "", tc.width, tc.height)
		if result == "" {
			t.Errorf("Expected non-empty result for width=%d, height=%d", tc.width, tc.height)
		}
	}
}
