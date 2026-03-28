package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gbrindisi/littlefactory/internal/driver"
	"github.com/gbrindisi/littlefactory/internal/tasks"
)

func TestSummarizeTasks_AllDone(t *testing.T) {
	taskList := []tasks.Task{
		{ID: "1", Title: "Task A", Status: "done"},
		{ID: "2", Title: "Task B", Status: "done"},
		{ID: "3", Title: "Task C", Status: "done"},
	}

	s := summarizeTasks("feature-a", taskList)

	if s.Name != "feature-a" {
		t.Errorf("expected name 'feature-a', got %q", s.Name)
	}
	if s.Total != 3 {
		t.Errorf("expected total 3, got %d", s.Total)
	}
	if s.Done != 3 {
		t.Errorf("expected done 3, got %d", s.Done)
	}
	if s.InProgress != "" {
		t.Errorf("expected no in_progress, got %q", s.InProgress)
	}
}

func TestSummarizeTasks_Mixed(t *testing.T) {
	taskList := []tasks.Task{
		{ID: "1", Title: "Task A", Status: "done"},
		{ID: "2", Title: "Task B", Status: "in_progress"},
		{ID: "3", Title: "Task C", Status: "todo"},
		{ID: "4", Title: "Task D", Status: "done"},
		{ID: "5", Title: "Task E", Status: "done"},
		{ID: "6", Title: "Task F", Status: "todo"},
		{ID: "7", Title: "Task G", Status: "todo"},
	}

	s := summarizeTasks("feature-b", taskList)

	if s.Total != 7 {
		t.Errorf("expected total 7, got %d", s.Total)
	}
	if s.Done != 3 {
		t.Errorf("expected done 3, got %d", s.Done)
	}
	if s.InProgress != "Task B" {
		t.Errorf("expected in_progress 'Task B', got %q", s.InProgress)
	}
}

func TestSummarizeTasks_Empty(t *testing.T) {
	s := summarizeTasks("empty", nil)

	if s.Total != 0 {
		t.Errorf("expected total 0, got %d", s.Total)
	}
	if s.Done != 0 {
		t.Errorf("expected done 0, got %d", s.Done)
	}
}

