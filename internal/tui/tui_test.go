package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gbrindisi/littlefactory/internal/config"
	"github.com/gbrindisi/littlefactory/internal/driver"
	"github.com/gbrindisi/littlefactory/internal/tasks"
)

// testConfig returns a test config for use in tests
func testConfig() *config.Config {
	return &config.Config{
		StateDir: ".littlefactory",
	}
}

func TestNew(t *testing.T) {
	eventChan := make(chan interface{})
	defer close(eventChan)

	model := New(eventChan, testConfig(), "/test/project")

	if model == nil {
		t.Fatal("Expected non-nil model")
	}
	if model.tasks == nil {
		t.Error("Expected tasks to be initialized")
	}
	if len(model.tasks) != 0 {
		t.Errorf("Expected empty tasks initially, got %d", len(model.tasks))
	}
	if model.outputPanel == nil {
		t.Error("Expected outputPanel to be initialized")
	}
	if !model.autoFollow {
		t.Error("Expected autoFollow to be true initially")
	}
	if model.eventChan == nil {
		t.Error("Expected eventChan to be set")
	}
}

func TestModel_Init(t *testing.T) {
	eventChan := make(chan interface{})
	defer close(eventChan)

	model := New(eventChan, testConfig(), "/test/project")
	cmd := model.Init()

	if cmd == nil {
		t.Error("Expected Init to return a command")
	}
}

func TestModel_Update_KeyMsg_Quit(t *testing.T) {
	eventChan := make(chan interface{})
	defer close(eventChan)

	model := New(eventChan, testConfig(), "/test/project")

	// Test 'q' key
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if updatedModel == nil {
		t.Error("Expected non-nil model")
	}
	if cmd == nil {
		t.Error("Expected quit command")
	}

	// Test Ctrl+C
	updatedModel, cmd = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if updatedModel == nil {
		t.Error("Expected non-nil model")
	}
	if cmd == nil {
		t.Error("Expected quit command")
	}
}

func TestModel_Update_KeyMsg_ToggleAutoFollow(t *testing.T) {
	eventChan := make(chan interface{})
	defer close(eventChan)

	model := New(eventChan, testConfig(), "/test/project")
	initialAutoFollow := model.autoFollow

	// Press 'f' to toggle
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	m := updatedModel.(*Model)

	if m.autoFollow == initialAutoFollow {
		t.Errorf("Expected autoFollow to toggle from %v to %v", initialAutoFollow, !initialAutoFollow)
	}

	// Press 'f' again to toggle back
	updatedModel2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	m2 := updatedModel2.(*Model)

	if m2.autoFollow != initialAutoFollow {
		t.Errorf("Expected autoFollow to toggle back to %v", initialAutoFollow)
	}
}

func TestModel_Update_KeyMsg_ViewportScroll(t *testing.T) {
	eventChan := make(chan interface{})
	defer close(eventChan)

	model := New(eventChan, testConfig(), "/test/project")
	model.autoFollow = true

	// Press 'up' arrow (should disable auto-follow)
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	m := updatedModel.(*Model)
	if m.autoFollow {
		t.Error("Expected autoFollow to be disabled after pressing up arrow")
	}

	// Reset
	model.autoFollow = true

	// Press 'pgup' (should disable auto-follow)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	m = updatedModel.(*Model)
	if m.autoFollow {
		t.Error("Expected autoFollow to be disabled after pressing pgup")
	}
}

func TestModel_Update_WindowSizeMsg(t *testing.T) {
	eventChan := make(chan interface{})
	defer close(eventChan)

	model := New(eventChan, testConfig(), "/test/project")

	msg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(*Model)

	if m.width != 100 {
		t.Errorf("Expected width to be 100, got %d", m.width)
	}
	if m.height != 30 {
		t.Errorf("Expected height to be 30, got %d", m.height)
	}
}

func TestModel_Update_RunStartedMsg(t *testing.T) {
	eventChan := make(chan interface{})
	defer close(eventChan)

	model := New(eventChan, testConfig(), "/test/project")

	msg := driver.RunStartedMsg{
		MaxIterations: 10,
		ReadyCount:    5,
	}
	updatedModel, cmd := model.Update(msg)
	m := updatedModel.(*Model)

	if m.maxIterations != 10 {
		t.Errorf("Expected maxIterations to be 10, got %d", m.maxIterations)
	}
	if cmd == nil {
		t.Error("Expected command to wait for next event")
	}
}

func TestModel_Update_IterationStartedMsg(t *testing.T) {
	eventChan := make(chan interface{})
	defer close(eventChan)

	model := New(eventChan, testConfig(), "/test/project")
	model.progressContent = "previous content"

	msg := driver.IterationStartedMsg{
		Iteration: 3,
		TaskID:    "task-123",
		TaskTitle: "Test Task",
	}
	updatedModel, cmd := model.Update(msg)
	m := updatedModel.(*Model)

	if m.iteration != 3 {
		t.Errorf("Expected iteration to be 3, got %d", m.iteration)
	}
	if m.activeTaskID != "task-123" {
		t.Errorf("Expected activeTaskID to be 'task-123', got %s", m.activeTaskID)
	}
	// Progress content is not cleared on iteration start anymore (file-based display)
	if cmd == nil {
		t.Error("Expected command to wait for next event")
	}
}

