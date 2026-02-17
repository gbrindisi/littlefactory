package driver

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gbrindisi/littlefactory/internal/agent"
	"github.com/gbrindisi/littlefactory/internal/config"
	"github.com/gbrindisi/littlefactory/internal/tasks"
)

// Note: MockAgent and MockTaskSource are defined in mocks_test.go

func TestNewDriver(t *testing.T) {
	ag := &MockAgent{}
	ts := &MockTaskSource{}
	cfg := &config.Config{MaxIterations: 5, Timeout: 60}

	d := NewDriver(ag, ts, cfg, "/test/project", nil)

	if d.agent != ag {
		t.Error("agent not set correctly")
	}
	if d.taskSource != ts {
		t.Error("taskSource not set correctly")
	}
	if d.config != cfg {
		t.Error("config not set correctly")
	}
	if d.projectRoot != "/test/project" {
		t.Errorf("projectRoot = %q, want %q", d.projectRoot, "/test/project")
	}
}

func TestIsComplete_NoTasks(t *testing.T) {
	ts := &MockTaskSource{ReadyTasks: []tasks.Task{}}
	d := &Driver{taskSource: ts}

	if !d.IsComplete() {
		t.Error("IsComplete() = false, want true when no tasks")
	}
}

func TestIsComplete_HasTasks(t *testing.T) {
	ts := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-1", Title: "Test"}},
	}
	d := &Driver{taskSource: ts}

	if d.IsComplete() {
		t.Error("IsComplete() = true, want false when tasks exist")
	}
}

func TestIsComplete_Error(t *testing.T) {
	ts := &MockTaskSource{ReadyErr: errors.New("error")}
	d := &Driver{taskSource: ts}

	// On error, should return false (assume not complete)
	if d.IsComplete() {
		t.Error("IsComplete() = true, want false on error")
	}
}

func TestRunIteration_Success(t *testing.T) {
	tmpDir := t.TempDir()

	ag := &MockAgent{
		RunFunc: func(ctx context.Context, prompt string, output io.Writer) (agent.AgentResult, error) {
			return agent.AgentResult{
				ExitCode:    0,
				OutputLines: 5,
				OutputBytes: 50,
			}, nil
		},
	}
	ts := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-1", Title: "Test Task"}},
	}
	cfg := &config.Config{MaxIterations: 1, Timeout: 60}

	d := NewDriver(ag, ts, cfg, tmpDir, nil)
	d.metadata = &RunMetadata{
		RunID:      "test-run",
		StartedAt:  time.Now(),
		Status:     RunStatusRunning,
		Iterations: []IterationMetadata{},
	}

	d.RunIteration(context.Background(), 1)

	if len(d.metadata.Iterations) != 1 {
		t.Fatalf("expected 1 iteration, got %d", len(d.metadata.Iterations))
	}

	iter := d.metadata.Iterations[0]
	if iter.Status != IterationStatusCompleted {
		t.Errorf("iteration status = %q, want %q", iter.Status, IterationStatusCompleted)
	}
	if iter.TaskID == nil || *iter.TaskID != "task-1" {
		t.Error("task ID not set correctly")
	}
	if iter.ExitCode == nil || *iter.ExitCode != 0 {
		t.Error("exit code not set correctly")
	}
}

func TestRunIteration_AgentFailure(t *testing.T) {
	tmpDir := t.TempDir()

	ag := &MockAgent{
		RunFunc: func(ctx context.Context, prompt string, output io.Writer) (agent.AgentResult, error) {
			return agent.AgentResult{ExitCode: 1}, nil
		},
	}
	ts := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-1", Title: "Test Task"}},
	}
	cfg := &config.Config{MaxIterations: 1, Timeout: 60}

	d := NewDriver(ag, ts, cfg, tmpDir, nil)
	d.metadata = &RunMetadata{
		RunID:      "test-run",
		StartedAt:  time.Now(),
		Status:     RunStatusRunning,
		Iterations: []IterationMetadata{},
	}

	d.RunIteration(context.Background(), 1)

	iter := d.metadata.Iterations[0]
	if iter.Status != IterationStatusFailed {
		t.Errorf("iteration status = %q, want %q", iter.Status, IterationStatusFailed)
	}
}

