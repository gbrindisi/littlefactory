//go:build integration

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// binaryPath is set by TestMain to the compiled littlefactory binary.
var binaryPath string

func TestMain(m *testing.M) {
	// Build binary to a temp directory
	tmpDir, err := os.MkdirTemp("", "lf-integration-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp dir: %v\n", err)
		os.Exit(1)
	}

	binaryPath = filepath.Join(tmpDir, "littlefactory")
	cmd := exec.Command("go", "build", "-tags=integration", "-o", binaryPath, "./")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build binary: %v\n", err)
		os.RemoveAll(tmpDir)
		os.Exit(1)
	}

	code := m.Run()
	os.RemoveAll(tmpDir)
	os.Exit(code)
}

// scaffoldRepo creates a git repo with a Factoryfile, change directory, and tasks.json.
// Returns the repo path.
func scaffoldRepo(t *testing.T, changeName string, taskList []map[string]interface{}) string {
	t.Helper()

	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Init git repo
	gitRunInteg(t, repoDir, "init")
	gitRunInteg(t, repoDir, "config", "user.email", "test@test.com")
	gitRunInteg(t, repoDir, "config", "user.name", "Test")
	gitRunInteg(t, repoDir, "checkout", "-b", "main")

	// Write Factoryfile
	factoryfile := `max_iterations: 10
timeout: 30
default_agent: echo

agents:
  echo:
    command: "echo done"
`
	if err := os.WriteFile(filepath.Join(repoDir, "Factoryfile"), []byte(factoryfile), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create change directory with tasks.json
	changeDir := filepath.Join(repoDir, ".littlefactory", "changes", changeName)
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

	// Initial commit
	gitRunInteg(t, repoDir, "add", ".")
	gitRunInteg(t, repoDir, "commit", "-m", "initial commit")

	return repoDir
}

// gitRunInteg executes a git command in the given directory and fails the test on error.
func gitRunInteg(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, out)
	}
}

// runBinary executes the littlefactory binary with args in the given directory.
// Returns stdout, stderr, and exit code.
func runBinary(t *testing.T, dir string, args ...string) (stdout string, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = dir

	var stdoutBuf, stderrBuf strings.Builder
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run binary: %v", err)
		}
	}

	return stdoutBuf.String(), stderrBuf.String(), exitCode
}

func defaultTasks() []map[string]interface{} {
	return []map[string]interface{}{
		{"id": "t1", "title": "Task One", "status": "todo"},
		{"id": "t2", "title": "Task Two", "status": "todo", "blockers": []string{"t1"}},
	}
}

// TestWorkflow_RunWithWorktree tests that run -c X -w creates a worktree,
// completes tasks in it, and exits 0.
func TestWorkflow_RunWithWorktree(t *testing.T) {
	changeName := "run-wt-test"
	repoDir := scaffoldRepo(t, changeName, defaultTasks())

	stdout, stderr, exitCode := runBinary(t, repoDir, "run", "-c", changeName, "-w")
	if exitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
	}

	// Verify worktree exists
	parentDir := filepath.Dir(repoDir)
	wtDir := filepath.Join(parentDir, changeName)
	if _, err := os.Stat(wtDir); os.IsNotExist(err) {
		t.Fatalf("expected worktree at %s to exist", wtDir)
	}

	// Verify tasks are done in the worktree
	tasksPath := filepath.Join(wtDir, ".littlefactory", "changes", changeName, "tasks.json")
	data, err := os.ReadFile(tasksPath)
	if err != nil {
		t.Fatalf("failed to read tasks.json in worktree: %v", err)
	}

	var tf struct {
		Tasks []struct {
			Status string `json:"status"`
		} `json:"tasks"`
	}
	if err := json.Unmarshal(data, &tf); err != nil {
		t.Fatal(err)
	}
	for i, task := range tf.Tasks {
		if task.Status != "done" {
			t.Errorf("task %d: expected status 'done', got %q", i, task.Status)
		}
	}
}

