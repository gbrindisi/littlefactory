// Package driver provides the loop orchestrator and metadata tracking.
package driver

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/yourusername/littlefactory/internal/agent"
	"github.com/yourusername/littlefactory/internal/config"
	"github.com/yourusername/littlefactory/internal/tasks"
	"github.com/yourusername/littlefactory/internal/template"
)

// Driver is the core orchestrator that runs the autonomous agent loop.
// It coordinates task retrieval, template rendering, agent execution,
// and metadata tracking through sequential iterations.
type Driver struct {
	agent        agent.Agent
	taskSource   tasks.TaskSource
	config       *config.Config
	projectRoot  string
	metadata     *RunMetadata
	eventChan    chan interface{}
	changeName   string
	worktreePath string
}

// NewDriver creates a new Driver with the given dependencies.
// The eventChan parameter is used to send events to the TUI.
// If eventChan is nil, events are not emitted (for backward compatibility).
func NewDriver(ag agent.Agent, ts tasks.TaskSource, cfg *config.Config, projectRoot string, eventChan chan interface{}) *Driver {
	return &Driver{
		agent:       ag,
		taskSource:  ts,
		config:      cfg,
		projectRoot: projectRoot,
		eventChan:   eventChan,
	}
}

// SetChangeName sets the openspec change name for this run.
func (d *Driver) SetChangeName(name string) {
	d.changeName = name
}

// SetWorktreePath sets the worktree directory path for this run.
// When set, the agent will execute in this directory.
func (d *Driver) SetWorktreePath(path string) {
	d.worktreePath = path
}

// ChangeName returns the configured change name, if any.
func (d *Driver) ChangeName() string {
	return d.changeName
}

// WorktreePath returns the configured worktree path, if any.
func (d *Driver) WorktreePath() string {
	return d.worktreePath
}

// effectiveRoot returns the working root directory for file operations.
// Returns worktreePath if set, otherwise projectRoot.
// This ensures state files (progress.md, metadata) are written to the
// worktree when running in worktree mode.
func (d *Driver) effectiveRoot() string {
	if d.worktreePath != "" {
		return d.worktreePath
	}
	return d.projectRoot
}

// emit sends an event to the event channel if it exists.
// This is a helper to safely emit events without checking for nil everywhere.
func (d *Driver) emit(event interface{}) {
	if d.eventChan != nil {
		d.eventChan <- event
	}
}

// Run executes the main agent loop up to maxIterations.
// It returns the final RunStatus when the run completes.
// If WorktreePath is set, the process working directory is changed
// to the worktree directory before the loop starts and restored on return.
func (d *Driver) Run(ctx context.Context) RunStatus {
	// Switch to worktree directory if configured
	if d.worktreePath != "" {
		origDir, err := os.Getwd()
		if err != nil {
			return RunStatusFailed
		}
		if err := os.Chdir(d.worktreePath); err != nil {
			return RunStatusFailed
		}
		defer os.Chdir(origDir)
	}

	// Initialize run metadata
	d.metadata = &RunMetadata{
		RunID:         GenerateRunID(),
		StartedAt:     time.Now(),
		Status:        RunStatusRunning,
		MaxIterations: d.config.MaxIterations,
		Iterations:    []IterationMetadata{},
	}

	// Initialize progress file
	if err := InitProgressFile(d.effectiveRoot(), d.config); err != nil {
		// Log but continue - not fatal
	}

	// Save initial metadata
	_ = SaveMetadata(d.effectiveRoot(), d.config, d.metadata)

	// Get initial ready task count and emit start event
	readyTasks, err := d.taskSource.Ready()
	if err != nil {
		// Log error but continue - might recover
	}
	d.emit(RunStartedMsg{
		MaxIterations: d.config.MaxIterations,
		ReadyCount:    len(readyTasks),
	})

	// Emit initial task list for TUI
	allTasks, err := d.taskSource.List()
	if err == nil {
		d.emit(TasksRefreshedMsg{Tasks: allTasks})
	}

	// Exit early if no ready tasks at start
	if len(readyTasks) == 0 {
		d.metadata.Status = RunStatusCompleted
		d.FinalizeRun()
		return d.metadata.Status
	}

	// Main iteration loop
	for iterNum := 1; iterNum <= d.config.MaxIterations; iterNum++ {
		// Check for context cancellation (SIGINT)
		if ctx.Err() != nil {
			d.metadata.Status = RunStatusCancelled
			break
		}

		// Check if all tasks are complete
		if d.IsComplete() {
			d.metadata.Status = RunStatusCompleted
			break
		}

		// Run the iteration
		iterStatus := d.RunIteration(ctx, iterNum)

		// Save metadata after each iteration for crash resilience
		_ = SaveMetadata(d.effectiveRoot(), d.config, d.metadata)

		// Check for completion after successful iteration
		if iterStatus == IterationStatusCompleted && d.IsComplete() {
			d.metadata.Status = RunStatusCompleted
			break
		}
	}

	// Finalize the run
	d.FinalizeRun()

	// Emit run complete event
	d.emit(RunCompleteMsg{
		Status:   d.metadata.Status,
		Metadata: d.metadata,
	})

	return d.metadata.Status
}