func TestRunIteration_NoReadyTasks(t *testing.T) {
	tmpDir := t.TempDir()

	ag := &MockAgent{}
	ts := &MockTaskSource{ReadyTasks: []tasks.Task{}}
	cfg := &config.Config{MaxIterations: 1, Timeout: 60}

	d := NewDriver(ag, ts, cfg, tmpDir, nil)
	d.metadata = &RunMetadata{
		RunID:      "test-run",
		StartedAt:  time.Now(),
		Status:     RunStatusRunning,
		Iterations: []IterationMetadata{},
	}

	d.RunIteration(context.Background(), 1)

	iter := d.metadata.Iterations[0]
	if iter.Status != IterationStatusFailed {
		t.Errorf("iteration status = %q, want %q", iter.Status, IterationStatusFailed)
	}
	if iter.ErrorMessage == nil || *iter.ErrorMessage != "no ready tasks available" {
		t.Error("expected 'no ready tasks available' error message")
	}
}

func TestRunIteration_TaskSourceError(t *testing.T) {
	tmpDir := t.TempDir()

	ag := &MockAgent{}
	ts := &MockTaskSource{ReadyErr: errors.New("task source error")}
	cfg := &config.Config{MaxIterations: 1, Timeout: 60}

	d := NewDriver(ag, ts, cfg, tmpDir, nil)
	d.metadata = &RunMetadata{
		RunID:      "test-run",
		StartedAt:  time.Now(),
		Status:     RunStatusRunning,
		Iterations: []IterationMetadata{},
	}

	d.RunIteration(context.Background(), 1)

	iter := d.metadata.Iterations[0]
	if iter.Status != IterationStatusFailed {
		t.Errorf("iteration status = %q, want %q", iter.Status, IterationStatusFailed)
	}
}

func TestRunIteration_Timeout(t *testing.T) {
	tmpDir := t.TempDir()

	ag := &MockAgent{
		RunFunc: func(ctx context.Context, prompt string, output io.Writer) (agent.AgentResult, error) {
			// Simulate timeout by waiting for context
			<-ctx.Done()
			return agent.AgentResult{ExitCode: -1}, ctx.Err()
		},
	}
	ts := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-1", Title: "Test Task"}},
	}
	cfg := &config.Config{MaxIterations: 1, Timeout: 1} // 1 second timeout

	d := NewDriver(ag, ts, cfg, tmpDir, nil)
	d.metadata = &RunMetadata{
		RunID:      "test-run",
		StartedAt:  time.Now(),
		Status:     RunStatusRunning,
		Iterations: []IterationMetadata{},
	}

	d.RunIteration(context.Background(), 1)

	iter := d.metadata.Iterations[0]
	if iter.Status != IterationStatusTimeout {
		t.Errorf("iteration status = %q, want %q", iter.Status, IterationStatusTimeout)
	}
}

func TestFinalizeRun(t *testing.T) {
	d := &Driver{
		projectRoot: t.TempDir(),
		config:      &config.Config{MaxIterations: 10, Timeout: 60},
	}

	startTime := time.Now().Add(-10 * time.Second)
	duration := 5.0
	d.metadata = &RunMetadata{
		RunID:     "test-run",
		StartedAt: startTime,
		Status:    RunStatusRunning,
		Iterations: []IterationMetadata{
			{
				IterationNumber: 1,
				Status:          IterationStatusCompleted,
				DurationSeconds: &duration,
			},
			{
				IterationNumber: 2,
				Status:          IterationStatusFailed,
				DurationSeconds: &duration,
			},
		},
	}

	d.FinalizeRun()

	if d.metadata.EndedAt == nil {
		t.Error("EndedAt not set")
	}
	if d.metadata.TotalDurationSeconds == nil {
		t.Error("TotalDurationSeconds not set")
	}
	if d.metadata.TotalIterations != 2 {
		t.Errorf("TotalIterations = %d, want 2", d.metadata.TotalIterations)
	}
	if d.metadata.SuccessfulIterations != 1 {
		t.Errorf("SuccessfulIterations = %d, want 1", d.metadata.SuccessfulIterations)
	}
	if d.metadata.FailedIterations != 1 {
		t.Errorf("FailedIterations = %d, want 1", d.metadata.FailedIterations)
	}
}

