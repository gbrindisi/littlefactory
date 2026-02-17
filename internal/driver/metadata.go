// Package driver provides the loop orchestrator and metadata tracking.
package driver

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/yourusername/littlefactory/internal/config"
)

// pythonISOFormat matches Python's datetime.isoformat() output for naive datetimes.
const pythonISOFormat = "2006-01-02T15:04:05.999999"

// RunStatus represents the status of an entire run.
type RunStatus string

const (
	RunStatusRunning   RunStatus = "running"
	RunStatusCompleted RunStatus = "completed"
	RunStatusCancelled RunStatus = "cancelled"
	RunStatusFailed    RunStatus = "failed"
)

// IterationStatus represents the status of a single iteration.
type IterationStatus string

const (
	IterationStatusRunning   IterationStatus = "running"
	IterationStatusCompleted IterationStatus = "completed"
	IterationStatusFailed    IterationStatus = "failed"
	IterationStatusTimeout   IterationStatus = "timeout"
)

// IterationMetadata tracks per-iteration metadata including timing, status,
// task info, and output metrics.
type IterationMetadata struct {
	IterationNumber int             `json:"iteration_number"`
	StartedAt       time.Time       `json:"started_at"`
	EndedAt         *time.Time      `json:"ended_at,omitempty"`
	DurationSeconds *float64        `json:"duration_seconds,omitempty"`
	Status          IterationStatus `json:"status"`
	TaskID          *string         `json:"task_id,omitempty"`
	TaskTitle       *string         `json:"task_title,omitempty"`
	ExitCode        *int            `json:"exit_code,omitempty"`
	ErrorMessage    *string         `json:"error_message,omitempty"`
	OutputLines     *int            `json:"output_lines,omitempty"`
	OutputBytes     *int            `json:"output_bytes,omitempty"`
}

// RunMetadata tracks run-level metadata including run ID, timestamps,
// status, and aggregate statistics.
type RunMetadata struct {
	RunID                       string              `json:"run_id"`
	StartedAt                   time.Time           `json:"started_at"`
	EndedAt                     *time.Time          `json:"ended_at,omitempty"`
	Status                      RunStatus           `json:"status"`
	MaxIterations               int                 `json:"max_iterations"`
	TotalIterations             int                 `json:"total_iterations"`
	SuccessfulIterations        int                 `json:"successful_iterations"`
	FailedIterations            int                 `json:"failed_iterations"`
	TotalDurationSeconds        *float64            `json:"total_duration_seconds,omitempty"`
	AvgIterationDurationSeconds *float64            `json:"avg_iteration_duration_seconds,omitempty"`
	Iterations                  []IterationMetadata `json:"iterations"`
}

// runMetadataJSON is the JSON representation of RunMetadata with explicit field order.
type runMetadataJSON struct {
	RunID                       string                   `json:"run_id"`
	StartedAt                   string                   `json:"started_at"`
	EndedAt                     *string                  `json:"ended_at"`
	Status                      RunStatus                `json:"status"`
	MaxIterations               int                      `json:"max_iterations"`
	TotalIterations             int                      `json:"total_iterations"`
	SuccessfulIterations        int                      `json:"successful_iterations"`
	FailedIterations            int                      `json:"failed_iterations"`
	TotalDurationSeconds        *float64                 `json:"total_duration_seconds"`
	AvgIterationDurationSeconds *float64                 `json:"avg_iteration_duration_seconds"`
	Iterations                  []iterationMetadataJSON  `json:"iterations"`
}