func TestFormatSummary_Basic(t *testing.T) {
	s := taskSummary{Name: "feature-a", Total: 7, Done: 3}
	got := formatSummary(s)
	expected := "feature-a: 3/7 done"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFormatSummary_Complete(t *testing.T) {
	s := taskSummary{Name: "feature-a", Total: 5, Done: 5}
	got := formatSummary(s)
	expected := "feature-a: 5/5 done [complete]"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFormatSummary_InProgress(t *testing.T) {
	s := taskSummary{Name: "feature-a", Total: 7, Done: 3, InProgress: "Task B"}
	got := formatSummary(s)
	expected := `feature-a: 3/7 done (in_progress: "Task B")`
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFormatSummary_ZeroTotalNotComplete(t *testing.T) {
	s := taskSummary{Name: "empty", Total: 0, Done: 0}
	got := formatSummary(s)
	expected := "empty: 0/0 done"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestReadTasksFromPath_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	tasksPath := filepath.Join(tmpDir, "tasks.json")

	data, _ := json.Marshal(map[string]interface{}{
		"tasks": []map[string]string{
			{"id": "1", "title": "Task A", "status": "done"},
			{"id": "2", "title": "Task B", "status": "todo"},
		},
	})
	if err := os.WriteFile(tasksPath, data, 0644); err != nil {
		t.Fatal(err)
	}

	taskList, err := readTasksFromPath(tasksPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(taskList) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(taskList))
	}
	if taskList[0].Title != "Task A" {
		t.Errorf("expected first task title 'Task A', got %q", taskList[0].Title)
	}
}

func TestReadTasksFromPath_NotFound(t *testing.T) {
	_, err := readTasksFromPath("/nonexistent/tasks.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
	if !os.IsNotExist(err) {
		t.Errorf("expected IsNotExist error, got: %v", err)
	}
}

func TestReadTasksFromPath_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	tasksPath := filepath.Join(tmpDir, "tasks.json")

	if err := os.WriteFile(tasksPath, []byte("not json"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := readTasksFromPath(tasksPath)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestStatusCmd_HasChangeFlag(t *testing.T) {
	flag := statusCmd.Flags().Lookup("change")
	if flag == nil {
		t.Fatal("expected --change flag to be registered")
	}
	if flag.Shorthand != "c" {
		t.Errorf("expected shorthand 'c', got '%s'", flag.Shorthand)
	}
}

func TestStatusCmd_HasAllFlag(t *testing.T) {
	flag := statusCmd.Flags().Lookup("all")
	if flag == nil {
		t.Fatal("expected --all flag to be registered")
	}
}

func TestStatusCmd_HasVerboseFlag(t *testing.T) {
	flag := statusCmd.Flags().Lookup("verbose")
	if flag == nil {
		t.Fatal("expected --verbose flag to be registered")
	}
	if flag.Shorthand != "v" {
		t.Errorf("expected shorthand 'v', got '%s'", flag.Shorthand)
	}
}

func TestPrintVerboseTasks(t *testing.T) {
	// This test verifies printVerboseTasks doesn't panic with various statuses.
	// Output goes to stdout so we just ensure no crash.
	taskList := []tasks.Task{
		{ID: "1", Title: "Done Task", Status: "done"},
		{ID: "2", Title: "Active Task", Status: "in_progress"},
		{ID: "3", Title: "Pending Task", Status: "todo"},
	}

	// Redirect stdout to /dev/null for this test
	old := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = old }()

	printVerboseTasks(taskList)
}

func TestFormatSummaryWithMeta_NilMeta(t *testing.T) {
	s := taskSummary{Name: "feature-a", Total: 5, Done: 3, InProgress: "Task B"}
	got := formatSummaryWithMeta(s, nil)
	expected := `feature-a: 3/5 done (in_progress: "Task B")`
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFormatSummaryWithMeta_NilMetaComplete(t *testing.T) {
	s := taskSummary{Name: "feature-a", Total: 5, Done: 5}
	got := formatSummaryWithMeta(s, nil)
	expected := "feature-a: 5/5 done [complete]"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFormatSummaryWithMeta_Running(t *testing.T) {
	s := taskSummary{Name: "feature-a", Total: 5, Done: 3, InProgress: "Task B"}
	meta := &driver.RunMetadata{Status: driver.RunStatusRunning}
	got := formatSummaryWithMeta(s, meta)
	expected := `feature-a: 3/5 done [running] (in_progress: "Task B")`
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFormatSummaryWithMeta_Failed(t *testing.T) {
	s := taskSummary{Name: "feature-a", Total: 5, Done: 3}
	meta := &driver.RunMetadata{Status: driver.RunStatusFailed}
	got := formatSummaryWithMeta(s, meta)
	expected := "feature-a: 3/5 done [failed]"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFormatSummaryWithMeta_ReadyToMerge(t *testing.T) {
	s := taskSummary{Name: "feature-a", Total: 5, Done: 5}
	meta := &driver.RunMetadata{Status: driver.RunStatusCompleted}
	got := formatSummaryWithMeta(s, meta)
	expected := "feature-a: 5/5 done [completed] [ready to merge]"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFormatSummaryWithMeta_CompletedNotAllDone(t *testing.T) {
	s := taskSummary{Name: "feature-a", Total: 5, Done: 3}
	meta := &driver.RunMetadata{Status: driver.RunStatusCompleted}
	got := formatSummaryWithMeta(s, meta)
	expected := "feature-a: 3/5 done [completed]"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFormatSummaryWithMeta_Cancelled(t *testing.T) {
	s := taskSummary{Name: "feature-a", Total: 5, Done: 2}
	meta := &driver.RunMetadata{Status: driver.RunStatusCancelled}
	got := formatSummaryWithMeta(s, meta)
	expected := "feature-a: 2/5 done [cancelled]"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
