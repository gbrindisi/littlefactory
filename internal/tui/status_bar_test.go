package tui

import (
	"strings"
	"testing"

	"github.com/yourusername/littlefactory/internal/driver"
	"github.com/yourusername/littlefactory/internal/tasks"
)

func TestRenderStatusBar_EmptyTaskList(t *testing.T) {
	result := renderStatusBar(
		[]tasks.Task{},
		1,
		10,
		true,
		false,
		driver.RunStatusCompleted,
		80,
	)

	// Should show 0 for all task counts
	if !strings.Contains(result, "0 done") {
		t.Errorf("Expected '0 done' for empty task list")
	}
	if !strings.Contains(result, "0 pending") {
		t.Errorf("Expected '0 pending' for empty task list")
	}
	if !strings.Contains(result, "0 blocked") {
		t.Errorf("Expected '0 blocked' for empty task list")
	}
}

func TestRenderStatusBar_TaskCounting(t *testing.T) {
	taskList := []tasks.Task{
		{ID: "task-1", Status: "done"},
		{ID: "task-2", Status: "closed"},
		{ID: "task-3", Status: "pending"},
		{ID: "task-4", Status: "in_progress"},
		{ID: "task-5", Status: "blocked"},
		{ID: "task-6", Status: "blocked"},
	}

	result := renderStatusBar(
		taskList,
		3,
		10,
		true,
		false,
		driver.RunStatusCompleted,
		80,
	)

	// Should count: 2 done (done + closed), 2 pending (pending + in_progress), 2 blocked
	if !strings.Contains(result, "2 done") {
		t.Errorf("Expected '2 done' (done + closed), got: %s", result)
	}
	if !strings.Contains(result, "2 pending") {
		t.Errorf("Expected '2 pending' (pending + in_progress), got: %s", result)
	}
	if !strings.Contains(result, "2 blocked") {
		t.Errorf("Expected '2 blocked', got: %s", result)
	}
}

func TestRenderStatusBar_IterationDisplay(t *testing.T) {
	result := renderStatusBar(
		[]tasks.Task{},
		5,
		10,
		true,
		false,
		driver.RunStatusCompleted,
		80,
	)

	// Should show current iteration
	if !strings.Contains(result, "Iteration 5/10") {
		t.Errorf("Expected 'Iteration 5/10' in status bar, got: %s", result)
	}
}

func TestRenderStatusBar_AutoFollowOn(t *testing.T) {
	result := renderStatusBar(
		[]tasks.Task{},
		1,
		10,
		true, // autoFollow = true
		false,
		driver.RunStatusCompleted,
		80,
	)

	// Should indicate auto-follow is on
	if !strings.Contains(result, "follow(on)") {
		t.Errorf("Expected 'follow(on)' when autoFollow is true, got: %s", result)
	}
}

func TestRenderStatusBar_AutoFollowOff(t *testing.T) {
	result := renderStatusBar(
		[]tasks.Task{},
		1,
		10,
		false, // autoFollow = false
		false,
		driver.RunStatusCompleted,
		80,
	)

	// Should indicate auto-follow is off
	if !strings.Contains(result, "follow(off)") {
		t.Errorf("Expected 'follow(off)' when autoFollow is false, got: %s", result)
	}
}

func TestRenderStatusBar_RunComplete(t *testing.T) {
	result := renderStatusBar(
		[]tasks.Task{},
		10,
		10,
		true,
		true, // runComplete = true
		driver.RunStatusCompleted,
		80,
	)

	// Should show run complete status
	if !strings.Contains(result, "Run complete") {
		t.Errorf("Expected 'Run complete' when run is complete, got: %s", result)
	}
	if !strings.Contains(result, string(driver.RunStatusCompleted)) {
		t.Errorf("Expected final status to be shown, got: %s", result)
	}
}

func TestRenderStatusBar_RunNotComplete(t *testing.T) {
	result := renderStatusBar(
		[]tasks.Task{},
		3,
		10,
		true,
		false, // runComplete = false
		driver.RunStatusCompleted,
		80,
	)

	// Should show iteration info instead of complete message
	if strings.Contains(result, "Run complete") {
		t.Errorf("Should not show 'Run complete' when run is not complete")
	}
	if !strings.Contains(result, "Iteration") {
		t.Errorf("Expected 'Iteration' in status bar when run is not complete, got: %s", result)
	}
}

func TestRenderStatusBar_KeyboardHints(t *testing.T) {
	result := renderStatusBar(
		[]tasks.Task{},
		1,
		10,
		true,
		false,
		driver.RunStatusCompleted,
		80,
	)

	// Should show keyboard hints
	if !strings.Contains(result, "q:quit") {
		t.Errorf("Expected 'q:quit' hint in status bar")
	}
	if !strings.Contains(result, "f:follow") {
		t.Errorf("Expected 'f:follow' hint in status bar")
	}
	if !strings.Contains(result, "scroll") {
		t.Errorf("Expected scroll hint in status bar")
	}
}

func TestRenderStatusBar_VariousStatuses(t *testing.T) {
	statuses := []driver.RunStatus{
		driver.RunStatusCompleted,
		driver.RunStatusFailed,
		driver.RunStatusCancelled,
		driver.RunStatusRunning,
	}

	for _, status := range statuses {
		result := renderStatusBar(
			[]tasks.Task{},
			10,
			10,
			true,
			true,
			status,
			80,
		)

		// Each status should be displayed
		if !strings.Contains(result, string(status)) {
			t.Errorf("Expected status '%s' to be displayed in: %s", status, result)
		}
	}
}

func TestRenderStatusBar_Width(t *testing.T) {
	// Test with various widths to ensure no panic
	widths := []int{40, 80, 120, 160}

	for _, width := range widths {
		result := renderStatusBar(
			[]tasks.Task{{ID: "task-1", Status: "pending"}},
			1,
			10,
			true,
			false,
			driver.RunStatusCompleted,
			width,
		)

		if result == "" {
			t.Errorf("Expected non-empty result for width=%d", width)
		}
	}
}
