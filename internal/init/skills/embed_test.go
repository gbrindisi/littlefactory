package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractSkills(t *testing.T) {
	tmpDir := t.TempDir()

	if err := ExtractSkills(tmpDir); err != nil {
		t.Fatalf("ExtractSkills failed: %v", err)
	}

	// No embedded skills exist, so .littlefactory/skills/ should be created but empty
	skillsDir := filepath.Join(tmpDir, ".littlefactory", "skills")
	info, err := os.Stat(skillsDir)
	if err != nil {
		t.Fatalf("expected %s to exist: %v", skillsDir, err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %s to be a directory", skillsDir)
	}

	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		t.Fatalf("failed to read skills directory: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected no skill directories, got %d", len(entries))
	}
}

func TestExtractSkillsIdempotent(t *testing.T) {
	tmpDir := t.TempDir()

	// Extract twice to verify idempotency (no errors on overwrite)
	if err := ExtractSkills(tmpDir); err != nil {
		t.Fatalf("first ExtractSkills failed: %v", err)
	}
	if err := ExtractSkills(tmpDir); err != nil {
		t.Fatalf("second ExtractSkills failed: %v", err)
	}
}
