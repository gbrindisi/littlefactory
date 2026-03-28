package template

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gbrindisi/littlefactory/internal/tasks"
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

func TestEmbeddedVerifierTemplateLoaded(t *testing.T) {
	if embeddedVerifierTemplate == "" {
		t.Fatal("embedded verifier template should not be empty")
	}
	if len(embeddedVerifierTemplate) < 100 {
		t.Fatalf("embedded verifier template seems too short: %d bytes", len(embeddedVerifierTemplate))
	}
}

func TestLoadVerifierWithLocalOverride(t *testing.T) {
	tmpDir := t.TempDir()
	agentsDir := filepath.Join(tmpDir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatalf("failed to create agents dir: %v", err)
	}

	localContent := "# Custom Verifier\n{change_name}"
	if err := os.WriteFile(filepath.Join(agentsDir, "VERIFIER.md"), []byte(localContent), 0644); err != nil {
		t.Fatalf("failed to write local verifier template: %v", err)
	}

	result, err := LoadVerifier(tmpDir)
	if err != nil {
		t.Fatalf("LoadVerifier returned error: %v", err)
	}

	if result != localContent {
		t.Errorf("LoadVerifier should return local override:\ngot:  %q\nwant: %q", result, localContent)
	}
}

func TestLoadVerifierWithoutLocalOverride(t *testing.T) {
	tmpDir := t.TempDir()

	result, err := LoadVerifier(tmpDir)
	if err != nil {
		t.Fatalf("LoadVerifier returned error: %v", err)
	}

	if result != embeddedVerifierTemplate {
		t.Errorf("LoadVerifier should return embedded template when local not found:\ngot length:  %d\nwant length: %d", len(result), len(embeddedVerifierTemplate))
	}
}

func TestRenderVerifierWithContext(t *testing.T) {
	tmpl := `Change: {change_name}
Path: {change_path}
Proposal: {proposal_path}
Specs: {specs_paths}
Design: {design_path}
Tasks: {tasks_path}`

	ctx := &ChangeContext{
		ChangeName:   "my-change",
		ChangePath:   ".littlefactory/changes/my-change/",
		ProposalPath: ".littlefactory/changes/my-change/proposal.md",
		SpecsPaths:   "specs/foo/spec.md, specs/bar/spec.md",
		DesignPath:   ".littlefactory/changes/my-change/design.md",
		TasksPath:    ".littlefactory/changes/my-change/tasks.json",
	}

	result := RenderVerifier(tmpl, ctx)

	expected := `Change: my-change
Path: .littlefactory/changes/my-change/
Proposal: .littlefactory/changes/my-change/proposal.md
Specs: specs/foo/spec.md, specs/bar/spec.md
Design: .littlefactory/changes/my-change/design.md
Tasks: .littlefactory/changes/my-change/tasks.json`

	if result != expected {
		t.Errorf("RenderVerifier mismatch:\ngot:  %q\nwant: %q", result, expected)
	}
}

func TestRenderVerifierWithNilContext(t *testing.T) {
	tmpl := `{change_name} - {change_path}`

	result := RenderVerifier(tmpl, nil)

	if result != tmpl {
		t.Errorf("RenderVerifier with nil context should return template unchanged:\ngot:  %q\nwant: %q", result, tmpl)
	}
}
