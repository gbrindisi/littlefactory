package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	if err := Run(dir); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Verify Factoryfile was created
	factoryfilePath := filepath.Join(dir, "Factoryfile")
	content, err := os.ReadFile(factoryfilePath)
	if err != nil {
		t.Fatalf("expected Factoryfile to exist: %v", err)
	}
	if string(content) != DefaultFactoryfile {
		t.Error("Factoryfile content does not match default")
	}

	// Verify AGENTS.md was created
	agentsPath := filepath.Join(dir, "AGENTS.md")
	if _, err := os.Stat(agentsPath); err != nil {
		t.Fatalf("expected AGENTS.md to exist: %v", err)
	}

	// Verify CLAUDE.md is a symlink to AGENTS.md
	claudePath := filepath.Join(dir, "CLAUDE.md")
	target, err := os.Readlink(claudePath)
	if err != nil {
		t.Fatalf("expected CLAUDE.md to be a symlink: %v", err)
	}
	if target != "AGENTS.md" {
		t.Errorf("expected CLAUDE.md -> AGENTS.md, got -> %s", target)
	}

	// Verify .gitignore was created with required entries
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

func TestRun_WithExistingClaudeMD(t *testing.T) {
	dir := t.TempDir()

	// Create existing CLAUDE.md
	claudePath := filepath.Join(dir, "CLAUDE.md")
	originalContent := "# My Custom Instructions\n\nDo things my way.\n"
	if err := os.WriteFile(claudePath, []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Run(dir); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// CLAUDE.md should have been migrated to AGENTS.md
	agentsPath := filepath.Join(dir, "AGENTS.md")
	content, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("expected AGENTS.md to exist: %v", err)
	}
	if string(content) != originalContent {
		t.Errorf("AGENTS.md should contain original CLAUDE.md content")
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

func TestRun_FailsIfFactoryfileExists(t *testing.T) {
	dir := t.TempDir()

	// Create existing Factoryfile
	factoryfilePath := filepath.Join(dir, "Factoryfile")
	if err := os.WriteFile(factoryfilePath, []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Run(dir)
	if err == nil {
		t.Fatal("expected Run to fail when Factoryfile already exists")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' error, got: %v", err)
	}
}

func TestRun_FailsIfFactoryfileYAMLExists(t *testing.T) {
	dir := t.TempDir()

	// Create existing Factoryfile.yaml
	factoryfilePath := filepath.Join(dir, "Factoryfile.yaml")
	if err := os.WriteFile(factoryfilePath, []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Run(dir)
	if err == nil {
		t.Fatal("expected Run to fail when Factoryfile.yaml already exists")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' error, got: %v", err)
	}
}

func TestRun_WithExistingGitignore(t *testing.T) {
	dir := t.TempDir()

	// Create existing .gitignore with custom content
	gitignorePath := filepath.Join(dir, ".gitignore")
	existingContent := "node_modules/\n*.log\n"
	if err := os.WriteFile(gitignorePath, []byte(existingContent), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Run(dir); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Existing content should be preserved
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}
	if !strings.HasPrefix(string(content), existingContent) {
		t.Error("existing .gitignore content was not preserved")
	}
	if !strings.Contains(string(content), ".littlefactory/run_metadata.json") {
		t.Error(".gitignore missing .littlefactory/run_metadata.json")
	}
	if !strings.Contains(string(content), ".littlefactory/tasks.json") {
		t.Error(".gitignore missing .littlefactory/tasks.json")
	}
}

func TestRun_WithoutClaudeDir_SkipsSymlinks(t *testing.T) {
	dir := t.TempDir()

	// No .claude/ directory exists
	if err := Run(dir); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// .claude/skills/ should NOT exist
	claudeSkillsPath := filepath.Join(dir, ".claude", "skills")
	if _, err := os.Stat(claudeSkillsPath); err == nil {
		t.Error("expected .claude/skills/ to NOT exist when .claude/ is absent")
	}
}
