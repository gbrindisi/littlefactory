package agentsmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSetup_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	result, err := Setup(dir)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	if result.Action != ActionCreated {
		t.Errorf("expected action %q, got %q", ActionCreated, result.Action)
	}

	// Verify AGENTS.md was created with default content
	agentsPath := filepath.Join(dir, "AGENTS.md")
	content, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("failed to read AGENTS.md: %v", err)
	}
	if string(content) != DefaultContent {
		t.Errorf("AGENTS.md content does not match default content")
	}

	// Verify CLAUDE.md is a symlink to AGENTS.md
	claudePath := filepath.Join(dir, "CLAUDE.md")
	target, err := os.Readlink(claudePath)
	if err != nil {
		t.Fatalf("CLAUDE.md is not a symlink: %v", err)
	}
	if target != "AGENTS.md" {
		t.Errorf("expected symlink target %q, got %q", "AGENTS.md", target)
	}

	// Verify reading through the symlink returns the same content
	claudeContent, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatalf("failed to read through CLAUDE.md symlink: %v", err)
	}
	if string(claudeContent) != DefaultContent {
		t.Errorf("CLAUDE.md symlink content does not match AGENTS.md")
	}
}

func TestSetup_ClaudeMDOnly(t *testing.T) {
	dir := t.TempDir()

	originalContent := "# My Claude Instructions\n\nCustom content here.\n"
	claudePath := filepath.Join(dir, "CLAUDE.md")
	if err := os.WriteFile(claudePath, []byte(originalContent), 0o644); err != nil {
		t.Fatalf("failed to create CLAUDE.md: %v", err)
	}

	result, err := Setup(dir)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	if result.Action != ActionMigrated {
		t.Errorf("expected action %q, got %q", ActionMigrated, result.Action)
	}

	// Verify AGENTS.md has the original CLAUDE.md content
	agentsPath := filepath.Join(dir, "AGENTS.md")
	content, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("failed to read AGENTS.md: %v", err)
	}
	if string(content) != originalContent {
		t.Errorf("AGENTS.md content does not match original CLAUDE.md content")
	}

	// Verify CLAUDE.md is now a symlink
	target, err := os.Readlink(claudePath)
	if err != nil {
		t.Fatalf("CLAUDE.md is not a symlink: %v", err)
	}
	if target != "AGENTS.md" {
		t.Errorf("expected symlink target %q, got %q", "AGENTS.md", target)
	}
}

func TestSetup_BothFilesExist(t *testing.T) {
	dir := t.TempDir()

	agentsContent := "# AGENTS content\n\nAgents stuff.\n"
	claudeContent := "# CLAUDE content\n\nClaude stuff.\n"

	agentsPath := filepath.Join(dir, "AGENTS.md")
	claudePath := filepath.Join(dir, "CLAUDE.md")

	if err := os.WriteFile(agentsPath, []byte(agentsContent), 0o644); err != nil {
		t.Fatalf("failed to create AGENTS.md: %v", err)
	}
	if err := os.WriteFile(claudePath, []byte(claudeContent), 0o644); err != nil {
		t.Fatalf("failed to create CLAUDE.md: %v", err)
	}

	result, err := Setup(dir)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	if result.Action != ActionMerged {
		t.Errorf("expected action %q, got %q", ActionMerged, result.Action)
	}

	// Verify merged content
	merged, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("failed to read merged AGENTS.md: %v", err)
	}

	mergedStr := string(merged)
	if !strings.Contains(mergedStr, agentsContent) {
		t.Error("merged content does not contain original AGENTS.md content")
	}
	if !strings.Contains(mergedStr, claudeContent) {
		t.Error("merged content does not contain original CLAUDE.md content")
	}
	if !strings.Contains(mergedStr, "<!-- Merged from CLAUDE.md -->") {
		t.Error("merged content does not contain merge separator")
	}

	// Verify the exact merge format
	expected := agentsContent + MergeSeparator + claudeContent
	if mergedStr != expected {
		t.Errorf("merged content does not match expected format.\nExpected:\n%s\nGot:\n%s", expected, mergedStr)
	}

	// Verify CLAUDE.md is now a symlink
	target, err := os.Readlink(claudePath)
	if err != nil {
		t.Fatalf("CLAUDE.md is not a symlink: %v", err)
	}
	if target != "AGENTS.md" {
		t.Errorf("expected symlink target %q, got %q", "AGENTS.md", target)
	}
}

