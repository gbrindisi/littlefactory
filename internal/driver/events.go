// Package driver provides the loop orchestrator and metadata tracking.
package driver

import "github.com/gbrindisi/littlefactory/internal/tasks"

// RunStartedMsg is sent when the driver run begins.
// It contains information about the planned run.
type RunStartedMsg struct {
	MaxIterations int
	ReadyCount    int
}

// IterationStartedMsg is sent when a new iteration begins.
// It contains information about the task being executed.
type IterationStartedMsg struct {
	Iteration int
	TaskID    string
	TaskTitle string
}

// OutputMsg contains a chunk of output data from the agent.
// The TUI receives these messages in real-time as the agent produces output.
type OutputMsg struct {
	Data []byte
}

// IterationCompleteMsg is sent when an iteration finishes.
// It contains the final status of the iteration.
type IterationCompleteMsg struct {
	Status IterationStatus
}

// TasksRefreshedMsg is sent after each iteration with the updated task list.
// The TUI uses this to refresh the task panel.
type TasksRefreshedMsg struct {
	Tasks []tasks.Task
}

// RunCompleteMsg is sent when the entire run finishes.
// It contains the final status and metadata for display.
type RunCompleteMsg struct {
	Status   RunStatus
	Metadata *RunMetadata
}

// outputWriter is an io.Writer that emits OutputMsg events.
// It's used to stream agent output to the TUI in real-time.
type outputWriter struct {
	eventChan chan interface{}
}

// Write implements io.Writer by emitting OutputMsg events.
func (w *outputWriter) Write(p []byte) (n int, err error) {
	if w.eventChan != nil {
		// Make a copy of the data to avoid races
		data := make([]byte, len(p))
		copy(data, p)
		w.eventChan <- OutputMsg{Data: data}
	}
	return len(p), nil
}

// newOutputWriter creates a new outputWriter that emits to the given channel.
func newOutputWriter(eventChan chan interface{}) *outputWriter {
	return &outputWriter{eventChan: eventChan}
}
