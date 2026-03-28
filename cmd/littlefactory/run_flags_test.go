package main

import (
	"encoding/json"
	"os"
	"os/exec"
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
	if got := err.Error(); got != `change "nonexistent" not found at .littlefactory/changes/nonexistent/` {
		t.Errorf("unexpected error: %s", got)
	}
}

func TestValidateChangeFlags_ChangeNoTasksJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// Create change directory but no tasks.json
	changeDir := filepath.Join(tmpDir, ".littlefactory", "changes", "incomplete-change")
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
	changeDir := filepath.Join(tmpDir, ".littlefactory", "changes", "feature-a")
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
	changeDir := filepath.Join(tmpDir, ".littlefactory", "changes", "feature-a")
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

// gitRun executes a git command in the given directory and fails the test on error.
func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func TestPrepareWorktree_ReusesExisting(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Initialize git repo
	gitRun(t, repoDir, "init")
	gitRun(t, repoDir, "config", "user.email", "test@test.com")
	gitRun(t, repoDir, "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(repoDir, "README.md"), []byte("init"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, repoDir, "add", ".")
	gitRun(t, repoDir, "commit", "-m", "initial commit")

	// Create a worktree manually first
	wtDir := filepath.Join(tmpDir, "my-feature")
	gitRun(t, repoDir, "worktree", "add", wtDir, "-b", "my-feature")

	// Now prepareWorktree should reuse it instead of erroring
	got, err := prepareWorktree(repoDir, "my-feature", tmpDir)
	if err != nil {
		t.Fatalf("expected reuse, got error: %v", err)
	}

	// Resolve symlinks for comparison (macOS /var -> /private/var)
	resolvedGot, _ := filepath.EvalSymlinks(got)
	resolvedExpected, _ := filepath.EvalSymlinks(wtDir)
	if resolvedGot != resolvedExpected {
		t.Errorf("expected reused path %q, got %q", resolvedExpected, resolvedGot)
	}
}

func TestPrepareWorktree_CreatesNew(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatal(err)
	}

	gitRun(t, repoDir, "init")
	gitRun(t, repoDir, "config", "user.email", "test@test.com")
	gitRun(t, repoDir, "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(repoDir, "README.md"), []byte("init"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, repoDir, "add", ".")
	gitRun(t, repoDir, "commit", "-m", "initial commit")

	// prepareWorktree should create a new worktree
	got, err := prepareWorktree(repoDir, "new-feature", tmpDir)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "new-feature")
	if got != expectedPath {
		t.Errorf("expected path %q, got %q", expectedPath, got)
	}

	// Verify directory exists
	if _, err := os.Stat(got); os.IsNotExist(err) {
		t.Error("worktree directory was not created")
	}
}