func TestRun_CompletesWhenNoTasks(t *testing.T) {
	tmpDir := t.TempDir()

	ag := &MockAgent{}
	ts := &MockTaskSource{ReadyTasks: []tasks.Task{}} // No tasks
	cfg := &config.Config{MaxIterations: 10, Timeout: 60}

	d := NewDriver(ag, ts, cfg, tmpDir, nil)

	status := d.Run(context.Background())

	if status != RunStatusCompleted {
		t.Errorf("Run() = %q, want %q", status, RunStatusCompleted)
	}
	if len(d.metadata.Iterations) != 0 {
		t.Errorf("expected 0 iterations, got %d", len(d.metadata.Iterations))
	}
}

func TestRun_ContextCancellation(t *testing.T) {
	tmpDir := t.TempDir()

	ag := &MockAgent{}
	ts := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-1", Title: "Test"}},
	}
	cfg := &config.Config{MaxIterations: 10, Timeout: 60}

	d := NewDriver(ag, ts, cfg, tmpDir, nil)

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	status := d.Run(ctx)

	if status != RunStatusCancelled {
		t.Errorf("Run() = %q, want %q", status, RunStatusCancelled)
	}
}

func TestMetadata_ReturnsMetadata(t *testing.T) {
	d := &Driver{
		metadata: &RunMetadata{RunID: "test-123"},
	}

	meta := d.Metadata()

	if meta == nil {
		t.Fatal("Metadata() returned nil")
	}
	if meta.RunID != "test-123" {
		t.Errorf("RunID = %q, want %q", meta.RunID, "test-123")
	}
}

// TestRunIteration_AgentError tests error returned by agent (distinct from non-zero exit code)
func TestRunIteration_AgentError(t *testing.T) {
	tmpDir := t.TempDir()

	ag := &MockAgent{
		RunFunc: func(ctx context.Context, prompt string, output io.Writer) (agent.AgentResult, error) {
			return agent.AgentResult{}, errors.New("agent execution error")
		},
	}
	ts := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-1", Title: "Test Task"}},
	}
	cfg := &config.Config{MaxIterations: 1, Timeout: 60}

	d := NewDriver(ag, ts, cfg, tmpDir, nil)
	d.metadata = &RunMetadata{
		RunID:      "test-run",
		StartedAt:  time.Now(),
		Status:     RunStatusRunning,
		Iterations: []IterationMetadata{},
	}

	d.RunIteration(context.Background(), 1)

	iter := d.metadata.Iterations[0]
	if iter.Status != IterationStatusFailed {
		t.Errorf("iteration status = %q, want %q", iter.Status, IterationStatusFailed)
	}
	if iter.ErrorMessage == nil || *iter.ErrorMessage != "agent execution error" {
		t.Errorf("expected 'agent execution error' error message, got %v", iter.ErrorMessage)
	}
}

// TestRun_KeepGoingOnFailure verifies that the driver continues iterating after a failure
func TestRun_KeepGoingOnFailure(t *testing.T) {
	tmpDir := t.TempDir()

	callCount := 0
	ag := &MockAgent{
		RunFunc: func(ctx context.Context, prompt string, output io.Writer) (agent.AgentResult, error) {
			callCount++
			if callCount == 1 {
				// First iteration fails
				return agent.AgentResult{ExitCode: 1}, nil
			}
			// Second iteration succeeds
			return agent.AgentResult{ExitCode: 0}, nil
		},
	}

	// TaskSource always returns tasks (we're testing keep-going, not task completion)
	ts := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-1", Title: "Test Task"}},
	}

	cfg := &config.Config{MaxIterations: 2, Timeout: 60}

	d := NewDriver(ag, ts, cfg, tmpDir, nil)

	status := d.Run(context.Background())

	// Even with failure in first iteration, driver should run both iterations
	// Final status is "completed" because the run itself completed (reached max iterations)
	// This is the "keep-going" strategy - failures don't abort the run
	if d.metadata.TotalIterations != 2 {
		t.Errorf("TotalIterations = %d, want 2 (keep-going should continue after failure)", d.metadata.TotalIterations)
	}
	if d.metadata.FailedIterations != 1 {
		t.Errorf("FailedIterations = %d, want 1", d.metadata.FailedIterations)
	}
	if d.metadata.SuccessfulIterations != 1 {
		t.Errorf("SuccessfulIterations = %d, want 1", d.metadata.SuccessfulIterations)
	}
	// Status is "completed" because the run itself completed (max iterations reached)
	// The keep-going behavior means failures don't change the run status
	if status != RunStatusCompleted {
		t.Errorf("Run() = %q, want %q", status, RunStatusCompleted)
	}
}

