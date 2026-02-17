package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpgrade_FailsIfOpenspecNotInstalled(t *testing.T) {
	dir := t.TempDir()

	// Create Factoryfile so we get past that check.
	if err := os.WriteFile(filepath.Join(dir, "Factoryfile"), []byte(DefaultFactoryfile), 0o644); err != nil {
		t.Fatal(err)
	}

	// Set PATH to an empty directory so openspec cannot be found.
	emptyDir := t.TempDir()
	t.Setenv("PATH", emptyDir)

	err := Upgrade(dir)
	if err == nil {
		t.Fatal("expected Upgrade to fail when openspec is not in PATH")
	}
	if !strings.Contains(err.Error(), "openspec is not installed") {
		t.Errorf("expected 'openspec is not installed' error, got: %v", err)
	}

	// Verify no new files were created (Factoryfile was pre-existing).
	if _, err := os.Stat(filepath.Join(dir, "AGENTS.md")); err == nil {
		t.Error("AGENTS.md should not have been created when openspec check fails")
	}
}

func TestUpgrade_FailsWithoutFactoryfile(t *testing.T) {
	requireOpenSpec(t)
	dir := t.TempDir()

	err := Upgrade(dir)
	if err == nil {
		t.Fatal("expected Upgrade to fail without Factoryfile")
	}
	if !strings.Contains(err.Error(), "No Factoryfile found") {
		t.Errorf("expected 'No Factoryfile found' error, got: %v", err)
	}
}

func TestUpgrade_WithFactoryfile(t *testing.T) {
	requireOpenSpec(t)
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

	// Verify OpenSpec schema was extracted
	schemaDir := filepath.Join(dir, "openspec", "schemas", "littlefactory")
	if _, err := os.Stat(schemaDir); err != nil {
		t.Fatalf("expected openspec schema directory to exist: %v", err)
	}

	// Verify OpenSpec config was created
	configPath := filepath.Join(dir, "openspec", "config.yaml")
	configContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("expected openspec config to exist: %v", err)
	}
	if string(configContent) != "schema: littlefactory\n" {
		t.Errorf("expected default openspec config, got: %s", string(configContent))
	}
}

func TestUpgrade_WithFactoryfileYAML(t *testing.T) {
	requireOpenSpec(t)
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
	requireOpenSpec(t)
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
	requireOpenSpec(t)
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
	requireOpenSpec(t)
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

	// No embedded skills exist, so no symlinks should be created
	claudeSkillsPath := filepath.Join(dir, ".claude", "skills")
	if _, err := os.Stat(claudeSkillsPath); err == nil {
		entries, readErr := os.ReadDir(claudeSkillsPath)
		if readErr != nil {
			t.Fatalf("failed to read .claude/skills/: %v", readErr)
		}
		if len(entries) != 0 {
			t.Errorf("expected no skill symlinks, got %d", len(entries))
		}
	}
}

func TestUpgrade_DoesNotCreateFactoryfile(t *testing.T) {
	requireOpenSpec(t)
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
