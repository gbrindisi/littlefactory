package driver

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/gbrindisi/littlefactory/internal/config"
)

func TestGenerateRunID(t *testing.T) {
	runID := GenerateRunID()

	// Verify format matches YYYYMMDD-HHMMSS pattern
	pattern := regexp.MustCompile(`^\d{8}-\d{6}$`)
	if !pattern.MatchString(runID) {
		t.Errorf("GenerateRunID() = %q, expected format YYYYMMDD-HHMMSS", runID)
	}

	// Verify it's a valid parseable time
	_, err := time.Parse("20060102-150405", runID)
	if err != nil {
		t.Errorf("GenerateRunID() returned unparseable time %q: %v", runID, err)
	}
}

func TestSaveMetadata(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Create test config
	cfg := &config.Config{
		StateDir: ".littlefactory",
	}

	// Create test metadata
	now := time.Now()
	metadata := &RunMetadata{
		RunID:                "20240130-150405",
		StartedAt:            now,
		Status:               RunStatusCompleted,
		MaxIterations:        10,
		TotalIterations:      2,
		SuccessfulIterations: 2,
		FailedIterations:     0,
		Iterations:           []IterationMetadata{},
	}

	// Save metadata
	err := SaveMetadata(tmpDir, cfg, metadata)
	if err != nil {
		t.Fatalf("SaveMetadata() error = %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(tmpDir, ".littlefactory", "run_metadata.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("SaveMetadata() did not create file at %s", filePath)
	}

	// Verify content is valid JSON
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read metadata file: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Errorf("SaveMetadata() wrote invalid JSON: %v", err)
	}

	// Verify run_id field
	if parsed["run_id"] != "20240130-150405" {
		t.Errorf("run_id = %v, expected 20240130-150405", parsed["run_id"])
	}
}

func TestSaveMetadataCreatesStateDir(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		StateDir: ".littlefactory",
	}

	metadata := &RunMetadata{
		RunID:     "20240130-150405",
		StartedAt: time.Now(),
		Status:    RunStatusRunning,
	}

	// state directory doesn't exist yet
	stateDir := filepath.Join(tmpDir, ".littlefactory")
	if _, err := os.Stat(stateDir); !os.IsNotExist(err) {
		t.Fatal(".littlefactory/ directory should not exist initially")
	}

	// SaveMetadata should create it
	err := SaveMetadata(tmpDir, cfg, metadata)
	if err != nil {
		t.Fatalf("SaveMetadata() error = %v", err)
	}

	if _, err := os.Stat(stateDir); os.IsNotExist(err) {
		t.Error("SaveMetadata() did not create .littlefactory/ directory")
	}
}

func TestCalculateAggregateStats(t *testing.T) {
	duration1 := 10.5
	duration2 := 20.5
	duration3 := 15.0

	metadata := &RunMetadata{
		RunID:     "20240130-150405",
		StartedAt: time.Now(),
		Status:    RunStatusCompleted,
		Iterations: []IterationMetadata{
			{IterationNumber: 1, Status: IterationStatusCompleted, DurationSeconds: &duration1},
			{IterationNumber: 2, Status: IterationStatusCompleted, DurationSeconds: &duration2},
			{IterationNumber: 3, Status: IterationStatusFailed, DurationSeconds: &duration3},
			{IterationNumber: 4, Status: IterationStatusTimeout, DurationSeconds: nil},
		},
	}

	metadata.CalculateAggregateStats()

	// Verify total iterations
	if metadata.TotalIterations != 4 {
		t.Errorf("TotalIterations = %d, expected 4", metadata.TotalIterations)
	}

	// Verify successful iterations (only completed count as successful)
	if metadata.SuccessfulIterations != 2 {
		t.Errorf("SuccessfulIterations = %d, expected 2", metadata.SuccessfulIterations)
	}

	// Verify failed iterations (failed + timeout)
	if metadata.FailedIterations != 2 {
		t.Errorf("FailedIterations = %d, expected 2", metadata.FailedIterations)
	}

	// Verify average duration (only iterations with duration)
	if metadata.AvgIterationDurationSeconds == nil {
		t.Error("AvgIterationDurationSeconds should not be nil")
	} else {
		expected := (10.5 + 20.5 + 15.0) / 3.0
		if *metadata.AvgIterationDurationSeconds != expected {
			t.Errorf("AvgIterationDurationSeconds = %f, expected %f", *metadata.AvgIterationDurationSeconds, expected)
		}
	}
}

