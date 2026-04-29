package worktree

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(b)
}

func TestSeed_CopiesChangeAndSpecs(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "src")
	dst := filepath.Join(tmpDir, "dst")

	// Source change artifacts
	writeFile(t, filepath.Join(src, ".littlefactory", "changes", "feat-x", "proposal.md"), "proposal")
	writeFile(t, filepath.Join(src, ".littlefactory", "changes", "feat-x", "tasks.json"), `[]`)
	writeFile(t, filepath.Join(src, ".littlefactory", "changes", "feat-x", "specs", "core", "spec.md"), "delta")

	// Project specs
	writeFile(t, filepath.Join(src, ".littlefactory", "specs", "core", "spec.md"), "core spec")

	// Unrelated change should not be copied
	writeFile(t, filepath.Join(src, ".littlefactory", "changes", "other", "proposal.md"), "other")

	// Per-run state should not be copied
	writeFile(t, filepath.Join(src, ".littlefactory", "progress.md"), "progress")
	writeFile(t, filepath.Join(src, ".littlefactory", "run_metadata.json"), "{}")

	if err := os.MkdirAll(dst, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := Seed(src, dst, "feat-x"); err != nil {
		t.Fatalf("Seed failed: %v", err)
	}

	// Change artifacts present
	if got := readFile(t, filepath.Join(dst, ".littlefactory", "changes", "feat-x", "proposal.md")); got != "proposal" {
		t.Fatalf("proposal not copied: %q", got)
	}
	if got := readFile(t, filepath.Join(dst, ".littlefactory", "changes", "feat-x", "tasks.json")); got != "[]" {
		t.Fatalf("tasks.json not copied: %q", got)
	}
	if got := readFile(t, filepath.Join(dst, ".littlefactory", "changes", "feat-x", "specs", "core", "spec.md")); got != "delta" {
		t.Fatalf("nested spec not copied: %q", got)
	}

	// Project specs present
	if got := readFile(t, filepath.Join(dst, ".littlefactory", "specs", "core", "spec.md")); got != "core spec" {
		t.Fatalf("project spec not copied: %q", got)
	}

	// Unrelated change should NOT be copied
	if _, err := os.Stat(filepath.Join(dst, ".littlefactory", "changes", "other")); !os.IsNotExist(err) {
		t.Fatalf("unrelated change directory should not be copied, stat err=%v", err)
	}

	// Per-run state should NOT be copied
	if _, err := os.Stat(filepath.Join(dst, ".littlefactory", "progress.md")); !os.IsNotExist(err) {
		t.Fatalf("progress.md should not be copied, stat err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(dst, ".littlefactory", "run_metadata.json")); !os.IsNotExist(err) {
		t.Fatalf("run_metadata.json should not be copied, stat err=%v", err)
	}
}

func TestSeed_MissingChangeDir(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "src")
	dst := filepath.Join(tmpDir, "dst")
	if err := os.MkdirAll(src, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		t.Fatal(err)
	}

	// No .littlefactory at all — should not error.
	if err := Seed(src, dst, "feat-x"); err != nil {
		t.Fatalf("Seed should succeed when sources are missing: %v", err)
	}
}

func TestSeed_OnlySpecsPresent(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "src")
	dst := filepath.Join(tmpDir, "dst")

	writeFile(t, filepath.Join(src, ".littlefactory", "specs", "core", "spec.md"), "core")
	if err := os.MkdirAll(dst, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := Seed(src, dst, "feat-x"); err != nil {
		t.Fatalf("Seed failed: %v", err)
	}

	if got := readFile(t, filepath.Join(dst, ".littlefactory", "specs", "core", "spec.md")); got != "core" {
		t.Fatalf("specs not copied: %q", got)
	}
}
