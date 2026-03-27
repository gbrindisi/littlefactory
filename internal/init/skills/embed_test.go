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

	// Verify .littlefactory/skills/ was created with embedded skills
	skillsDir := filepath.Join(tmpDir, ".littlefactory", "skills")
	info, err := os.Stat(skillsDir)
	if err != nil {
		t.Fatalf("expected %s to exist: %v", skillsDir, err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %s to be a directory", skillsDir)
	}

	// Verify lf-explore skill was extracted
	exploreDir := filepath.Join(skillsDir, "lf-explore")
	if _, err := os.Stat(exploreDir); err != nil {
		t.Fatalf("expected lf-explore skill directory to exist: %v", err)
	}
	skillFile := filepath.Join(exploreDir, "SKILL.md")
	if _, err := os.Stat(skillFile); err != nil {
		t.Fatalf("expected lf-explore/SKILL.md to exist: %v", err)
	}

	// Verify lf-formalize skill was extracted
	formalizeDir := filepath.Join(skillsDir, "lf-formalize")
	if _, err := os.Stat(formalizeDir); err != nil {
		t.Fatalf("expected lf-formalize skill directory to exist: %v", err)
	}
	formalizeSkillFile := filepath.Join(formalizeDir, "SKILL.md")
	if _, err := os.Stat(formalizeSkillFile); err != nil {
		t.Fatalf("expected lf-formalize/SKILL.md to exist: %v", err)
	}

	// Verify lf-do skill was extracted
	doDir := filepath.Join(skillsDir, "lf-do")
	if _, err := os.Stat(doDir); err != nil {
		t.Fatalf("expected lf-do skill directory to exist: %v", err)
	}
	doSkillFile := filepath.Join(doDir, "SKILL.md")
	if _, err := os.Stat(doSkillFile); err != nil {
		t.Fatalf("expected lf-do/SKILL.md to exist: %v", err)
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
