package template

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/littlefactory/internal/tasks"
)

func TestEmbeddedTemplateLoaded(t *testing.T) {
	if embeddedTemplate == "" {
		t.Fatal("embedded template should not be empty")
	}
	if len(embeddedTemplate) < 100 {
		t.Fatalf("embedded template seems too short: %d bytes", len(embeddedTemplate))
	}
}

func TestRenderWithTask(t *testing.T) {
	tmpl := `Task: {task_id}
Title: {task_title}
Description: {task_description}`

	task := &tasks.Task{
		ID:          "test-123",
		Title:       "Test Task",
		Description: "This is a test description",
	}

	result := Render(tmpl, task)

	expected := `Task: test-123
Title: Test Task
Description: This is a test description`

	if result != expected {
		t.Errorf("Render mismatch:\ngot:  %q\nwant: %q", result, expected)
	}
}

func TestRenderWithNilTask(t *testing.T) {
	tmpl := `Task: {task_id}
Title: {task_title}
Description: {task_description}`

	result := Render(tmpl, nil)

	if result != tmpl {
		t.Errorf("Render with nil task should return template unchanged:\ngot:  %q\nwant: %q", result, tmpl)
	}
}

func TestRenderMultipleOccurrences(t *testing.T) {
	tmpl := `{task_id} - {task_id} - {task_id}`

	task := &tasks.Task{
		ID: "abc",
	}

	result := Render(tmpl, task)
	expected := "abc - abc - abc"

	if result != expected {
		t.Errorf("Render should replace all occurrences:\ngot:  %q\nwant: %q", result, expected)
	}
}

func TestLoadWithLocalOverride(t *testing.T) {
	// Create a temp state directory with agents/WORKER.md
	tmpDir := t.TempDir()
	agentsDir := filepath.Join(tmpDir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatalf("failed to create agents dir: %v", err)
	}

	localContent := "# Local Override Template\n{task_id}"
	localPath := filepath.Join(agentsDir, "WORKER.md")
	if err := os.WriteFile(localPath, []byte(localContent), 0644); err != nil {
		t.Fatalf("failed to write local template: %v", err)
	}

	result, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if result != localContent {
		t.Errorf("Load should return local override:\ngot:  %q\nwant: %q", result, localContent)
	}
}

func TestLoadWithoutLocalOverride(t *testing.T) {
	// Create a temp state directory without agents/WORKER.md
	tmpDir := t.TempDir()

	result, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if result != embeddedTemplate {
		t.Errorf("Load should return embedded template when local not found:\ngot length:  %d\nwant length: %d", len(result), len(embeddedTemplate))
	}
}

func TestLoadWithEmptyAgentsDir(t *testing.T) {
	// Create a temp state directory with empty agents dir (no WORKER.md)
	tmpDir := t.TempDir()
	agentsDir := filepath.Join(tmpDir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatalf("failed to create agents dir: %v", err)
	}

	result, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if result != embeddedTemplate {
		t.Errorf("Load should return embedded template when local not found:\ngot length:  %d\nwant length: %d", len(result), len(embeddedTemplate))
	}
}