// MarshalJSON implements custom JSON marshaling for RunMetadata to ensure
// timestamps are serialized in ISO8601 format matching Python's isoformat().
func (r *RunMetadata) MarshalJSON() ([]byte, error) {
	aux := runMetadataJSON{
		RunID:                       r.RunID,
		StartedAt:                   r.StartedAt.Format(pythonISOFormat),
		Status:                      r.Status,
		MaxIterations:               r.MaxIterations,
		TotalIterations:             r.TotalIterations,
		SuccessfulIterations:        r.SuccessfulIterations,
		FailedIterations:            r.FailedIterations,
		TotalDurationSeconds:        r.TotalDurationSeconds,
		AvgIterationDurationSeconds: r.AvgIterationDurationSeconds,
		Iterations:                  make([]iterationMetadataJSON, len(r.Iterations)),
	}
	if r.EndedAt != nil {
		endedAt := r.EndedAt.Format(pythonISOFormat)
		aux.EndedAt = &endedAt
	}
	for i, iter := range r.Iterations {
		aux.Iterations[i] = iter.toJSON()
	}
	return json.Marshal(aux)
}

// iterationMetadataJSON is the JSON representation of IterationMetadata with explicit field order.
type iterationMetadataJSON struct {
	IterationNumber int             `json:"iteration_number"`
	StartedAt       string          `json:"started_at"`
	EndedAt         *string         `json:"ended_at"`
	DurationSeconds *float64        `json:"duration_seconds"`
	Status          IterationStatus `json:"status"`
	TaskID          *string         `json:"task_id"`
	TaskTitle       *string         `json:"task_title"`
	ExitCode        *int            `json:"exit_code"`
	ErrorMessage    *string         `json:"error_message"`
	OutputLines     *int            `json:"output_lines"`
	OutputBytes     *int            `json:"output_bytes"`
}

// toJSON converts IterationMetadata to its JSON representation.
func (i *IterationMetadata) toJSON() iterationMetadataJSON {
	aux := iterationMetadataJSON{
		IterationNumber: i.IterationNumber,
		StartedAt:       i.StartedAt.Format(pythonISOFormat),
		DurationSeconds: i.DurationSeconds,
		Status:          i.Status,
		TaskID:          i.TaskID,
		TaskTitle:       i.TaskTitle,
		ExitCode:        i.ExitCode,
		ErrorMessage:    i.ErrorMessage,
		OutputLines:     i.OutputLines,
		OutputBytes:     i.OutputBytes,
	}
	if i.EndedAt != nil {
		endedAt := i.EndedAt.Format(pythonISOFormat)
		aux.EndedAt = &endedAt
	}
	return aux
}

// MarshalJSON implements custom JSON marshaling for IterationMetadata to ensure
// timestamps are serialized in ISO8601 format matching Python's isoformat().
func (i *IterationMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.toJSON())
}

// ptr is a helper function to create a pointer to a value.
func ptr[T any](v T) *T {
	return &v
}

// GenerateRunID creates a run ID using the format YYYYMMDD-HHMMSS.
func GenerateRunID() string {
	return time.Now().Format("20060102-150405")
}

// SaveMetadata writes the run metadata to <state_dir>/run_metadata.json.
// The projectRoot parameter is used to construct the path, and cfg provides
// the state directory configuration.
func SaveMetadata(projectRoot string, cfg *config.Config, metadata *RunMetadata) error {
	stateDir := filepath.Join(projectRoot, cfg.StateDir)
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(stateDir, "run_metadata.json")
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// CalculateAggregateStats updates the run metadata with aggregate statistics
// computed from the iteration data. This includes total_iterations,
// successful_iterations, failed_iterations, and avg_iteration_duration_seconds.
func (r *RunMetadata) CalculateAggregateStats() {
	r.TotalIterations = len(r.Iterations)

	var successful, failed int
	var totalDuration float64
	var completedCount int

	for _, iter := range r.Iterations {
		switch iter.Status {
		case IterationStatusCompleted:
			successful++
		case IterationStatusFailed, IterationStatusTimeout:
			failed++
		}

		if iter.DurationSeconds != nil {
			totalDuration += *iter.DurationSeconds
			completedCount++
		}
	}

	r.SuccessfulIterations = successful
	r.FailedIterations = failed

	if completedCount > 0 {
		avgDuration := totalDuration / float64(completedCount)
		r.AvgIterationDurationSeconds = &avgDuration
	}
}