func TestCalculateAggregateStatsEmpty(t *testing.T) {
	metadata := &RunMetadata{
		RunID:      "20240130-150405",
		StartedAt:  time.Now(),
		Status:     RunStatusCompleted,
		Iterations: []IterationMetadata{},
	}

	metadata.CalculateAggregateStats()

	if metadata.TotalIterations != 0 {
		t.Errorf("TotalIterations = %d, expected 0", metadata.TotalIterations)
	}
	if metadata.SuccessfulIterations != 0 {
		t.Errorf("SuccessfulIterations = %d, expected 0", metadata.SuccessfulIterations)
	}
	if metadata.FailedIterations != 0 {
		t.Errorf("FailedIterations = %d, expected 0", metadata.FailedIterations)
	}
	if metadata.AvgIterationDurationSeconds != nil {
		t.Errorf("AvgIterationDurationSeconds = %v, expected nil", metadata.AvgIterationDurationSeconds)
	}
}

func TestRunMetadataMarshalJSON(t *testing.T) {
	startedAt := time.Date(2024, 1, 30, 15, 4, 5, 0, time.UTC)
	endedAt := time.Date(2024, 1, 30, 15, 10, 0, 0, time.UTC)
	duration := 355.0

	metadata := &RunMetadata{
		RunID:                "20240130-150405",
		StartedAt:            startedAt,
		EndedAt:              &endedAt,
		Status:               RunStatusCompleted,
		MaxIterations:        10,
		TotalIterations:      2,
		SuccessfulIterations: 2,
		FailedIterations:     0,
		TotalDurationSeconds: &duration,
		Iterations:           []IterationMetadata{},
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify ISO8601 timestamp format (Python isoformat style, no timezone)
	startedAtStr, ok := parsed["started_at"].(string)
	if !ok {
		t.Error("started_at should be a string")
	} else if startedAtStr != "2024-01-30T15:04:05" {
		t.Errorf("started_at = %q, expected 2024-01-30T15:04:05", startedAtStr)
	}

	endedAtStr, ok := parsed["ended_at"].(string)
	if !ok {
		t.Error("ended_at should be a string")
	} else if endedAtStr != "2024-01-30T15:10:00" {
		t.Errorf("ended_at = %q, expected 2024-01-30T15:10:00", endedAtStr)
	}
}

func TestIterationMetadataMarshalJSON(t *testing.T) {
	startedAt := time.Date(2024, 1, 30, 15, 4, 5, 0, time.UTC)
	endedAt := time.Date(2024, 1, 30, 15, 5, 0, 0, time.UTC)
	duration := 55.0
	exitCode := 0
	outputLines := 100
	outputBytes := 5000
	taskID := "task-123"
	taskTitle := "Test Task"

	iter := &IterationMetadata{
		IterationNumber: 1,
		StartedAt:       startedAt,
		EndedAt:         &endedAt,
		DurationSeconds: &duration,
		Status:          IterationStatusCompleted,
		TaskID:          &taskID,
		TaskTitle:       &taskTitle,
		ExitCode:        &exitCode,
		OutputLines:     &outputLines,
		OutputBytes:     &outputBytes,
	}

	data, err := json.Marshal(iter)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify ISO8601 timestamp format (Python isoformat style, no timezone)
	startedAtStr, ok := parsed["started_at"].(string)
	if !ok {
		t.Error("started_at should be a string")
	} else if startedAtStr != "2024-01-30T15:04:05" {
		t.Errorf("started_at = %q, expected 2024-01-30T15:04:05", startedAtStr)
	}

	endedAtStr, ok := parsed["ended_at"].(string)
	if !ok {
		t.Error("ended_at should be a string")
	} else if endedAtStr != "2024-01-30T15:05:00" {
		t.Errorf("ended_at = %q, expected 2024-01-30T15:05:00", endedAtStr)
	}

	// Verify other fields are present
	if parsed["iteration_number"] != float64(1) {
		t.Errorf("iteration_number = %v, expected 1", parsed["iteration_number"])
	}
	if parsed["task_id"] != "task-123" {
		t.Errorf("task_id = %v, expected task-123", parsed["task_id"])
	}
}

func TestRunMetadataMarshalJSONNullEndedAt(t *testing.T) {
	startedAt := time.Date(2024, 1, 30, 15, 4, 5, 0, time.UTC)

	metadata := &RunMetadata{
		RunID:      "20240130-150405",
		StartedAt:  startedAt,
		EndedAt:    nil, // Still running
		Status:     RunStatusRunning,
		Iterations: []IterationMetadata{},
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// ended_at should be present with null value (matching Python behavior)
	if val, exists := parsed["ended_at"]; !exists {
		t.Error("ended_at should be present")
	} else if val != nil {
		t.Errorf("ended_at should be null, got %v", val)
	}
}
