package driver

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/gbrindisi/littlefactory/internal/agent"
	"github.com/gbrindisi/littlefactory/internal/config"
	"github.com/gbrindisi/littlefactory/internal/tasks"
)

func TestStatusLine_IterationStart(t *testing.T) {
	tmpDir := t.TempDir()
	var buf bytes.Buffer

	ag := &MockAgent{
		RunFunc: func(ctx context.Context, prompt string, output io.Writer) (agent.AgentResult, error) {
			return agent.AgentResult{ExitCode: 0}, nil
		},
	}
	ts := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-abc", Title: "Build feature"}},
	}
	cfg := &config.Config{MaxIterations: 5, Timeout: 60}

	d := NewDriver(ag, ts, cfg, tmpDir)
	d.SetStatusOutput(&buf)
	d.metadata = &RunMetadata{
		RunID:      "test-run",
		StartedAt:  time.Now(),
		Status:     RunStatusRunning,
		Iterations: []IterationMetadata{},
	}

	d.RunIteration(context.Background(), 2)

	output := buf.String()
	expected := "[2/5] Starting: Build feature (task-abc)"
	if !strings.Contains(output, expected) {
		t.Errorf("output %q does not contain %q", output, expected)
	}
}

func TestStatusLine_IterationCompleted(t *testing.T) {
	tmpDir := t.TempDir()
	var buf bytes.Buffer

	ag := &MockAgent{
		RunFunc: func(ctx context.Context, prompt string, output io.Writer) (agent.AgentResult, error) {
			return agent.AgentResult{ExitCode: 0}, nil
		},
	}
	ts := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-1", Title: "Test"}},
	}
	cfg := &config.Config{MaxIterations: 3, Timeout: 60}

	d := NewDriver(ag, ts, cfg, tmpDir)
	d.SetStatusOutput(&buf)
	d.metadata = &RunMetadata{
		RunID:      "test-run",
		StartedAt:  time.Now(),
		Status:     RunStatusRunning,
		Iterations: []IterationMetadata{},
	}

	d.RunIteration(context.Background(), 1)

	output := buf.String()
	expected := "[1/3] Completed"
	if !strings.Contains(output, expected) {
		t.Errorf("output %q does not contain %q", output, expected)
	}
}

func TestStatusLine_IterationFailed(t *testing.T) {
	tmpDir := t.TempDir()
	var buf bytes.Buffer

	ag := &MockAgent{
		RunFunc: func(ctx context.Context, prompt string, output io.Writer) (agent.AgentResult, error) {
			return agent.AgentResult{}, errors.New("something broke")
		},
	}
	ts := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-1", Title: "Test"}},
	}
	cfg := &config.Config{MaxIterations: 3, Timeout: 60}

	d := NewDriver(ag, ts, cfg, tmpDir)
	d.SetStatusOutput(&buf)
	d.metadata = &RunMetadata{
		RunID:      "test-run",
		StartedAt:  time.Now(),
		Status:     RunStatusRunning,
		Iterations: []IterationMetadata{},
	}

	d.RunIteration(context.Background(), 1)

	output := buf.String()
	expected := "[1/3] Failed: something broke"
	if !strings.Contains(output, expected) {
		t.Errorf("output %q does not contain %q", output, expected)
	}
}

func TestStatusLine_IterationFailedExitCode(t *testing.T) {
	tmpDir := t.TempDir()
	var buf bytes.Buffer

	ag := &MockAgent{
		RunFunc: func(ctx context.Context, prompt string, output io.Writer) (agent.AgentResult, error) {
			return agent.AgentResult{ExitCode: 1}, nil
		},
	}
	ts := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-1", Title: "Test"}},
	}
	cfg := &config.Config{MaxIterations: 3, Timeout: 60}

	d := NewDriver(ag, ts, cfg, tmpDir)
	d.SetStatusOutput(&buf)
	d.metadata = &RunMetadata{
		RunID:      "test-run",
		StartedAt:  time.Now(),
		Status:     RunStatusRunning,
		Iterations: []IterationMetadata{},
	}

	d.RunIteration(context.Background(), 1)

	output := buf.String()
	expected := "[1/3] Failed: exit code 1"
	if !strings.Contains(output, expected) {
		t.Errorf("output %q does not contain %q", output, expected)
	}
}

func TestStatusLine_IterationTimedOut(t *testing.T) {
	tmpDir := t.TempDir()
	var buf bytes.Buffer

	ag := &MockAgent{
		RunFunc: func(ctx context.Context, prompt string, output io.Writer) (agent.AgentResult, error) {
			<-ctx.Done()
			return agent.AgentResult{ExitCode: -1}, ctx.Err()
		},
	}
	ts := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-1", Title: "Test"}},
	}
	cfg := &config.Config{MaxIterations: 3, Timeout: 1}

	d := NewDriver(ag, ts, cfg, tmpDir)
	d.SetStatusOutput(&buf)
	d.metadata = &RunMetadata{
		RunID:      "test-run",
		StartedAt:  time.Now(),
		Status:     RunStatusRunning,
		Iterations: []IterationMetadata{},
	}

	d.RunIteration(context.Background(), 2)

	output := buf.String()
	expected := "[2/3] Timed out"
	if !strings.Contains(output, expected) {
		t.Errorf("output %q does not contain %q", output, expected)
	}
}

func TestStatusLine_RunSummaryCompleted(t *testing.T) {
	tmpDir := t.TempDir()
	var buf bytes.Buffer

	ag := &MockAgent{
		RunFunc: func(ctx context.Context, prompt string, output io.Writer) (agent.AgentResult, error) {
			return agent.AgentResult{ExitCode: 0}, nil
		},
	}
	ts := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-1", Title: "Test"}},
	}
	cfg := &config.Config{MaxIterations: 1, Timeout: 60}

	d := NewDriver(ag, ts, cfg, tmpDir)
	d.SetStatusOutput(&buf)

	d.Run(context.Background())

	output := buf.String()
	expected := "Run complete: completed (1/1 iterations)"
	if !strings.Contains(output, expected) {
		t.Errorf("output %q does not contain %q", output, expected)
	}
}

func TestStatusLine_RunSummaryCancelled(t *testing.T) {
	tmpDir := t.TempDir()
	var buf bytes.Buffer

	ag := &MockAgent{}
	ts := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-1", Title: "Test"}},
	}
	cfg := &config.Config{MaxIterations: 10, Timeout: 60}

	d := NewDriver(ag, ts, cfg, tmpDir)
	d.SetStatusOutput(&buf)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	d.Run(ctx)

	output := buf.String()
	if !strings.Contains(output, "Run cancelled") {
		t.Errorf("output %q does not contain %q", output, "Run cancelled")
	}
}

func TestStatusLine_RunSummaryNoTasks(t *testing.T) {
	tmpDir := t.TempDir()
	var buf bytes.Buffer

	ag := &MockAgent{}
	ts := &MockTaskSource{ReadyTasks: []tasks.Task{}}
	cfg := &config.Config{MaxIterations: 5, Timeout: 60}

	d := NewDriver(ag, ts, cfg, tmpDir)
	d.SetStatusOutput(&buf)

	d.Run(context.Background())

	output := buf.String()
	expected := "Run complete: completed (0/5 iterations)"
	if !strings.Contains(output, expected) {
		t.Errorf("output %q does not contain %q", output, expected)
	}
}
