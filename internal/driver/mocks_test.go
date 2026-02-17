package driver

import (
	"context"
	"io"
	"testing"

	"github.com/yourusername/littlefactory/internal/agent"
	"github.com/yourusername/littlefactory/internal/tasks"
)

// MockAgent implements the Agent interface for testing.
// It allows configuring the behavior of Run() via the RunFunc field.
type MockAgent struct {
	RunFunc func(ctx context.Context, prompt string, output io.Writer) (agent.AgentResult, error)
}

// Run implements the Agent interface.
func (m *MockAgent) Run(ctx context.Context, prompt string, output io.Writer) (agent.AgentResult, error) {
	if m.RunFunc != nil {
		return m.RunFunc(ctx, prompt, output)
	}
	return agent.AgentResult{ExitCode: 0, OutputLines: 10, OutputBytes: 100}, nil
}

// MockTaskSource implements the TaskSource interface for testing.
// It allows configuring behavior via exported fields.
type MockTaskSource struct {
	ReadyTasks []tasks.Task
	ReadyErr   error
	ListTasks  []tasks.Task
	ListErr    error
	ShowTask   *tasks.Task
	ShowErr    error
	ClaimErr   error
	CloseErr   error
	ResetErr   error
	ReadyCount int // tracks how many times Ready() was called
	ListCount  int // tracks how many times List() was called
	ClaimedIDs []string
	ClosedIDs  []string
	ResetIDs   []string
}

// Ready implements the TaskSource interface.
func (m *MockTaskSource) Ready() ([]tasks.Task, error) {
	m.ReadyCount++
	if m.ReadyErr != nil {
		return nil, m.ReadyErr
	}
	return m.ReadyTasks, nil
}

// List implements the TaskSource interface.
func (m *MockTaskSource) List() ([]tasks.Task, error) {
	m.ListCount++
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	return m.ListTasks, nil
}

// Show implements the TaskSource interface.
func (m *MockTaskSource) Show(id string) (*tasks.Task, error) {
	if m.ShowErr != nil {
		return nil, m.ShowErr
	}
	if m.ShowTask != nil {
		return m.ShowTask, nil
	}
	// Return a default task based on ID
	return &tasks.Task{
		ID:          id,
		Title:       "Test Task",
		Description: "Test Description",
	}, nil
}

// Claim implements the TaskSource interface.
func (m *MockTaskSource) Claim(id string) error {
	m.ClaimedIDs = append(m.ClaimedIDs, id)
	return m.ClaimErr
}

// Close implements the TaskSource interface.
func (m *MockTaskSource) Close(id, reason string) error {
	m.ClosedIDs = append(m.ClosedIDs, id)
	return m.CloseErr
}

// Reset implements the TaskSource interface.
func (m *MockTaskSource) Reset(id string) error {
	m.ResetIDs = append(m.ResetIDs, id)
	return m.ResetErr
}

// Compile-time interface verification
var _ agent.Agent = (*MockAgent)(nil)
var _ tasks.TaskSource = (*MockTaskSource)(nil)

// TestMockAgentImplementsInterface verifies MockAgent satisfies Agent interface.
func TestMockAgentImplementsInterface(t *testing.T) {
	var _ agent.Agent = (*MockAgent)(nil)
}

// TestMockTaskSourceImplementsInterface verifies MockTaskSource satisfies TaskSource interface.
func TestMockTaskSourceImplementsInterface(t *testing.T) {
	var _ tasks.TaskSource = (*MockTaskSource)(nil)
}

// TestMockAgentDefaultBehavior verifies default behavior when RunFunc is nil.
func TestMockAgentDefaultBehavior(t *testing.T) {
	mock := &MockAgent{}
	result, err := mock.Run(context.Background(), "test prompt", io.Discard)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", result.ExitCode)
	}
	if result.OutputLines != 10 {
		t.Errorf("OutputLines = %d, want 10", result.OutputLines)
	}
	if result.OutputBytes != 100 {
		t.Errorf("OutputBytes = %d, want 100", result.OutputBytes)
	}
}