// IsComplete checks if there are no more ready tasks.
// Returns true if the task source has no tasks available for work.
func (d *Driver) IsComplete() bool {
	readyTasks, err := d.taskSource.Ready()
	if err != nil {
		// On error, assume not complete to allow retry
		return false
	}
	return len(readyTasks) == 0
}

// RunIteration executes a single iteration with the given number.
// It retrieves a task, renders the prompt, executes the agent,
// and tracks the iteration metadata. Returns the iteration status.
func (d *Driver) RunIteration(ctx context.Context, iterNum int) IterationStatus {
	// Initialize iteration metadata
	iter := IterationMetadata{
		IterationNumber: iterNum,
		StartedAt:       time.Now(),
		Status:          IterationStatusRunning,
	}

	// Get the next ready task
	readyTasks, err := d.taskSource.Ready()
	if err != nil {
		iter.Status = IterationStatusFailed
		iter.ErrorMessage = ptr(err.Error())
		d.finalizeIteration(&iter)
		return iter.Status
	}

	if len(readyTasks) == 0 {
		// No tasks available - shouldn't happen if IsComplete() was checked
		iter.Status = IterationStatusFailed
		iter.ErrorMessage = ptr("no ready tasks available")
		d.finalizeIteration(&iter)
		return iter.Status
	}

	// Take the first ready task
	nextTask := readyTasks[0]
	iter.TaskID = ptr(nextTask.ID)
	iter.TaskTitle = ptr(nextTask.Title)

	// Emit iteration started event
	d.emit(IterationStartedMsg{
		Iteration: iterNum,
		TaskID:    nextTask.ID,
		TaskTitle: nextTask.Title,
	})

	// Claim the task (mark as in_progress)
	if err := d.taskSource.Claim(nextTask.ID); err != nil {
		iter.Status = IterationStatusFailed
		iter.ErrorMessage = ptr(err.Error())
		d.finalizeIteration(&iter)
		return iter.Status
	}

	// Get full task details
	taskDetails, err := d.taskSource.Show(nextTask.ID)
	if err != nil {
		iter.Status = IterationStatusFailed
		iter.ErrorMessage = ptr(err.Error())
		d.finalizeIteration(&iter)
		return iter.Status
	}

	// Load and render the template
	tmpl, err := template.Load(filepath.Join(d.effectiveRoot(), d.config.StateDir))
	if err != nil {
		iter.Status = IterationStatusFailed
		iter.ErrorMessage = ptr(err.Error())
		d.finalizeIteration(&iter)
		return iter.Status
	}

	prompt := template.Render(tmpl, taskDetails)

	// Create timeout context for this iteration
	iterCtx, cancel := context.WithTimeout(ctx, time.Duration(d.config.Timeout)*time.Second)
	defer cancel()

	// Create output writer that emits events
	var outputDest io.Writer
	if d.eventChan != nil {
		// TUI mode: only emit to event channel (stdout corrupts alternate screen)
		outputDest = newOutputWriter(d.eventChan)
	} else {
		// Non-TUI mode: write directly to stdout
		outputDest = os.Stdout
	}

	// Execute the agent
	// Output is streamed in real-time via io.Writer
	result, err := d.agent.Run(iterCtx, prompt, outputDest)

	// Record output metrics
	iter.OutputLines = ptr(result.OutputLines)
	iter.OutputBytes = ptr(result.OutputBytes)
	iter.ExitCode = ptr(result.ExitCode)

	// Determine status based on error and result
	// Check parent context first (SIGINT), then iteration context (timeout)
	switch {
	case ctx.Err() == context.Canceled:
		// Parent context cancelled (SIGINT) - iteration was interrupted
		iter.Status = IterationStatusFailed
		iter.ErrorMessage = ptr("interrupted by user")
	case iterCtx.Err() == context.DeadlineExceeded:
		iter.Status = IterationStatusTimeout
	case err != nil:
		iter.Status = IterationStatusFailed
		iter.ErrorMessage = ptr(err.Error())
	case result.ExitCode != 0:
		iter.Status = IterationStatusFailed
	default:
		iter.Status = IterationStatusCompleted
	}

	// Update task state based on iteration result
	if iter.Status == IterationStatusCompleted {
		// Mark task as done on success
		_ = d.taskSource.Close(nextTask.ID, "Completed")
	} else {
		// Reset task to todo on failure or timeout
		_ = d.taskSource.Reset(nextTask.ID)
	}

	// Emit iteration complete event
	d.emit(IterationCompleteMsg{Status: iter.Status})

	// Finalize iteration
	d.finalizeIteration(&iter)

	// Append to progress file
	statusStr := string(iter.Status)
	_ = AppendSessionToProgress(d.effectiveRoot(), d.config, iterNum, *iter.TaskID, statusStr)

	// Emit tasks refreshed event with updated task list
	allTasks, err := d.taskSource.List()
	if err == nil {
		d.emit(TasksRefreshedMsg{Tasks: allTasks})
	}

	return iter.Status
}

