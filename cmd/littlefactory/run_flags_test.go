package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateChangeFlags_WorktreeRequiresChange(t *testing.T) {
	err := validateChangeFlags("/tmp/fake", "", "", true)
	if err == nil {
		t.Fatal("expected error when -w is set without -c")
	}
	expected := "the --worktree flag requires --change to specify the branch name"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestValidateChangeFlags_NoFlags(t *testing.T) {
	err := validateChangeFlags("/tmp/fake", "", "", false)
	if err != nil {
		t.Fatalf("expected no error with no flags, got: %v", err)
	}
}

func TestValidateChangeFlags_ChangeNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	err := validateChangeFlags(tmpDir, "nonexistent", "", false)
	if err == nil {
		t.Fatal("expected error for nonexistent change")
	}
	if got := err.Error(); got != `change "nonexistent" not found at openspec/changes/nonexistent/` {
		t.Errorf("unexpected error: %s", got)
	}
}

func TestValidateChangeFlags_ChangeNoTasksJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// Create change directory but no tasks.json
	changeDir := filepath.Join(tmpDir, "openspec", "changes", "incomplete-change")
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatal(err)
	}

	err := validateChangeFlags(tmpDir, "incomplete-change", "", false)
	if err == nil {
		t.Fatal("expected error for missing tasks.json")
	}
	if got := err.Error(); got != `no tasks.json found for change "incomplete-change"` {
		t.Errorf("unexpected error: %s", got)
	}
}

func TestValidateChangeFlags_ValidChange(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid change directory with tasks.json
	changeDir := filepath.Join(tmpDir, "openspec", "changes", "feature-a")
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatal(err)
	}

	tasksData, _ := json.Marshal(map[string]interface{}{
		"tasks": []map[string]string{
			{"id": "001", "title": "Test", "status": "todo"},
		},
	})
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.json"), tasksData, 0644); err != nil {
		t.Fatal(err)
	}

	err := validateChangeFlags(tmpDir, "feature-a", "", false)
	if err != nil {
		t.Fatalf("expected no error for valid change, got: %v", err)
	}
}

func TestValidateChangeFlags_ValidChangeWithWorktree(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid change directory with tasks.json
	changeDir := filepath.Join(tmpDir, "openspec", "changes", "feature-a")
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatal(err)
	}

	tasksData, _ := json.Marshal(map[string]interface{}{
		"tasks": []map[string]string{
			{"id": "001", "title": "Test", "status": "todo"},
		},
	})
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.json"), tasksData, 0644); err != nil {
		t.Fatal(err)
	}

	// -c with -w should pass validation (worktree creation checks happen later)
	err := validateChangeFlags(tmpDir, "feature-a", "", true)
	if err != nil {
		t.Fatalf("expected no error for valid change with -w, got: %v", err)
	}
}

func TestRunCmd_HasChangeFlag(t *testing.T) {
	flag := runCmd.Flags().Lookup("change")
	if flag == nil {
		t.Fatal("expected --change flag to be registered")
	}
	if flag.Shorthand != "c" {
		t.Errorf("expected shorthand 'c', got '%s'", flag.Shorthand)
	}
}

func TestRunCmd_HasWorktreeFlag(t *testing.T) {
	flag := runCmd.Flags().Lookup("worktree")
	if flag == nil {
		t.Fatal("expected --worktree flag to be registered")
	}
	if flag.Shorthand != "w" {
		t.Errorf("expected shorthand 'w', got '%s'", flag.Shorthand)
	}
}

func TestRunCmd_HasTasksFlag(t *testing.T) {
	flag := runCmd.Flags().Lookup("tasks")
	if flag == nil {
		t.Fatal("expected --tasks flag to be registered")
	}
	if flag.Shorthand != "t" {
		t.Errorf("expected shorthand 't', got '%s'", flag.Shorthand)
	}
}

func TestValidateChangeFlags_TasksFileNotFound(t *testing.T) {
	err := validateChangeFlags("/tmp/fake", "", "/nonexistent/tasks.json", false)
	if err == nil {
		t.Fatal("expected error for nonexistent tasks file")
	}
	expected := "tasks file not found: /nonexistent/tasks.json"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestValidateChangeFlags_TasksFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	tasksFile := filepath.Join(tmpDir, "custom-tasks.json")
	if err := os.WriteFile(tasksFile, []byte(`{"tasks":[]}`), 0644); err != nil {
		t.Fatal(err)
	}

	err := validateChangeFlags("/tmp/fake", "", tasksFile, false)
	if err != nil {
		t.Fatalf("expected no error for existing tasks file, got: %v", err)
	}
}

func TestValidateChangeFlags_TasksOverridesChange(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a valid tasks file
	tasksFile := filepath.Join(tmpDir, "custom-tasks.json")
	if err := os.WriteFile(tasksFile, []byte(`{"tasks":[]}`), 0644); err != nil {
		t.Fatal(err)
	}

	// --tasks with a nonexistent --change should still pass because --tasks takes priority
	err := validateChangeFlags(tmpDir, "nonexistent-change", tasksFile, false)
	if err != nil {
		t.Fatalf("expected no error when --tasks is valid (overrides --change), got: %v", err)
	}
}

func TestValidateChangeFlags_TasksNotFoundErrorMessage(t *testing.T) {
	err := validateChangeFlags("/tmp/fake", "", "nonexistent.json", false)
	if err == nil {
		t.Fatal("expected error for nonexistent tasks file")
	}
	if !strings.Contains(err.Error(), "tasks file not found") {
		t.Errorf("expected 'tasks file not found' in error, got: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "nonexistent.json") {
		t.Errorf("expected file path in error, got: %s", err.Error())
	}
}
