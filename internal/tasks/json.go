package tasks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gbrindisi/littlefactory/internal/config"
)

// JSONTaskSource implements TaskSource using a local JSON file.
type JSONTaskSource struct {
	config      *config.Config
	projectRoot string
	tasksPath   string
}

// tasksFile represents the structure of the tasks.json file.
type tasksFile struct {
	Tasks []Task `json:"tasks"`
}

// validStatuses are the allowed values for a task's Status field.
var validStatuses = map[string]bool{
	"todo":        true,
	"in_progress": true,
	"done":        true,
}

// ValidateTasks validates the structure and content of a task list.
// It returns a multi-error with all validation failures, or nil if valid.
// filePath is used only for error messages.
func ValidateTasks(taskList []Task, filePath string) error {
	var errs []string

	// Track seen IDs for uniqueness check
	seenIDs := make(map[string]int) // id -> first index

	for i, task := range taskList {
		// Required field: id
		if task.ID == "" {
			errs = append(errs, fmt.Sprintf("task at index %d: missing required field \"id\"", i))
			continue // skip further checks for this task without an ID
		}

		// Unique ID check
		if firstIdx, exists := seenIDs[task.ID]; exists {
			errs = append(errs, fmt.Sprintf("task %q (index %d): duplicate id (first seen at index %d)", task.ID, i, firstIdx))
		} else {
			seenIDs[task.ID] = i
		}

		// Required field: title
		if task.Title == "" {
			errs = append(errs, fmt.Sprintf("task %q (index %d): missing required field \"title\"", task.ID, i))
		}

		// Required field: status
		if task.Status == "" {
			errs = append(errs, fmt.Sprintf("task %q (index %d): missing required field \"status\"", task.ID, i))
		} else if !validStatuses[task.Status] {
			errs = append(errs, fmt.Sprintf("task %q (index %d): invalid status %q (must be: todo, in_progress, done)", task.ID, i, task.Status))
		}
	}

	// Blocker chain validation (only if we have tasks and no ID issues that would break chain walking)
	if len(taskList) > 0 {
		errs = append(errs, validateBlockerChain(taskList, seenIDs)...)
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("error loading tasks from %s:\n  %s", filePath, strings.Join(errs, "\n  "))
}

// validateBlockerChain validates that tasks form a strict linear sequence via blockers.
func validateBlockerChain(taskList []Task, idIndex map[string]int) []string {
	var errs []string

	// Find root tasks (empty blockers)
	var roots []string
	for _, task := range taskList {
		if task.ID == "" {
			continue
		}
		if len(task.Blockers) == 0 {
			roots = append(roots, task.ID)
		}
	}

	// Validate exactly one root
	if len(roots) == 0 {
		errs = append(errs, "invalid sequence: no root task (all tasks have blockers)")
		return errs
	}
	if len(roots) > 1 {
		errs = append(errs, fmt.Sprintf("invalid sequence: multiple root tasks (tasks with no blockers): %s", formatIDList(roots)))
		return errs
	}

	// Validate single blocker per non-root task and blocker references
	blockerTargets := make(map[string]string) // blocker -> task that blocks on it
	for _, task := range taskList {
		if task.ID == "" {
			continue
		}
		if len(task.Blockers) > 1 {
			errs = append(errs, fmt.Sprintf("invalid sequence: task %q has %d blockers, expected 0 or 1", task.ID, len(task.Blockers)))
			continue
		}
		if len(task.Blockers) == 1 {
			blockerID := task.Blockers[0]
			if _, exists := idIndex[blockerID]; !exists {
				errs = append(errs, fmt.Sprintf("task %q (index %d): blocker %q does not exist", task.ID, idIndex[task.ID], blockerID))
			}
			blockerTargets[blockerID] = task.ID
		}
	}

	// Walk chain from root and check all tasks are reachable
	if len(errs) == 0 {
		visited := make(map[string]bool)
		current := roots[0]
		for current != "" {
			if visited[current] {
				break
			}
			visited[current] = true
			next, ok := blockerTargets[current]
			if !ok {
				break
			}
			current = next
		}

		var orphaned []string
		for _, task := range taskList {
			if task.ID != "" && !visited[task.ID] {
				orphaned = append(orphaned, task.ID)
			}
		}
		if len(orphaned) > 0 {
			errs = append(errs, fmt.Sprintf("invalid sequence: tasks not reachable from root: %s", formatIDList(orphaned)))
		}
	}

	return errs
}

// formatIDList formats a list of IDs as quoted, comma-separated strings.
func formatIDList(ids []string) string {
	quoted := make([]string, len(ids))
	for i, id := range ids {
		quoted[i] = fmt.Sprintf("%q", id)
	}
	return strings.Join(quoted, ", ")
}

// NewJSONTaskSource creates a new JSONTaskSource.
// projectRoot is the root directory of the project.
// cfg is the configuration containing the state directory path.
// It reads and validates the tasks file on creation.
func NewJSONTaskSource(projectRoot string, cfg *config.Config) (*JSONTaskSource, error) {
	tasksPath := filepath.Join(projectRoot, cfg.StateDir, "tasks.json")
	j := &JSONTaskSource{
		config:      cfg,
		projectRoot: projectRoot,
		tasksPath:   tasksPath,
	}

	// Read and validate tasks at load time
	taskList, err := j.readTasks()
	if err != nil {
		return nil, err
	}

	if err := ValidateTasks(taskList, tasksPath); err != nil {
		return nil, err
	}

	return j, nil
}

// NewJSONTaskSourceWithPath creates a new JSONTaskSource using an explicit file path.
// This is used when the task source comes from a non-default location (e.g., an openspec change).
// It checks that the file exists and validates its content on creation.
func NewJSONTaskSourceWithPath(tasksPath string) (*JSONTaskSource, error) {
	// Check file existence
	if _, err := os.Stat(tasksPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("tasks file not found: %s", tasksPath)
	}

	j := &JSONTaskSource{
		tasksPath: tasksPath,
	}

	// Read and validate tasks at load time
	taskList, err := j.readTasks()
	if err != nil {
		return nil, err
	}

	if err := ValidateTasks(taskList, tasksPath); err != nil {
		return nil, err
	}

	return j, nil
}

// readTasks reads and parses the tasks.json file.
func (j *JSONTaskSource) readTasks() ([]Task, error) {
	data, err := os.ReadFile(j.tasksPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If file doesn't exist, return empty list
			return []Task{}, nil
		}
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}

	var tf tasksFile
	if err := json.Unmarshal(data, &tf); err != nil {
		return nil, fmt.Errorf("failed to parse tasks file: %w", err)
	}

	return tf.Tasks, nil
}