func TestSetup_AlreadyConfigured(t *testing.T) {
	dir := t.TempDir()

	agentsPath := filepath.Join(dir, "AGENTS.md")
	claudePath := filepath.Join(dir, "CLAUDE.md")

	content := "# Already set up\n"
	if err := os.WriteFile(agentsPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create AGENTS.md: %v", err)
	}
	if err := os.Symlink("AGENTS.md", claudePath); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	result, err := Setup(dir)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	if result.Action != ActionSkipped {
		t.Errorf("expected action %q, got %q", ActionSkipped, result.Action)
	}

	// Verify AGENTS.md content unchanged
	readContent, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("failed to read AGENTS.md: %v", err)
	}
	if string(readContent) != content {
		t.Errorf("AGENTS.md content was modified when it should have been skipped")
	}
}

func TestSetup_AgentsMDOnlyNoClaudeMD(t *testing.T) {
	dir := t.TempDir()

	agentsPath := filepath.Join(dir, "AGENTS.md")
	content := "# Existing AGENTS.md\n"
	if err := os.WriteFile(agentsPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create AGENTS.md: %v", err)
	}

	result, err := Setup(dir)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	if result.Action != ActionCreated {
		t.Errorf("expected action %q, got %q", ActionCreated, result.Action)
	}

	// Verify CLAUDE.md symlink was created
	claudePath := filepath.Join(dir, "CLAUDE.md")
	target, err := os.Readlink(claudePath)
	if err != nil {
		t.Fatalf("CLAUDE.md is not a symlink: %v", err)
	}
	if target != "AGENTS.md" {
		t.Errorf("expected symlink target %q, got %q", "AGENTS.md", target)
	}

	// Verify AGENTS.md content unchanged
	readContent, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("failed to read AGENTS.md: %v", err)
	}
	if string(readContent) != content {
		t.Errorf("AGENTS.md content was modified when it should not have been")
	}
}

func TestSetup_Idempotent(t *testing.T) {
	dir := t.TempDir()

	// First run: creates both files
	result1, err := Setup(dir)
	if err != nil {
		t.Fatalf("first Setup failed: %v", err)
	}
	if result1.Action != ActionCreated {
		t.Errorf("first run: expected action %q, got %q", ActionCreated, result1.Action)
	}

	// Second run: should skip
	result2, err := Setup(dir)
	if err != nil {
		t.Fatalf("second Setup failed: %v", err)
	}
	if result2.Action != ActionSkipped {
		t.Errorf("second run: expected action %q, got %q", ActionSkipped, result2.Action)
	}
}

func TestDefaultContent(t *testing.T) {
	if !strings.Contains(DefaultContent, "# Agent Instructions") {
		t.Error("DefaultContent missing header")
	}
	if !strings.Contains(DefaultContent, "littlefactory") {
		t.Error("DefaultContent missing littlefactory reference")
	}
	if !strings.Contains(DefaultContent, ".littlefactory/tasks.json") {
		t.Error("DefaultContent missing tasks.json reference")
	}
}

func TestMergeSeparator(t *testing.T) {
	if !strings.Contains(MergeSeparator, "---") {
		t.Error("MergeSeparator missing horizontal rule")
	}
	if !strings.Contains(MergeSeparator, "Merged from CLAUDE.md") {
		t.Error("MergeSeparator missing merge attribution")
	}
}