func TestDriver_SetChangeName(t *testing.T) {
	d := &Driver{}
	d.SetChangeName("feature-a")
	if d.ChangeName() != "feature-a" {
		t.Errorf("ChangeName() = %q, want %q", d.ChangeName(), "feature-a")
	}
}

func TestDriver_SetWorktreePath(t *testing.T) {
	d := &Driver{}
	d.SetWorktreePath("/tmp/worktrees/feature-a")
	if d.WorktreePath() != "/tmp/worktrees/feature-a" {
		t.Errorf("WorktreePath() = %q, want %q", d.WorktreePath(), "/tmp/worktrees/feature-a")
	}
}

func TestDriver_DefaultChangeNameEmpty(t *testing.T) {
	d := &Driver{}
	if d.ChangeName() != "" {
		t.Errorf("default ChangeName() = %q, want empty", d.ChangeName())
	}
}

func TestDriver_DefaultWorktreePathEmpty(t *testing.T) {
	d := &Driver{}
	if d.WorktreePath() != "" {
		t.Errorf("default WorktreePath() = %q, want empty", d.WorktreePath())
	}
}

func TestRun_WorktreePathSwitchesDirectory(t *testing.T) {
	worktreeDir := t.TempDir()
	projectDir := t.TempDir()

	// Save original working directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	var capturedDir string
	ag := &MockAgent{
		RunFunc: func(ctx context.Context, prompt string, output io.Writer) (agent.AgentResult, error) {
			// Capture the working directory during agent execution
			capturedDir, _ = os.Getwd()
			return agent.AgentResult{ExitCode: 0}, nil
		},
	}
	ts := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-1", Title: "Test"}},
	}
	cfg := &config.Config{MaxIterations: 1, Timeout: 60}

	d := NewDriver(ag, ts, cfg, projectDir, nil)
	d.SetWorktreePath(worktreeDir)

	d.Run(context.Background())

	// Resolve symlinks for comparison (macOS /var -> /private/var)
	resolvedWorktree, _ := filepath.EvalSymlinks(worktreeDir)
	resolvedCaptured, _ := filepath.EvalSymlinks(capturedDir)

	if resolvedCaptured != resolvedWorktree {
		t.Errorf("agent ran in %q, want %q", resolvedCaptured, resolvedWorktree)
	}

	// Verify working directory is restored after Run
	currentDir, _ := os.Getwd()
	resolvedCurrent, _ := filepath.EvalSymlinks(currentDir)
	resolvedOrig, _ := filepath.EvalSymlinks(origDir)
	if resolvedCurrent != resolvedOrig {
		t.Errorf("working directory not restored: got %q, want %q", resolvedCurrent, resolvedOrig)
	}
}

func TestRun_WorktreePathInvalidFails(t *testing.T) {
	projectDir := t.TempDir()

	ag := &MockAgent{}
	ts := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-1", Title: "Test"}},
	}
	cfg := &config.Config{MaxIterations: 1, Timeout: 60}

	d := NewDriver(ag, ts, cfg, projectDir, nil)
	d.SetWorktreePath("/nonexistent/path/that/does/not/exist")

	status := d.Run(context.Background())

	if status != RunStatusFailed {
		t.Errorf("Run() = %q, want %q for invalid worktree path", status, RunStatusFailed)
	}
}

func TestRun_NoWorktreePathPreservesDirectory(t *testing.T) {
	projectDir := t.TempDir()

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	ag := &MockAgent{}
	ts := &MockTaskSource{ReadyTasks: []tasks.Task{}} // No tasks, exits early
	cfg := &config.Config{MaxIterations: 1, Timeout: 60}

	d := NewDriver(ag, ts, cfg, projectDir, nil)
	// No worktree path set

	d.Run(context.Background())

	currentDir, _ := os.Getwd()
	resolvedCurrent, _ := filepath.EvalSymlinks(currentDir)
	resolvedOrig, _ := filepath.EvalSymlinks(origDir)
	if resolvedCurrent != resolvedOrig {
		t.Errorf("working directory changed without worktree: got %q, want %q", resolvedCurrent, resolvedOrig)
	}
}
