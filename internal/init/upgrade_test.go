package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpgrade_FailsWithoutFactoryfile(t *testing.T) {
	dir := t.TempDir()

	err := Upgrade(dir)
	if err == nil {
		t.Fatal("expected Upgrade to fail without Factoryfile")
	}
	if !strings.Contains(err.Error(), "no Factoryfile found") {
		t.Errorf("expected 'no Factoryfile found' error, got: %v", err)
	}
}

func TestUpgrade_WithFactoryfile(t *testing.T) {
	dir := t.TempDir()

	// Create Factoryfile (required for upgrade)
	factoryfilePath := filepath.Join(dir, "Factoryfile")
	if err := os.WriteFile(factoryfilePath, []byte(DefaultFactoryfile), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Upgrade(dir); err != nil {
		t.Fatalf("Upgrade failed: %v", err)
	}

	// Verify AGENTS.md was created
	agentsPath := filepath.Join(dir, "AGENTS.md")
	if _, err := os.Stat(agentsPath); err != nil {
		t.Fatalf("expected AGENTS.md to exist: %v", err)
	}

	// Verify CLAUDE.md is a symlink
	claudePath := filepath.Join(dir, "CLAUDE.md")
	target, err := os.Readlink(claudePath)
	if err != nil {
		t.Fatalf("expected CLAUDE.md to be a symlink: %v", err)
	}
	if target != "AGENTS.md" {
		t.Errorf("expected CLAUDE.md -> AGENTS.md, got -> %s", target)
	}

	// Verify .gitignore was created
	gitignorePath := filepath.Join(dir, ".gitignore")
	gitignoreContent, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("expected .gitignore to exist: %v", err)
	}
	if !strings.Contains(string(gitignoreContent), ".littlefactory/run_metadata.json") {
		t.Error(".gitignore missing .littlefactory/run_metadata.json")
	}
	if !strings.Contains(string(gitignoreContent), ".littlefactory/tasks.json") {
		t.Error(".gitignore missing .littlefactory/tasks.json")
	}

	// Verify changes directory was created
	changesDir := filepath.Join(dir, ".littlefactory", "changes")
	if _, err := os.Stat(changesDir); err != nil {
		t.Fatalf("expected .littlefactory/changes/ to exist: %v", err)
	}
}

func TestUpgrade_WithFactoryfileYAML(t *testing.T) {
	dir := t.TempDir()

	// Create Factoryfile.yaml (alternate name)
	factoryfilePath := filepath.Join(dir, "Factoryfile.yaml")
	if err := os.WriteFile(factoryfilePath, []byte(DefaultFactoryfile), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Upgrade(dir); err != nil {
		t.Fatalf("Upgrade failed: %v", err)
	}

	// Verify core artifacts created
	if _, err := os.Stat(filepath.Join(dir, "AGENTS.md")); err != nil {
		t.Fatalf("expected AGENTS.md to exist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".gitignore")); err != nil {
		t.Fatalf("expected .gitignore to exist: %v", err)
	}
}

func TestUpgrade_Idempotent(t *testing.T) {
	dir := t.TempDir()

	// Create Factoryfile
	factoryfilePath := filepath.Join(dir, "Factoryfile")
	if err := os.WriteFile(factoryfilePath, []byte(DefaultFactoryfile), 0o644); err != nil {
		t.Fatal(err)
	}

	// First upgrade
	if err := Upgrade(dir); err != nil {
		t.Fatalf("first Upgrade failed: %v", err)
	}

	// Capture state after first upgrade
	agentsContent1, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	gitignoreContent1, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}

	// Second upgrade should be idempotent
	if err := Upgrade(dir); err != nil {
		t.Fatalf("second Upgrade failed: %v", err)
	}

	// Verify state unchanged
	agentsContent2, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	gitignoreContent2, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}

	if string(agentsContent1) != string(agentsContent2) {
		t.Error("AGENTS.md content changed after idempotent upgrade")
	}
	if string(gitignoreContent1) != string(gitignoreContent2) {
		t.Error(".gitignore content changed after idempotent upgrade")
	}
}

func TestUpgrade_WithExistingClaudeMD(t *testing.T) {
	dir := t.TempDir()

	// Create Factoryfile
	factoryfilePath := filepath.Join(dir, "Factoryfile")
	if err := os.WriteFile(factoryfilePath, []byte(DefaultFactoryfile), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create existing CLAUDE.md
	claudePath := filepath.Join(dir, "CLAUDE.md")
	originalContent := "# Custom Instructions\n"
	if err := os.WriteFile(claudePath, []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Upgrade(dir); err != nil {
		t.Fatalf("Upgrade failed: %v", err)
	}

	// CLAUDE.md should have been migrated
	agentsPath := filepath.Join(dir, "AGENTS.md")
	content, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != originalContent {
		t.Error("AGENTS.md should contain migrated CLAUDE.md content")
	}

	// CLAUDE.md should be a symlink
	target, err := os.Readlink(claudePath)
	if err != nil {
		t.Fatalf("expected CLAUDE.md to be a symlink: %v", err)
	}
	if target != "AGENTS.md" {
		t.Errorf("expected CLAUDE.md -> AGENTS.md, got -> %s", target)
	}
}

func TestUpgrade_WithClaudeDir(t *testing.T) {
	dir := t.TempDir()

	// Create Factoryfile and .claude/ directory
	if err := os.WriteFile(filepath.Join(dir, "Factoryfile"), []byte(DefaultFactoryfile), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := Upgrade(dir); err != nil {
		t.Fatalf("Upgrade failed: %v", err)
	}

	// Verify embedded skill symlinks were created in .claude/skills/
	claudeSkillsPath := filepath.Join(dir, ".claude", "skills")
	entries, err := os.ReadDir(claudeSkillsPath)
	if err != nil {
		t.Fatalf("failed to read .claude/skills/: %v", err)
	}
	if len(entries) != 4 {
		t.Errorf("expected 4 skill symlinks (lf-do, lf-explore, lf-formalize, lf-verify), got %d", len(entries))
	}
}

func TestUpgrade_DoesNotCreateFactoryfile(t *testing.T) {
	dir := t.TempDir()

	// Create Factoryfile with custom content
	customContent := "max_iterations: 5\n"
	factoryfilePath := filepath.Join(dir, "Factoryfile")
	if err := os.WriteFile(factoryfilePath, []byte(customContent), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Upgrade(dir); err != nil {
		t.Fatalf("Upgrade failed: %v", err)
	}

	// Factoryfile should not have been modified
	content, err := os.ReadFile(factoryfilePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != customContent {
		t.Error("Upgrade should not modify existing Factoryfile")
	}
}