// TestWorkflow_RunReusesWorktree tests that a second run reuses an existing worktree.
func TestWorkflow_RunReusesWorktree(t *testing.T) {
	changeName := "reuse-wt-test"
	repoDir := scaffoldRepo(t, changeName, defaultTasks())

	// First run creates worktree and completes tasks
	_, _, exitCode := runBinary(t, repoDir, "run", "-c", changeName, "-w")
	if exitCode != 0 {
		t.Fatalf("first run failed with exit code %d", exitCode)
	}

	// Second run should reuse worktree
	stdout, stderr, exitCode := runBinary(t, repoDir, "run", "-c", changeName, "-w")
	if exitCode != 0 {
		t.Fatalf("second run failed with exit code %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
	}

	if !strings.Contains(stdout, "Reusing existing worktree") {
		t.Errorf("expected 'Reusing existing worktree' in stdout, got:\n%s", stdout)
	}
}

// TestWorkflow_StatusAll tests that status --all shows run state and [ready to merge].
func TestWorkflow_StatusAll(t *testing.T) {
	changeName := "status-all-test"
	repoDir := scaffoldRepo(t, changeName, defaultTasks())

	// Run to completion
	_, _, exitCode := runBinary(t, repoDir, "run", "-c", changeName, "-w")
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d", exitCode)
	}

	// Check status --all
	stdout, stderr, exitCode := runBinary(t, repoDir, "status", "--all")
	if exitCode != 0 {
		t.Fatalf("status --all failed with exit code %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
	}

	if !strings.Contains(stdout, changeName) {
		t.Errorf("expected change name %q in status output, got:\n%s", changeName, stdout)
	}
	if !strings.Contains(stdout, "[ready to merge]") {
		t.Errorf("expected '[ready to merge]' in status output, got:\n%s", stdout)
	}
}

// TestWorkflow_Verify tests that verify -c X runs in worktree context and exits 0.
func TestWorkflow_Verify(t *testing.T) {
	changeName := "verify-test"
	repoDir := scaffoldRepo(t, changeName, defaultTasks())

	// Run to create worktree
	_, _, exitCode := runBinary(t, repoDir, "run", "-c", changeName, "-w")
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d", exitCode)
	}

	// Verify
	stdout, stderr, exitCode := runBinary(t, repoDir, "verify", "-c", changeName)
	if exitCode != 0 {
		t.Fatalf("verify failed with exit code %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
	}

	if !strings.Contains(stdout, "Running verification in worktree") {
		t.Errorf("expected 'Running verification in worktree' in stdout, got:\n%s", stdout)
	}
}

// TestWorkflow_MergeFullCycle tests the full merge lifecycle: verify, merge, cleanup.
func TestWorkflow_MergeFullCycle(t *testing.T) {
	changeName := "merge-full-test"
	repoDir := scaffoldRepo(t, changeName, defaultTasks())

	// Run to completion
	_, _, exitCode := runBinary(t, repoDir, "run", "-c", changeName, "-w")
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d", exitCode)
	}

	// Merge
	stdout, stderr, exitCode := runBinary(t, repoDir, "merge", "-c", changeName)
	if exitCode != 0 {
		t.Fatalf("merge failed with exit code %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
	}

	// Verify worktree is removed
	parentDir := filepath.Dir(repoDir)
	wtDir := filepath.Join(parentDir, changeName)
	if _, err := os.Stat(wtDir); !os.IsNotExist(err) {
		t.Error("expected worktree directory to be removed")
	}

	// Verify branch is deleted
	cmd := exec.Command("git", "branch", "--list", changeName)
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(out)) != "" {
		t.Error("expected branch to be deleted")
	}

	// Verify we're on main
	cmd = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repoDir
	out, err = cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(out)) != "main" {
		t.Errorf("expected to be on main, got %s", strings.TrimSpace(string(out)))
	}
}

// TestWorkflow_MergeWithRebase tests that merge rebases when main has advanced.
func TestWorkflow_MergeWithRebase(t *testing.T) {
	changeName := "merge-rebase-test"
	repoDir := scaffoldRepo(t, changeName, defaultTasks())

	// Run to completion
	_, _, exitCode := runBinary(t, repoDir, "run", "-c", changeName, "-w")
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d", exitCode)
	}

	// Commit any unstaged files in the worktree (driver writes metadata)
	parentDir := filepath.Dir(repoDir)
	wtDir := filepath.Join(parentDir, changeName)
	gitRunInteg(t, wtDir, "add", ".")
	gitRunInteg(t, wtDir, "commit", "-m", "commit driver metadata")

	// Advance main with a non-conflicting commit
	if err := os.WriteFile(filepath.Join(repoDir, "main-advance.txt"), []byte("main"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRunInteg(t, repoDir, "add", ".")
	gitRunInteg(t, repoDir, "commit", "-m", "advance main")

	// Merge (should rebase first)
	stdout, stderr, exitCode := runBinary(t, repoDir, "merge", "-c", changeName)
	if exitCode != 0 {
		t.Fatalf("merge failed with exit code %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
	}

	// Verify both changesets exist on main
	if _, err := os.Stat(filepath.Join(repoDir, "main-advance.txt")); os.IsNotExist(err) {
		t.Error("expected main-advance.txt on main after merge")
	}

	// The worktree's changes are in .littlefactory/changes/ which should be on main after merge
	tasksPath := filepath.Join(repoDir, ".littlefactory", "changes", changeName, "tasks.json")
	data, err := os.ReadFile(tasksPath)
	if err != nil {
		t.Fatalf("expected tasks.json on main after merge: %v", err)
	}

	var tf struct {
		Tasks []struct {
			Status string `json:"status"`
		} `json:"tasks"`
	}
	if err := json.Unmarshal(data, &tf); err != nil {
		t.Fatal(err)
	}
	for i, task := range tf.Tasks {
		if task.Status != "done" {
			t.Errorf("task %d: expected status 'done' on main, got %q", i, task.Status)
		}
	}
}
