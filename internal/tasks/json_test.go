package tasks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/littlefactory/internal/config"
)

func TestJSONTaskSource_Interface(t *testing.T) {
	// Verify JSONTaskSource implements TaskSource interface
	var _ TaskSource = (*JSONTaskSource)(nil)
}

func TestJSONTaskSource_ReadWriteCycle(t *testing.T) {
	// Create temp directory for testing
	tmpDir, err := os.MkdirTemp("", "json-task-source-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test config
	cfg := &config.Config{
		StateDir: ".littlefactory",
	}

	// Create task source (no file exists yet, returns empty valid list)
	ts, err := NewJSONTaskSource(tmpDir, cfg)
	if err != nil {
		t.Fatalf("NewJSONTaskSource() failed: %v", err)
	}

	// Initially, should have no tasks
	tasks, err := ts.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}

	// Write some tasks (valid sequential chain)
	testTasks := []Task{
		{ID: "001", Title: "Task 1", Description: "First task", Status: "todo"},
		{ID: "002", Title: "Task 2", Description: "Second task", Status: "todo", Blockers: []string{"001"}},
		{ID: "003", Title: "Task 3", Description: "Third task", Status: "done", Blockers: []string{"002"}},
	}
	if err := ts.writeTasks(testTasks); err != nil {
		t.Fatalf("writeTasks() failed: %v", err)
	}

	// Verify .littlefactory directory was created
	dirPath := filepath.Join(tmpDir, ".littlefactory")
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Error(".littlefactory directory was not created")
	}

	// Read tasks back
	tasks, err = ts.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}
}

func TestJSONTaskSource_Ready(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "json-task-source-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		StateDir: ".littlefactory",
	}

	ts, err := NewJSONTaskSource(tmpDir, cfg)
	if err != nil {
		t.Fatalf("NewJSONTaskSource() failed: %v", err)
	}

	// Write tasks with different statuses (valid sequential chain)
	testTasks := []Task{
		{ID: "001", Title: "Done task", Description: "Already done", Status: "done"},
		{ID: "002", Title: "In progress", Description: "Working on it", Status: "in_progress", Blockers: []string{"001"}},
		{ID: "003", Title: "Ready task", Description: "Should be returned", Status: "todo", Blockers: []string{"002"}},
		{ID: "004", Title: "Another ready", Description: "Should not be returned", Status: "todo", Blockers: []string{"003"}},
	}
	if err := ts.writeTasks(testTasks); err != nil {
		t.Fatalf("writeTasks() failed: %v", err)
	}

	// Ready should return first "todo" task
	ready, err := ts.Ready()
	if err != nil {
		t.Fatalf("Ready() failed: %v", err)
	}
	if len(ready) != 1 {
		t.Errorf("Expected 1 ready task, got %d", len(ready))
	}
	if len(ready) > 0 && ready[0].ID != "003" {
		t.Errorf("Expected task 003, got %s", ready[0].ID)
	}
}

func TestJSONTaskSource_Show(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "json-task-source-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		StateDir: ".littlefactory",
	}

	ts, err := NewJSONTaskSource(tmpDir, cfg)
	if err != nil {
		t.Fatalf("NewJSONTaskSource() failed: %v", err)
	}

	testTasks := []Task{
		{ID: "001", Title: "Task 1", Description: "First task", Status: "todo"},
		{ID: "002", Title: "Task 2", Description: "Second task", Status: "todo", Blockers: []string{"001"}},
	}
	if err := ts.writeTasks(testTasks); err != nil {
		t.Fatalf("writeTasks() failed: %v", err)
	}

	// Show existing task
	task, err := ts.Show("002")
	if err != nil {
		t.Fatalf("Show() failed: %v", err)
	}
	if task.Title != "Task 2" {
		t.Errorf("Expected title 'Task 2', got '%s'", task.Title)
	}

	// Show non-existent task
	_, err = ts.Show("999")
	if err == nil {
		t.Error("Expected error for non-existent task, got nil")
	}
}

