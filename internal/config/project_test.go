package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindProjectRootFrom_CurrentDirectory(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	if err := os.WriteFile(factoryfile, []byte("# test factoryfile"), 0644); err != nil {
		t.Fatalf("failed to create Factoryfile: %v", err)
	}

	// Test from directory containing Factoryfile
	root, err := findProjectRootFrom(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if root != tmpDir {
		t.Errorf("expected %q, got %q", tmpDir, root)
	}
}

func TestFindProjectRootFrom_ParentDirectory(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	subDir := filepath.Join(tmpDir, "subdir", "nested")

	if err := os.WriteFile(factoryfile, []byte("# test factoryfile"), 0644); err != nil {
		t.Fatalf("failed to create Factoryfile: %v", err)
	}
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create nested dir: %v", err)
	}

	// Test from nested subdirectory
	root, err := findProjectRootFrom(subDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if root != tmpDir {
		t.Errorf("expected %q, got %q", tmpDir, root)
	}
}

func TestFindProjectRootFrom_NotFound(t *testing.T) {
	// Create temp directory without Factoryfile
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	// Test that error is returned when Factoryfile not found
	_, err := findProjectRootFrom(subDir)
	if err == nil {
		t.Fatal("expected error when Factoryfile not found, got nil")
	}

	if err != os.ErrNotExist {
		t.Errorf("expected os.ErrNotExist, got %v", err)
	}
}

func TestFindProjectRootFrom_FactoryfileIsDirectory(t *testing.T) {
	// Create temp directory with Factoryfile as a directory (not file)
	tmpDir := t.TempDir()
	factoryfileDir := filepath.Join(tmpDir, "Factoryfile")

	// Create Factoryfile as a directory, not a file
	if err := os.Mkdir(factoryfileDir, 0755); err != nil {
		t.Fatalf("failed to create Factoryfile dir: %v", err)
	}

	// Should NOT treat directory as valid marker, should return error
	_, err := findProjectRootFrom(tmpDir)
	if err == nil {
		t.Fatal("expected error when Factoryfile is a directory, got nil")
	}

	if err != os.ErrNotExist {
		t.Errorf("expected os.ErrNotExist, got %v", err)
	}
}

func TestFindProjectRootFrom_AbsolutePath(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	factoryfile := filepath.Join(tmpDir, "Factoryfile")
	if err := os.WriteFile(factoryfile, []byte("# test factoryfile"), 0644); err != nil {
		t.Fatalf("failed to create Factoryfile: %v", err)
	}

	// Test that returned path is absolute
	root, err := findProjectRootFrom(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !filepath.IsAbs(root) {
		t.Errorf("expected absolute path, got %q", root)
	}
}

func TestTasksDir(t *testing.T) {
	tests := []struct {
		projectRoot string
		expected    string
	}{
		{"/home/user/project", "/home/user/project/tasks"},
		{"/tmp/test", "/tmp/test/tasks"},
		{"/", "/tasks"},
	}

	for _, tt := range tests {
		result := TasksDir(tt.projectRoot)
		if result != tt.expected {
			t.Errorf("TasksDir(%q) = %q, want %q", tt.projectRoot, result, tt.expected)
		}
	}
}

func TestEnsureTasksDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Ensure tasks directory is created
	err := EnsureTasksDir(tmpDir)
	if err != nil {
		t.Fatalf("EnsureTasksDir failed: %v", err)
	}

	tasksPath := filepath.Join(tmpDir, "tasks")
	info, err := os.Stat(tasksPath)
	if err != nil {
		t.Fatalf("tasks directory not created: %v", err)
	}

	if !info.IsDir() {
		t.Error("tasks path is not a directory")
	}
}

func TestEnsureTasksDir_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	tasksPath := filepath.Join(tmpDir, "tasks")

	// Pre-create tasks directory
	if err := os.Mkdir(tasksPath, 0755); err != nil {
		t.Fatalf("failed to pre-create tasks dir: %v", err)
	}

	// Should not error when directory already exists
	err := EnsureTasksDir(tmpDir)
	if err != nil {
		t.Fatalf("EnsureTasksDir failed on existing dir: %v", err)
	}
}