func TestModel_Update_OutputMsg(t *testing.T) {
	eventChan := make(chan interface{})
	defer close(eventChan)

	model := New(eventChan, testConfig(), "/test/project")
	model.autoFollow = true

	msg := driver.OutputMsg{
		Data: []byte("Hello, World!"),
	}
	_, cmd := model.Update(msg)

	// OutputMsg is now a no-op since output is handled via file watching
	if cmd == nil {
		t.Error("Expected command to wait for next event")
	}

	// The progress content is not updated directly by OutputMsg anymore
	// It's updated via FileChangedMsg when the file watcher detects changes
}

func TestModel_Update_IterationCompleteMsg(t *testing.T) {
	eventChan := make(chan interface{})
	defer close(eventChan)

	model := New(eventChan, testConfig(), "/test/project")
	model.activeTaskID = "task-123"

	msg := driver.IterationCompleteMsg{
		Status: driver.IterationStatusCompleted,
	}
	updatedModel, cmd := model.Update(msg)
	m := updatedModel.(*Model)

	if m.activeTaskID != "" {
		t.Errorf("Expected activeTaskID to be cleared, got %s", m.activeTaskID)
	}
	if cmd == nil {
		t.Error("Expected command to wait for next event")
	}
}

func TestModel_Update_TasksRefreshedMsg(t *testing.T) {
	eventChan := make(chan interface{})
	defer close(eventChan)

	model := New(eventChan, testConfig(), "/test/project")

	newTasks := []tasks.Task{
		{ID: "task-1", Title: "Task 1", Status: "done"},
		{ID: "task-2", Title: "Task 2", Status: "pending"},
	}

	msg := driver.TasksRefreshedMsg{
		Tasks: newTasks,
	}
	updatedModel, cmd := model.Update(msg)
	m := updatedModel.(*Model)

	if len(m.tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(m.tasks))
	}
	if m.tasks[0].ID != "task-1" {
		t.Errorf("Expected first task ID to be 'task-1', got %s", m.tasks[0].ID)
	}
	if m.tasks[1].ID != "task-2" {
		t.Errorf("Expected second task ID to be 'task-2', got %s", m.tasks[1].ID)
	}
	if cmd == nil {
		t.Error("Expected command to wait for next event")
	}
}

func TestModel_Update_RunCompleteMsg(t *testing.T) {
	eventChan := make(chan interface{})
	defer close(eventChan)

	model := New(eventChan, testConfig(), "/test/project")

	msg := driver.RunCompleteMsg{
		Status:   driver.RunStatusCompleted,
		Metadata: nil,
	}
	updatedModel, cmd := model.Update(msg)
	m := updatedModel.(*Model)

	if !m.runComplete {
		t.Error("Expected runComplete to be true")
	}
	if m.finalStatus != driver.RunStatusCompleted {
		t.Errorf("Expected finalStatus to be RunStatusCompleted, got %s", m.finalStatus)
	}
	// cmd might be a batch, but no new waitForEvent should be added after run complete
	_ = cmd
}

func TestModel_View_Uninitialized(t *testing.T) {
	eventChan := make(chan interface{})
	defer close(eventChan)

	model := New(eventChan, testConfig(), "/test/project")
	// Don't set width/height

	view := model.View()
	if view != "Initializing..." {
		t.Errorf("Expected 'Initializing...' for uninitialized model, got %s", view)
	}
}

func TestModel_View_Initialized(t *testing.T) {
	eventChan := make(chan interface{})
	defer close(eventChan)

	model := New(eventChan, testConfig(), "/test/project")
	model.width = 100
	model.height = 30
	model.tasks = []tasks.Task{
		{ID: "task-1", Title: "Task 1", Status: "pending"},
	}

	view := model.View()

	// Should contain task information
	if view == "" {
		t.Error("Expected non-empty view")
	}
	// The view should not be "Initializing..." anymore
	if view == "Initializing..." {
		t.Error("Expected initialized view, not 'Initializing...'")
	}
}

func TestModel_RecalculateLayout(t *testing.T) {
	eventChan := make(chan interface{})
	defer close(eventChan)

	model := New(eventChan, testConfig(), "/test/project")
	model.width = 100
	model.height = 30

	model.recalculateLayout()

	// Output panel should get width = total width - left panel width
	// This is a simple test to ensure no panic occurs
	// Note: We can't directly check the output panel dimensions without exposing them
	// But we can verify the method doesn't panic
	if model.outputPanel == nil {
		t.Error("Expected outputPanel to remain initialized")
	}

	// Test with small width
	model.width = 10
	model.recalculateLayout()
	// Should handle gracefully (width might be 0 or negative, clamped to 0)
}

func TestWaitForEvent(t *testing.T) {
	eventChan := make(chan interface{}, 1)

	// Send a test message
	testMsg := driver.RunStartedMsg{MaxIterations: 5}
	eventChan <- testMsg

	cmd := waitForEvent(eventChan)
	if cmd == nil {
		t.Fatal("Expected non-nil command")
	}

	// Execute the command
	result := cmd()
	if result == nil {
		t.Error("Expected message from command")
	}

	// Verify it's the correct message
	if msg, ok := result.(driver.RunStartedMsg); !ok {
		t.Errorf("Expected RunStartedMsg, got %T", result)
	} else if msg.MaxIterations != 5 {
		t.Errorf("Expected MaxIterations to be 5, got %d", msg.MaxIterations)
	}

	// Close channel and test that it returns nil
	close(eventChan)
	cmd2 := waitForEvent(eventChan)
	result2 := cmd2()
	if result2 != nil {
		t.Errorf("Expected nil when channel is closed, got %v", result2)
	}
}