func TestJSONTaskSource_Claim(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "json-task-source-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		StateDir: ".littlefactory",
	}

	ts, err := NewJSONTaskSource(tmpDir, cfg)
	if err != nil {
		t.Fatalf("NewJSONTaskSource() failed: %v", err)
	}

	testTasks := []Task{
		{ID: "001", Title: "Task 1", Description: "First task", Status: "todo"},
	}
	if err := ts.writeTasks(testTasks); err != nil {
		t.Fatalf("writeTasks() failed: %v", err)
	}

	// Claim task
	if err := ts.Claim("001"); err != nil {
		t.Fatalf("Claim() failed: %v", err)
	}

	// Verify status changed
	task, err := ts.Show("001")
	if err != nil {
		t.Fatalf("Show() failed: %v", err)
	}
	if task.Status != "in_progress" {
		t.Errorf("Expected status 'in_progress', got '%s'", task.Status)
	}
}

func TestJSONTaskSource_Close(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "json-task-source-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		StateDir: ".littlefactory",
	}

	ts, err := NewJSONTaskSource(tmpDir, cfg)
	if err != nil {
		t.Fatalf("NewJSONTaskSource() failed: %v", err)
	}

	testTasks := []Task{
		{ID: "001", Title: "Task 1", Description: "First task", Status: "in_progress"},
	}
	if err := ts.writeTasks(testTasks); err != nil {
		t.Fatalf("writeTasks() failed: %v", err)
	}

	// Close task
	if err := ts.Close("001", "Completed"); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	// Verify status changed
	task, err := ts.Show("001")
	if err != nil {
		t.Fatalf("Show() failed: %v", err)
	}
	if task.Status != "done" {
		t.Errorf("Expected status 'done', got '%s'", task.Status)
	}
}

func TestJSONTaskSource_Reset(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "json-task-source-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		StateDir: ".littlefactory",
	}

	ts, err := NewJSONTaskSource(tmpDir, cfg)
	if err != nil {
		t.Fatalf("NewJSONTaskSource() failed: %v", err)
	}

	testTasks := []Task{
		{ID: "001", Title: "Task 1", Description: "First task", Status: "in_progress"},
	}
	if err := ts.writeTasks(testTasks); err != nil {
		t.Fatalf("writeTasks() failed: %v", err)
	}

	// Reset task
	if err := ts.Reset("001"); err != nil {
		t.Fatalf("Reset() failed: %v", err)
	}

	// Verify status changed
	task, err := ts.Show("001")
	if err != nil {
		t.Fatalf("Show() failed: %v", err)
	}
	if task.Status != "todo" {
		t.Errorf("Expected status 'todo', got '%s'", task.Status)
	}
}

func TestNewJSONTaskSourceWithPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "json-task-source-path-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a tasks.json file at a custom path
	tasksPath := filepath.Join(tmpDir, "custom", "tasks.json")
	if err := os.MkdirAll(filepath.Dir(tasksPath), 0755); err != nil {
		t.Fatal(err)
	}

	// Write test tasks (valid sequential chain)
	data := `{"tasks": [{"id": "001", "title": "Custom Task", "description": "From custom path", "status": "todo"}]}`
	if err := os.WriteFile(tasksPath, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	ts, err := NewJSONTaskSourceWithPath(tasksPath)
	if err != nil {
		t.Fatalf("NewJSONTaskSourceWithPath() failed: %v", err)
	}

	// Verify it reads from the custom path
	tasks, err := ts.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tasks))
	}
	if tasks[0].Title != "Custom Task" {
		t.Errorf("Expected title 'Custom Task', got '%s'", tasks[0].Title)
	}

	// Verify interface compliance
	var _ TaskSource = ts
}