// writeTasks writes the tasks to the tasks.json file.
// Creates the state directory if it doesn't exist.
func (j *JSONTaskSource) writeTasks(tasks []Task) error {
	// Ensure directory exists
	dir := filepath.Dir(j.tasksPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create tasks directory: %w", err)
	}

	tf := tasksFile{Tasks: tasks}
	data, err := json.MarshalIndent(tf, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	if err := os.WriteFile(j.tasksPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write tasks file: %w", err)
	}

	return nil
}

// Ready returns all tasks with status "todo".
// Returns the first task with status "todo" in array order.
func (j *JSONTaskSource) Ready() ([]Task, error) {
	tasks, err := j.readTasks()
	if err != nil {
		return nil, err
	}

	// Find first task with status "todo"
	for _, task := range tasks {
		if task.Status == "todo" {
			return []Task{task}, nil
		}
	}

	// No ready tasks
	return []Task{}, nil
}

// List returns all tasks from the JSON file.
func (j *JSONTaskSource) List() ([]Task, error) {
	return j.readTasks()
}

// Show retrieves full details for a task by ID.
func (j *JSONTaskSource) Show(id string) (*Task, error) {
	tasks, err := j.readTasks()
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		if task.ID == id {
			return &task, nil
		}
	}

	return nil, fmt.Errorf("task %s not found", id)
}

// Claim sets a task's status to "in_progress" and persists to file.
func (j *JSONTaskSource) Claim(id string) error {
	tasks, err := j.readTasks()
	if err != nil {
		return err
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = "in_progress"
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task %s not found", id)
	}

	return j.writeTasks(tasks)
}

// Close marks a task as complete (status "done") and persists to file.
func (j *JSONTaskSource) Close(id, reason string) error {
	tasks, err := j.readTasks()
	if err != nil {
		return err
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = "done"
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task %s not found", id)
	}

	return j.writeTasks(tasks)
}

// Reset sets a task's status back to "todo" and persists to file.
func (j *JSONTaskSource) Reset(id string) error {
	tasks, err := j.readTasks()
	if err != nil {
		return err
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = "todo"
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task %s not found", id)
	}

	return j.writeTasks(tasks)
}
