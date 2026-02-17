package driver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/littlefactory/internal/config"
)

func TestInitProgressFile_CreatesNewFile(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Create test config
	cfg := &config.Config{
		StateDir: ".littlefactory",
	}

	// Init progress file
	err := InitProgressFile(tmpDir, cfg)
	if err != nil {
		t.Fatalf("InitProgressFile failed: %v", err)
	}

	// Verify file exists
	filePath := filepath.Join(tmpDir, cfg.StateDir, ProgressFileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read progress file: %v", err)
	}

	// Verify header content
	contentStr := string(content)
	if !strings.HasPrefix(contentStr, "# Little Factory Progress Log\n") {
		t.Errorf("Expected header '# Little Factory Progress Log\\n', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "**Started:**") {
		t.Error("Expected '**Started:**' in header")
	}
	if !strings.Contains(contentStr, "---") {
		t.Error("Expected '---' separator in header")
	}
}

func TestInitProgressFile_PreservesExistingFile(t *testing.T) {
	// Create temp directory with existing progress file
	tmpDir := t.TempDir()
	cfg := &config.Config{
		StateDir: ".littlefactory",
	}
	stateDir := filepath.Join(tmpDir, cfg.StateDir)
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		t.Fatalf("Failed to create state dir: %v", err)
	}

	existingContent := "# Existing Progress\nSome existing data\n"
	filePath := filepath.Join(stateDir, ProgressFileName)
	if err := os.WriteFile(filePath, []byte(existingContent), 0644); err != nil {
		t.Fatalf("Failed to write existing file: %v", err)
	}

	// Init progress file
	err := InitProgressFile(tmpDir, cfg)
	if err != nil {
		t.Fatalf("InitProgressFile failed: %v", err)
	}

	// Verify content is unchanged
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read progress file: %v", err)
	}

	if string(content) != existingContent {
		t.Errorf("Expected existing content to be preserved, got: %s", string(content))
	}
}

func TestAppendSessionToProgress(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	cfg := &config.Config{
		StateDir: ".littlefactory",
	}
	stateDir := filepath.Join(tmpDir, cfg.StateDir)
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		t.Fatalf("Failed to create state dir: %v", err)
	}

	// Create initial progress file
	filePath := filepath.Join(stateDir, ProgressFileName)
	if err := os.WriteFile(filePath, []byte("# Header\n"), 0644); err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	// Append session
	err := AppendSessionToProgress(tmpDir, cfg, 1, "task-123", "completed")
	if err != nil {
		t.Fatalf("AppendSessionToProgress failed: %v", err)
	}

	// Verify content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "## Iteration 1") {
		t.Error("Expected iteration header")
	}
	if !strings.Contains(contentStr, "- **Task:** task-123") {
		t.Error("Expected task ID with bold label")
	}
	if !strings.Contains(contentStr, "- **Status:** completed") {
		t.Error("Expected status with bold label")
	}
	if !strings.Contains(contentStr, "---") {
		t.Error("Expected separator")
	}
}

func TestAppendSessionToProgress_MultipleAppends(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	cfg := &config.Config{
		StateDir: ".littlefactory",
	}
	stateDir := filepath.Join(tmpDir, cfg.StateDir)
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		t.Fatalf("Failed to create state dir: %v", err)
	}

	// Create initial progress file
	filePath := filepath.Join(stateDir, ProgressFileName)
	initialContent := "# Header\n---\n"
	if err := os.WriteFile(filePath, []byte(initialContent), 0644); err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	// Append multiple sessions
	if err := AppendSessionToProgress(tmpDir, cfg, 1, "task-1", "completed"); err != nil {
		t.Fatalf("First append failed: %v", err)
	}
	if err := AppendSessionToProgress(tmpDir, cfg, 2, "task-2", "completed"); err != nil {
		t.Fatalf("Second append failed: %v", err)
	}

	// Verify all content preserved
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)

	// Check initial content preserved
	if !strings.HasPrefix(contentStr, "# Header\n---\n") {
		t.Error("Initial content not preserved")
	}

	// Check both iterations present
	if !strings.Contains(contentStr, "## Iteration 1") {
		t.Error("First iteration missing")
	}
	if !strings.Contains(contentStr, "## Iteration 2") {
		t.Error("Second iteration missing")
	}
	if !strings.Contains(contentStr, "task-1") {
		t.Error("First task missing")
	}
	if !strings.Contains(contentStr, "task-2") {
		t.Error("Second task missing")
	}
}

func TestAppendSessionToProgress_CreatesFileIfNotExists(t *testing.T) {
	// Create temp directory but not the state dir
	tmpDir := t.TempDir()
	cfg := &config.Config{
		StateDir: ".littlefactory",
	}

	// Append session (file doesn't exist yet)
	err := AppendSessionToProgress(tmpDir, cfg, 1, "task-1", "completed")
	if err != nil {
		t.Fatalf("AppendSessionToProgress failed: %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(tmpDir, cfg.StateDir, ProgressFileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "## Iteration 1") {
		t.Error("Expected iteration content")
	}
}

func TestProgressFilePath(t *testing.T) {
	projectRoot := "/path/to/project"
	cfg := &config.Config{
		StateDir: ".littlefactory",
	}
	expected := filepath.Join(projectRoot, cfg.StateDir, ProgressFileName)
	actual := ProgressFilePath(projectRoot, cfg)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