func TestNewJSONTaskSourceWithPath_FileNotFound(t *testing.T) {
	_, err := NewJSONTaskSourceWithPath("/nonexistent/tasks.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
	if !strings.Contains(err.Error(), "tasks file not found") {
		t.Errorf("expected 'tasks file not found' in error, got: %v", err)
	}
}

func TestNewJSONTaskSourceWithPath_InvalidContent(t *testing.T) {
	tmpDir := t.TempDir()
	tasksPath := filepath.Join(tmpDir, "tasks.json")

	// Write invalid JSON
	if err := os.WriteFile(tasksPath, []byte("not json"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := NewJSONTaskSourceWithPath(tasksPath)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "failed to parse tasks file") {
		t.Errorf("expected parse error, got: %v", err)
	}
}

func TestNewJSONTaskSource_InvalidContent(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		StateDir: ".littlefactory",
	}

	// Create state directory with invalid tasks.json
	stateDir := filepath.Join(tmpDir, ".littlefactory")
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Write tasks with missing required fields
	data := `{"tasks": [{"id": "", "title": "No ID", "status": "todo"}]}`
	if err := os.WriteFile(filepath.Join(stateDir, "tasks.json"), []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := NewJSONTaskSource(tmpDir, cfg)
	if err == nil {
		t.Fatal("expected error for invalid tasks.json")
	}
	if !strings.Contains(err.Error(), "missing required field") {
		t.Errorf("expected validation error, got: %v", err)
	}
}

// --- ValidateTasks tests ---

func TestValidateTasks_ValidTasks(t *testing.T) {
	taskList := []Task{
		{ID: "001", Title: "Task 1", Status: "todo"},
		{ID: "002", Title: "Task 2", Status: "in_progress", Blockers: []string{"001"}},
		{ID: "003", Title: "Task 3", Status: "done", Blockers: []string{"002"}},
	}
	err := ValidateTasks(taskList, "test.json")
	if err != nil {
		t.Errorf("expected no error for valid tasks, got: %v", err)
	}
}

func TestValidateTasks_EmptyList(t *testing.T) {
	err := ValidateTasks([]Task{}, "test.json")
	if err != nil {
		t.Errorf("expected no error for empty list, got: %v", err)
	}
}

func TestValidateTasks_MissingID(t *testing.T) {
	taskList := []Task{
		{ID: "", Title: "No ID", Status: "todo"},
	}
	err := ValidateTasks(taskList, "test.json")
	if err == nil {
		t.Fatal("expected error for missing id")
	}
	if !strings.Contains(err.Error(), `missing required field "id"`) {
		t.Errorf("expected missing id error, got: %v", err)
	}
}

func TestValidateTasks_MissingTitle(t *testing.T) {
	taskList := []Task{
		{ID: "001", Title: "", Status: "todo"},
	}
	err := ValidateTasks(taskList, "test.json")
	if err == nil {
		t.Fatal("expected error for missing title")
	}
	if !strings.Contains(err.Error(), `missing required field "title"`) {
		t.Errorf("expected missing title error, got: %v", err)
	}
}

func TestValidateTasks_MissingStatus(t *testing.T) {
	taskList := []Task{
		{ID: "001", Title: "Task", Status: ""},
	}
	err := ValidateTasks(taskList, "test.json")
	if err == nil {
		t.Fatal("expected error for missing status")
	}
	if !strings.Contains(err.Error(), `missing required field "status"`) {
		t.Errorf("expected missing status error, got: %v", err)
	}
}

func TestValidateTasks_InvalidStatus(t *testing.T) {
	taskList := []Task{
		{ID: "001", Title: "Task", Status: "pending"},
	}
	err := ValidateTasks(taskList, "test.json")
	if err == nil {
		t.Fatal("expected error for invalid status")
	}
	if !strings.Contains(err.Error(), `invalid status "pending"`) {
		t.Errorf("expected invalid status error, got: %v", err)
	}
}

func TestValidateTasks_DuplicateID(t *testing.T) {
	taskList := []Task{
		{ID: "001", Title: "Task 1", Status: "todo"},
		{ID: "001", Title: "Task 2", Status: "todo", Blockers: []string{"001"}},
	}
	err := ValidateTasks(taskList, "test.json")
	if err == nil {
		t.Fatal("expected error for duplicate id")
	}
	if !strings.Contains(err.Error(), "duplicate id") {
		t.Errorf("expected duplicate id error, got: %v", err)
	}
}

func TestValidateTasks_MultipleRoots(t *testing.T) {
	taskList := []Task{
		{ID: "001", Title: "Task 1", Status: "todo"},
		{ID: "002", Title: "Task 2", Status: "todo"},
	}
	err := ValidateTasks(taskList, "test.json")
	if err == nil {
		t.Fatal("expected error for multiple roots")
	}
	if !strings.Contains(err.Error(), "multiple root tasks") {
		t.Errorf("expected multiple roots error, got: %v", err)
	}
}

func TestValidateTasks_NoRoot(t *testing.T) {
	taskList := []Task{
		{ID: "001", Title: "Task 1", Status: "todo", Blockers: []string{"002"}},
		{ID: "002", Title: "Task 2", Status: "todo", Blockers: []string{"001"}},
	}
	err := ValidateTasks(taskList, "test.json")
	if err == nil {
		t.Fatal("expected error for no root")
	}
	if !strings.Contains(err.Error(), "no root task") {
		t.Errorf("expected no root error, got: %v", err)
	}
}

func TestValidateTasks_MultipleBlockers(t *testing.T) {
	taskList := []Task{
		{ID: "001", Title: "Task 1", Status: "todo"},
		{ID: "002", Title: "Task 2", Status: "todo", Blockers: []string{"001"}},
		{ID: "003", Title: "Task 3", Status: "todo", Blockers: []string{"001", "002"}},
	}
	err := ValidateTasks(taskList, "test.json")
	if err == nil {
		t.Fatal("expected error for multiple blockers")
	}
	if !strings.Contains(err.Error(), "has 2 blockers") {
		t.Errorf("expected multiple blockers error, got: %v", err)
	}
}

func TestValidateTasks_NonExistentBlocker(t *testing.T) {
	taskList := []Task{
		{ID: "001", Title: "Task 1", Status: "todo"},
		{ID: "002", Title: "Task 2", Status: "todo", Blockers: []string{"999"}},
	}
	err := ValidateTasks(taskList, "test.json")
	if err == nil {
		t.Fatal("expected error for nonexistent blocker")
	}
	if !strings.Contains(err.Error(), `blocker "999" does not exist`) {
		t.Errorf("expected nonexistent blocker error, got: %v", err)
	}
}

func TestValidateTasks_OrphanedTasks(t *testing.T) {
	// 001 -> 002, but 003 points to itself (not reachable from root)
	taskList := []Task{
		{ID: "001", Title: "Task 1", Status: "todo"},
		{ID: "002", Title: "Task 2", Status: "todo", Blockers: []string{"001"}},
		{ID: "003", Title: "Task 3", Status: "todo", Blockers: []string{"003"}},
	}
	err := ValidateTasks(taskList, "test.json")
	if err == nil {
		t.Fatal("expected error for orphaned tasks")
	}
	if !strings.Contains(err.Error(), "not reachable from root") {
		t.Errorf("expected orphaned tasks error, got: %v", err)
	}
}

func TestValidateTasks_MultiError(t *testing.T) {
	taskList := []Task{
		{ID: "001", Title: "", Status: "invalid"},
		{ID: "001", Title: "Dup", Status: ""},
	}
	err := ValidateTasks(taskList, "test.json")
	if err == nil {
		t.Fatal("expected multi-error")
	}
	errStr := err.Error()
	// Should contain the file path header
	if !strings.Contains(errStr, "Error loading tasks from test.json:") {
		t.Errorf("expected file path header, got: %v", err)
	}
	// Should contain multiple errors
	if !strings.Contains(errStr, `missing required field "title"`) {
		t.Errorf("expected missing title error, got: %v", err)
	}
	if !strings.Contains(errStr, `invalid status "invalid"`) {
		t.Errorf("expected invalid status error, got: %v", err)
	}
	if !strings.Contains(errStr, "duplicate id") {
		t.Errorf("expected duplicate id error, got: %v", err)
	}
}

func TestValidateTasks_EmptyDescriptionAllowed(t *testing.T) {
	taskList := []Task{
		{ID: "001", Title: "Task 1", Description: "", Status: "todo"},
	}
	err := ValidateTasks(taskList, "test.json")
	if err != nil {
		t.Errorf("expected no error for empty description, got: %v", err)
	}
}

func TestValidateTasks_AllValidStatuses(t *testing.T) {
	taskList := []Task{
		{ID: "001", Title: "Task 1", Status: "todo"},
		{ID: "002", Title: "Task 2", Status: "in_progress", Blockers: []string{"001"}},
		{ID: "003", Title: "Task 3", Status: "done", Blockers: []string{"002"}},
	}
	err := ValidateTasks(taskList, "test.json")
	if err != nil {
		t.Errorf("expected no error for valid statuses, got: %v", err)
	}
}