// TestMockAgentCustomBehavior verifies custom RunFunc is called.
func TestMockAgentCustomBehavior(t *testing.T) {
	called := false
	mock := &MockAgent{
		RunFunc: func(ctx context.Context, prompt string, output io.Writer) (agent.AgentResult, error) {
			called = true
			if prompt != "test prompt" {
				t.Errorf("prompt = %q, want %q", prompt, "test prompt")
			}
			return agent.AgentResult{ExitCode: 42}, nil
		},
	}

	result, _ := mock.Run(context.Background(), "test prompt", io.Discard)

	if !called {
		t.Error("RunFunc was not called")
	}
	if result.ExitCode != 42 {
		t.Errorf("ExitCode = %d, want 42", result.ExitCode)
	}
}

// TestMockTaskSourceReadyBehavior verifies Ready() tracks call count.
func TestMockTaskSourceReadyBehavior(t *testing.T) {
	mock := &MockTaskSource{
		ReadyTasks: []tasks.Task{{ID: "task-1", Title: "Test"}},
	}

	if mock.ReadyCount != 0 {
		t.Errorf("ReadyCount = %d, want 0", mock.ReadyCount)
	}

	result, err := mock.Ready()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("len(result) = %d, want 1", len(result))
	}
	if mock.ReadyCount != 1 {
		t.Errorf("ReadyCount = %d, want 1", mock.ReadyCount)
	}
}

// TestMockTaskSourceShowBehavior verifies Show() returns configured task.
func TestMockTaskSourceShowBehavior(t *testing.T) {
	mock := &MockTaskSource{
		ShowTask: &tasks.Task{ID: "custom-id", Title: "Custom Task"},
	}

	task, err := mock.Show("any-id")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if task.ID != "custom-id" {
		t.Errorf("task.ID = %q, want %q", task.ID, "custom-id")
	}
}

// TestMockTaskSourceShowDefaultBehavior verifies Show() returns default task when ShowTask is nil.
func TestMockTaskSourceShowDefaultBehavior(t *testing.T) {
	mock := &MockTaskSource{}

	task, err := mock.Show("requested-id")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if task.ID != "requested-id" {
		t.Errorf("task.ID = %q, want %q", task.ID, "requested-id")
	}
}

// TestMockTaskSourceClaimBehavior verifies Claim() tracks claimed IDs.
func TestMockTaskSourceClaimBehavior(t *testing.T) {
	mock := &MockTaskSource{}

	mock.Claim("task-1")
	mock.Claim("task-2")

	if len(mock.ClaimedIDs) != 2 {
		t.Errorf("len(ClaimedIDs) = %d, want 2", len(mock.ClaimedIDs))
	}
	if mock.ClaimedIDs[0] != "task-1" {
		t.Errorf("ClaimedIDs[0] = %q, want %q", mock.ClaimedIDs[0], "task-1")
	}
	if mock.ClaimedIDs[1] != "task-2" {
		t.Errorf("ClaimedIDs[1] = %q, want %q", mock.ClaimedIDs[1], "task-2")
	}
}

// TestMockTaskSourceCloseBehavior verifies Close() tracks closed IDs.
func TestMockTaskSourceCloseBehavior(t *testing.T) {
	mock := &MockTaskSource{}

	mock.Close("task-1", "completed")
	mock.Close("task-2", "wont-fix")

	if len(mock.ClosedIDs) != 2 {
		t.Errorf("len(ClosedIDs) = %d, want 2", len(mock.ClosedIDs))
	}
	if mock.ClosedIDs[0] != "task-1" {
		t.Errorf("ClosedIDs[0] = %q, want %q", mock.ClosedIDs[0], "task-1")
	}
	if mock.ClosedIDs[1] != "task-2" {
		t.Errorf("ClosedIDs[1] = %q, want %q", mock.ClosedIDs[1], "task-2")
	}
}

// TestMockTaskSourceResetBehavior verifies Reset() tracks reset IDs.
func TestMockTaskSourceResetBehavior(t *testing.T) {
	mock := &MockTaskSource{}

	mock.Reset("task-1")
	mock.Reset("task-2")

	if len(mock.ResetIDs) != 2 {
		t.Errorf("len(ResetIDs) = %d, want 2", len(mock.ResetIDs))
	}
	if mock.ResetIDs[0] != "task-1" {
		t.Errorf("ResetIDs[0] = %q, want %q", mock.ResetIDs[0], "task-1")
	}
	if mock.ResetIDs[1] != "task-2" {
		t.Errorf("ResetIDs[1] = %q, want %q", mock.ResetIDs[1], "task-2")
	}
}
