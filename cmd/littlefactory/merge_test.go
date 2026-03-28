package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// initTestRepo creates a git repo with an initial commit and returns its path.
func initTestRepo(t *testing.T, baseDir string) string {
	t.Helper()
	repoDir := filepath.Join(baseDir, "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatal(err)
	}
	gitRun(t, repoDir, "init")
	gitRun(t, repoDir, "config", "user.email", "test@test.com")
	gitRun(t, repoDir, "config", "user.name", "Test")
	// Ensure we're on "main" branch
	gitRun(t, repoDir, "checkout", "-b", "main")
	if err := os.WriteFile(filepath.Join(repoDir, "README.md"), []byte("init"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, repoDir, "add", ".")
	gitRun(t, repoDir, "commit", "-m", "initial commit")
	return repoDir
}

// createWorktreeWithTasks sets up a worktree branch with a change dir and tasks.json.
// Returns the worktree path.
func createWorktreeWithTasks(t *testing.T, repoDir, changeName string, taskList []map[string]interface{}) string {
	t.Helper()
	tmpDir := filepath.Dir(repoDir)
	wtDir := filepath.Join(tmpDir, "wt-"+changeName)
	gitRun(t, repoDir, "worktree", "add", wtDir, "-b", changeName)

	// Create change directory with tasks.json in the worktree
	changeDir := filepath.Join(wtDir, ".littlefactory", "changes", changeName)
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	tasksData, err := json.MarshalIndent(map[string]interface{}{"tasks": taskList}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.json"), tasksData, 0o644); err != nil {
		t.Fatal(err)
	}

	// Commit the tasks.json in the worktree
	gitRun(t, wtDir, "add", ".")
	gitRun(t, wtDir, "commit", "-m", "add tasks")

	return wtDir
}

// --- Flag registration tests ---

func TestMergeCmd_HasChangeFlag(t *testing.T) {
	flag := mergeCmd.Flags().Lookup("change")
	if flag == nil {
		t.Fatal("expected --change flag to be registered")
	}
	if flag.Shorthand != "c" {
		t.Errorf("expected shorthand 'c', got '%s'", flag.Shorthand)
	}
}

func TestMergeCmd_ChangeFlagRequired(t *testing.T) {
	flag := mergeCmd.Flags().Lookup("change")
	if flag == nil {
		t.Fatal("expected --change flag to be registered")
	}
	annotations := flag.Annotations
	if annotations == nil {
		t.Fatal("expected --change flag to have annotations (required)")
	}
	if _, ok := annotations["cobra_annotation_bash_completion_one_required_flag"]; !ok {
		t.Error("expected --change flag to be marked as required")
	}
}

func TestMergeCmd_HasForceFlag(t *testing.T) {
	flag := mergeCmd.Flags().Lookup("force")
	if flag == nil {
		t.Fatal("expected --force flag to be registered")
	}
	if flag.Shorthand != "f" {
		t.Errorf("expected shorthand 'f', got '%s'", flag.Shorthand)
	}
}

func TestMergeCmd_HasMaxVerifyRetriesFlag(t *testing.T) {
	flag := mergeCmd.Flags().Lookup("max-verify-retries")
	if flag == nil {
		t.Fatal("expected --max-verify-retries flag to be registered")
	}
	if flag.DefValue != "3" {
		t.Errorf("expected default value '3', got '%s'", flag.DefValue)
	}
}

// --- checkAllTasksDone tests ---

func TestCheckAllTasksDone_AllDone(t *testing.T) {
	tmpDir := t.TempDir()
	changeName := "test-change"
	changeDir := filepath.Join(tmpDir, ".littlefactory", "changes", changeName)
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	tasks := []map[string]interface{}{
		{"id": "t1", "title": "Task 1", "status": "done"},
		{"id": "t2", "title": "Task 2", "status": "done", "blockers": []string{"t1"}},
	}
	data, _ := json.Marshal(map[string]interface{}{"tasks": tasks})
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	err := checkAllTasksDone(tmpDir, changeName)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestCheckAllTasksDone_IncompleteTasks(t *testing.T) {
	tmpDir := t.TempDir()
	changeName := "test-change"
	changeDir := filepath.Join(tmpDir, ".littlefactory", "changes", changeName)
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	tasks := []map[string]interface{}{
		{"id": "t1", "title": "Done Task", "status": "done"},
		{"id": "t2", "title": "Incomplete Task", "status": "todo", "blockers": []string{"t1"}},
	}
	data, _ := json.Marshal(map[string]interface{}{"tasks": tasks})
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	err := checkAllTasksDone(tmpDir, changeName)
	if err == nil {
		t.Fatal("expected error for incomplete tasks")
	}
	if !strings.Contains(err.Error(), "Incomplete Task") {
		t.Errorf("expected error to mention incomplete task, got: %v", err)
	}
	if !strings.Contains(err.Error(), "--force") {
		t.Errorf("expected error to mention --force, got: %v", err)
	}
}

func TestCheckAllTasksDone_NoTasksFile(t *testing.T) {
	tmpDir := t.TempDir()
	err := checkAllTasksDone(tmpDir, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing tasks.json")
	}
}

// --- rebaseIfNeeded tests ---

func TestRebaseIfNeeded_NoRebaseNeeded(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := initTestRepo(t, tmpDir)

	// Create worktree with a commit
	wtDir := filepath.Join(tmpDir, "wt-feature")
	gitRun(t, repoDir, "worktree", "add", wtDir, "-b", "feature")
	if err := os.WriteFile(filepath.Join(wtDir, "feature.txt"), []byte("feature"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, wtDir, "add", ".")
	gitRun(t, wtDir, "commit", "-m", "feature work")

	// main has NOT advanced, so no rebase needed
	err := rebaseIfNeeded(repoDir, wtDir, "feature")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestRebaseIfNeeded_RebaseRequired(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := initTestRepo(t, tmpDir)

	// Create worktree
	wtDir := filepath.Join(tmpDir, "wt-feature")
	gitRun(t, repoDir, "worktree", "add", wtDir, "-b", "feature")
	if err := os.WriteFile(filepath.Join(wtDir, "feature.txt"), []byte("feature"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, wtDir, "add", ".")
	gitRun(t, wtDir, "commit", "-m", "feature work")

	// Advance main with a non-conflicting commit
	if err := os.WriteFile(filepath.Join(repoDir, "main-only.txt"), []byte("main advance"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, repoDir, "add", ".")
	gitRun(t, repoDir, "commit", "-m", "main advance")

	err := rebaseIfNeeded(repoDir, wtDir, "feature")
	if err != nil {
		t.Fatalf("expected successful rebase, got: %v", err)
	}

	// Verify that main is now ancestor of feature
	cmd := exec.Command("git", "merge-base", "--is-ancestor", "main", "feature")
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		t.Fatal("expected main to be ancestor of feature after rebase")
	}
}

func TestRebaseIfNeeded_ConflictAborts(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := initTestRepo(t, tmpDir)

	// Create worktree and modify README.md
	wtDir := filepath.Join(tmpDir, "wt-feature")
	gitRun(t, repoDir, "worktree", "add", wtDir, "-b", "feature")
	if err := os.WriteFile(filepath.Join(wtDir, "README.md"), []byte("feature version"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, wtDir, "add", ".")
	gitRun(t, wtDir, "commit", "-m", "feature modifies readme")

	// Advance main with a conflicting commit on the same file
	if err := os.WriteFile(filepath.Join(repoDir, "README.md"), []byte("main version"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, repoDir, "add", ".")
	gitRun(t, repoDir, "commit", "-m", "main modifies readme")

	err := rebaseIfNeeded(repoDir, wtDir, "feature")
	if err == nil {
		t.Fatal("expected error due to rebase conflict")
	}
	if !strings.Contains(err.Error(), "rebase failed") {
		t.Errorf("expected 'rebase failed' in error, got: %v", err)
	}

	// Verify rebase was aborted (branch should be clean)
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = wtDir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git status failed: %v", err)
	}
	if strings.TrimSpace(string(out)) != "" {
		t.Errorf("expected clean worktree after rebase abort, got: %s", out)
	}
}

// --- mergeIntoMain tests ---

func TestMergeIntoMain_Success(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := initTestRepo(t, tmpDir)

	// Create a branch with a commit
	wtDir := filepath.Join(tmpDir, "wt-feature")
	gitRun(t, repoDir, "worktree", "add", wtDir, "-b", "feature")
	if err := os.WriteFile(filepath.Join(wtDir, "feature.txt"), []byte("feature"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, wtDir, "add", ".")
	gitRun(t, wtDir, "commit", "-m", "feature work")

	err := mergeIntoMain(repoDir, "feature")
	if err != nil {
		t.Fatalf("expected merge success, got: %v", err)
	}

	// Verify we're on main
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(out)) != "main" {
		t.Errorf("expected to be on main, got %s", out)
	}

	// Verify the merge commit exists (--no-ff creates a merge commit)
	cmd = exec.Command("git", "log", "--oneline", "-1")
	cmd.Dir = repoDir
	out, err = cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "Merge branch") {
		t.Errorf("expected merge commit, got: %s", out)
	}

	// Verify feature file exists on main
	if _, err := os.Stat(filepath.Join(repoDir, "feature.txt")); os.IsNotExist(err) {
		t.Error("expected feature.txt to exist on main after merge")
	}
}

// --- cleanupWorktreeAndBranch tests ---

func TestCleanupWorktreeAndBranch_Success(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := initTestRepo(t, tmpDir)

	// Create worktree and merge it first
	wtDir := filepath.Join(tmpDir, "wt-cleanup")
	gitRun(t, repoDir, "worktree", "add", wtDir, "-b", "cleanup-branch")
	if err := os.WriteFile(filepath.Join(wtDir, "file.txt"), []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, wtDir, "add", ".")
	gitRun(t, wtDir, "commit", "-m", "add file")

	// Merge first so branch -d succeeds
	gitRun(t, repoDir, "merge", "--no-ff", "cleanup-branch", "-m", "merge")

	cleanupWorktreeAndBranch(repoDir, wtDir, "cleanup-branch")

	// Verify worktree is removed
	if _, err := os.Stat(wtDir); !os.IsNotExist(err) {
		t.Error("expected worktree directory to be removed")
	}

	// Verify branch is deleted
	cmd := exec.Command("git", "branch", "--list", "cleanup-branch")
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(out)) != "" {
		t.Error("expected branch to be deleted")
	}
}

func TestCleanupWorktreeAndBranch_FailureDoesNotPanic(t *testing.T) {
	// Cleanup with invalid paths should warn but not panic
	cleanupWorktreeAndBranch("/nonexistent", "/nonexistent/wt", "nonexistent-branch")
	// If we get here without panic, the test passes
}

// --- Full merge integration test ---

func TestMergeIntegration_FullCycle(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := initTestRepo(t, tmpDir)

	changeName := "int-test-change"
	allDoneTasks := []map[string]interface{}{
		{"id": "t1", "title": "Task 1", "status": "done"},
		{"id": "t2", "title": "Task 2", "status": "done", "blockers": []string{"t1"}},
	}
	wtDir := createWorktreeWithTasks(t, repoDir, changeName, allDoneTasks)

	// Add a feature file
	if err := os.WriteFile(filepath.Join(wtDir, "feature.txt"), []byte("feature content"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, wtDir, "add", ".")
	gitRun(t, wtDir, "commit", "-m", "feature implementation")

	// 1. Check tasks are all done
	if err := checkAllTasksDone(wtDir, changeName); err != nil {
		t.Fatalf("tasks check failed: %v", err)
	}

	// 2. Rebase (no rebase needed since main hasn't advanced)
	if err := rebaseIfNeeded(repoDir, wtDir, changeName); err != nil {
		t.Fatalf("rebase failed: %v", err)
	}

	// 3. Merge
	if err := mergeIntoMain(repoDir, changeName); err != nil {
		t.Fatalf("merge failed: %v", err)
	}

	// 4. Cleanup
	cleanupWorktreeAndBranch(repoDir, wtDir, changeName)

	// Verify: on main, feature file exists, worktree gone, branch gone
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repoDir
	out, _ := cmd.Output()
	if strings.TrimSpace(string(out)) != "main" {
		t.Errorf("expected to be on main, got %s", out)
	}

	if _, err := os.Stat(filepath.Join(repoDir, "feature.txt")); os.IsNotExist(err) {
		t.Error("expected feature.txt on main after merge")
	}

	if _, err := os.Stat(wtDir); !os.IsNotExist(err) {
		t.Error("expected worktree to be removed")
	}

	cmd = exec.Command("git", "branch", "--list", changeName)
	cmd.Dir = repoDir
	out, _ = cmd.Output()
	if strings.TrimSpace(string(out)) != "" {
		t.Error("expected branch to be deleted")
	}
}

func TestMergeIntegration_RebaseThenMerge(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := initTestRepo(t, tmpDir)

	changeName := "rebase-test"
	allDoneTasks := []map[string]interface{}{
		{"id": "t1", "title": "Task 1", "status": "done"},
	}
	wtDir := createWorktreeWithTasks(t, repoDir, changeName, allDoneTasks)

	// Add feature file in worktree
	if err := os.WriteFile(filepath.Join(wtDir, "feature.txt"), []byte("feature"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, wtDir, "add", ".")
	gitRun(t, wtDir, "commit", "-m", "feature work")

	// Advance main with a non-conflicting commit
	if err := os.WriteFile(filepath.Join(repoDir, "main-advance.txt"), []byte("main"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, repoDir, "add", ".")
	gitRun(t, repoDir, "commit", "-m", "main advance")

	// Rebase should succeed
	if err := rebaseIfNeeded(repoDir, wtDir, changeName); err != nil {
		t.Fatalf("rebase failed: %v", err)
	}

	// Merge should succeed
	if err := mergeIntoMain(repoDir, changeName); err != nil {
		t.Fatalf("merge failed: %v", err)
	}

	// Verify both files exist on main
	for _, f := range []string{"feature.txt", "main-advance.txt"} {
		if _, err := os.Stat(filepath.Join(repoDir, f)); os.IsNotExist(err) {
			t.Errorf("expected %s on main after merge", f)
		}
	}

	// Cleanup
	cleanupWorktreeAndBranch(repoDir, wtDir, changeName)
}