// finalizeIteration records the end time and duration, then adds
// the iteration to the run metadata.
func (d *Driver) finalizeIteration(iter *IterationMetadata) {
	now := time.Now()
	iter.EndedAt = &now
	duration := now.Sub(iter.StartedAt).Seconds()
	iter.DurationSeconds = &duration
	d.metadata.Iterations = append(d.metadata.Iterations, *iter)
}

// FinalizeRun completes the run by setting end time, calculating
// aggregate stats, and determining final status.
func (d *Driver) FinalizeRun() {
	now := time.Now()
	d.metadata.EndedAt = &now

	// Calculate total duration
	totalDuration := now.Sub(d.metadata.StartedAt).Seconds()
	d.metadata.TotalDurationSeconds = &totalDuration

	// Calculate aggregate statistics
	d.metadata.CalculateAggregateStats()

	// Set final status if not already set (cancelled or completed)
	if d.metadata.Status == RunStatusRunning {
		// If we reached max iterations without completing, mark as completed
		// (the run itself completed, even if tasks remain)
		if len(d.metadata.Iterations) >= d.config.MaxIterations {
			d.metadata.Status = RunStatusCompleted
		}
	}

	// Save final metadata
	_ = SaveMetadata(d.effectiveRoot(), d.config, d.metadata)
}

// Metadata returns the current run metadata.
// This is useful for inspecting the state of the run.
func (d *Driver) Metadata() *RunMetadata {
	return d.metadata
}
