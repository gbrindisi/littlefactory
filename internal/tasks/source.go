// Package tasks defines the task source interface and related types.
package tasks

// Task represents a single task from the task source.
type Task struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Status      string   `json:"status"` // Valid values: "todo", "in_progress", "done"
	Labels      []string `json:"labels,omitempty"`
	Blockers    []string `json:"blockers,omitempty"`
}

// TaskSource defines the interface for task management operations.
// Implementations can integrate with JSON files, GitHub issues, or other task systems.
type TaskSource interface {
	// Ready returns all tasks with no blockers that are available for work.
	Ready() ([]Task, error)

	// List returns all tasks with their current status.
	// This includes closed, blocked, pending, and ready tasks.
	List() ([]Task, error)

	// Show retrieves full details for a task by ID.
	Show(id string) (*Task, error)

	// Claim marks a task as in progress.
	Claim(id string) error

	// Close marks a task as complete with the given reason.
	Close(id, reason string) error

	// Reset marks a task as todo (available for work).
	Reset(id string) error
}
