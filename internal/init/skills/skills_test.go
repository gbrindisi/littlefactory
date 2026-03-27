package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateSymlinks_NoClaudeDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .littlefactory/skills/ with a skill but no .claude/ dir
	skillDir := filepath.Join(tmpDir, ".littlefactory", "skills", "test-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := CreateSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("CreateSymlinks failed: %v", err)
	}

	if result.ClaudeDirExists {
		t.Fatal("expected ClaudeDirExists=false when .claude/ does not exist")
	}
	if len(result.Entries) != 0 {
		t.Fatalf("expected no entries, got %d", len(result.Entries))
	}
}

func TestCreateSymlinks_CreatesSymlinks(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .claude/ directory
	if err := os.MkdirAll(filepath.Join(tmpDir, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Create .littlefactory/skills/ with a skill
	skillDir := filepath.Join(tmpDir, ".littlefactory", "skills", "test-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := CreateSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("CreateSymlinks failed: %v", err)
	}

	if !result.ClaudeDirExists {
		t.Fatal("expected ClaudeDirExists=true")
	}
	if len(result.Created()) != 1 {
		t.Fatalf("expected 1 created entry, got %d", len(result.Created()))
	}
	if result.Created()[0].Name != "test-skill" {
		t.Fatalf("expected skill name 'test-skill', got %q", result.Created()[0].Name)
	}

	// Verify symlink exists and points to correct target
	linkPath := filepath.Join(tmpDir, ".claude", "skills", "test-skill")
	target, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("expected symlink at %s: %v", linkPath, err)
	}

	expectedTarget := filepath.Join("..", "..", ".littlefactory", "skills", "test-skill")
	if target != expectedTarget {
		t.Fatalf("expected symlink target %q, got %q", expectedTarget, target)
	}

	// Verify the symlink resolves to the actual skill directory.
	// Use EvalSymlinks on both to handle macOS /var -> /private/var.
	resolved, err := filepath.EvalSymlinks(linkPath)
	if err != nil {
		t.Fatalf("symlink does not resolve: %v", err)
	}
	expectedResolved, err := filepath.EvalSymlinks(filepath.Join(tmpDir, ".littlefactory", "skills", "test-skill"))
	if err != nil {
		t.Fatalf("could not resolve expected path: %v", err)
	}
	if resolved != expectedResolved {
		t.Fatalf("expected resolved path %q, got %q", expectedResolved, resolved)
	}
}

func TestCreateSymlinks_CreatesClaudeSkillsDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .claude/ but NOT .claude/skills/
	if err := os.MkdirAll(filepath.Join(tmpDir, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Create a skill
	skillDir := filepath.Join(tmpDir, ".littlefactory", "skills", "my-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := CreateSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("CreateSymlinks failed: %v", err)
	}

	// Verify .claude/skills/ was created
	info, err := os.Stat(filepath.Join(tmpDir, ".claude", "skills"))
	if err != nil {
		t.Fatalf("expected .claude/skills/ to exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected .claude/skills/ to be a directory")
	}
}

func TestCreateSymlinks_SkipsExistingSymlink(t *testing.T) {
	tmpDir := t.TempDir()

	// Set up directories
	if err := os.MkdirAll(filepath.Join(tmpDir, ".claude", "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	skillDir := filepath.Join(tmpDir, ".littlefactory", "skills", "test-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create an existing symlink at the target location
	linkPath := filepath.Join(tmpDir, ".claude", "skills", "test-skill")
	if err := os.Symlink("/some/other/target", linkPath); err != nil {
		t.Fatal(err)
	}

	result, err := CreateSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("CreateSymlinks failed: %v", err)
	}

	if len(result.Skipped()) != 1 {
		t.Fatalf("expected 1 skipped entry, got %d", len(result.Skipped()))
	}
	if len(result.Created()) != 0 {
		t.Fatalf("expected 0 created entries, got %d", len(result.Created()))
	}

	// Verify the existing symlink was NOT overwritten
	target, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatal(err)
	}
	if target != "/some/other/target" {
		t.Fatalf("symlink was overwritten: got target %q", target)
	}
}

func TestCreateSymlinks_SkipsExistingFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Set up directories
	if err := os.MkdirAll(filepath.Join(tmpDir, ".claude", "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	skillDir := filepath.Join(tmpDir, ".littlefactory", "skills", "test-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create a regular directory (not a symlink) at the target location
	if err := os.MkdirAll(filepath.Join(tmpDir, ".claude", "skills", "test-skill"), 0o755); err != nil {
		t.Fatal(err)
	}

	result, err := CreateSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("CreateSymlinks failed: %v", err)
	}

	if len(result.Skipped()) != 1 {
		t.Fatalf("expected 1 skipped entry, got %d", len(result.Skipped()))
	}
}

func TestCreateSymlinks_MultipleSkills(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .claude/
	if err := os.MkdirAll(filepath.Join(tmpDir, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Create multiple skills
	for _, name := range []string{"skill-a", "skill-b", "skill-c"} {
		skillDir := filepath.Join(tmpDir, ".littlefactory", "skills", name)
		if err := os.MkdirAll(skillDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("content"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	result, err := CreateSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("CreateSymlinks failed: %v", err)
	}

	if len(result.Created()) != 3 {
		t.Fatalf("expected 3 created entries, got %d", len(result.Created()))
	}

	// Verify each symlink
	for _, name := range []string{"skill-a", "skill-b", "skill-c"} {
		linkPath := filepath.Join(tmpDir, ".claude", "skills", name)
		target, err := os.Readlink(linkPath)
		if err != nil {
			t.Fatalf("expected symlink for %s: %v", name, err)
		}
		expectedTarget := filepath.Join("..", "..", ".littlefactory", "skills", name)
		if target != expectedTarget {
			t.Fatalf("skill %s: expected target %q, got %q", name, expectedTarget, target)
		}
	}
}

func TestCreateSymlinks_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .claude/
	if err := os.MkdirAll(filepath.Join(tmpDir, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Create a skill
	skillDir := filepath.Join(tmpDir, ".littlefactory", "skills", "test-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}

	// First call: should create
	result1, err := CreateSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("first CreateSymlinks failed: %v", err)
	}
	if len(result1.Created()) != 1 {
		t.Fatalf("first call: expected 1 created, got %d", len(result1.Created()))
	}

	// Second call: should skip
	result2, err := CreateSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("second CreateSymlinks failed: %v", err)
	}
	if len(result2.Skipped()) != 1 {
		t.Fatalf("second call: expected 1 skipped, got %d", len(result2.Skipped()))
	}
	if len(result2.Created()) != 0 {
		t.Fatalf("second call: expected 0 created, got %d", len(result2.Created()))
	}
}

func TestCreateSymlinks_NoSkillsDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .claude/ but no .littlefactory/skills/
	if err := os.MkdirAll(filepath.Join(tmpDir, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}

	result, err := CreateSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("CreateSymlinks failed: %v", err)
	}

	if !result.ClaudeDirExists {
		t.Fatal("expected ClaudeDirExists=true")
	}
	if len(result.Entries) != 0 {
		t.Fatalf("expected no entries, got %d", len(result.Entries))
	}
}

func TestCleanupOrphanedSymlinks_RemovesOrphans(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .claude/skills/ with openspec-* symlinks (orphans) and a legit lf-* symlink
	claudeSkills := filepath.Join(tmpDir, ".claude", "skills")
	if err := os.MkdirAll(claudeSkills, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"openspec-explore", "openspec-new-change", "openspec-archive-change"} {
		if err := os.Symlink("/nonexistent", filepath.Join(claudeSkills, name)); err != nil {
			t.Fatal(err)
		}
	}
	// A non-openspec symlink should be left alone
	if err := os.Symlink("/other", filepath.Join(claudeSkills, "lf-explore")); err != nil {
		t.Fatal(err)
	}

	removed, err := CleanupOrphanedSymlinks(tmpDir, "openspec-")
	if err != nil {
		t.Fatalf("CleanupOrphanedSymlinks failed: %v", err)
	}

	if len(removed) != 3 {
		t.Fatalf("expected 3 removed, got %d: %v", len(removed), removed)
	}

	// lf-explore should still exist
	if _, err := os.Lstat(filepath.Join(claudeSkills, "lf-explore")); err != nil {
		t.Fatal("lf-explore symlink was removed but should not have been")
	}
}

func TestCleanupOrphanedSymlinks_KeepsMatchingSkills(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .claude/skills/ with an openspec-* symlink
	claudeSkills := filepath.Join(tmpDir, ".claude", "skills")
	if err := os.MkdirAll(claudeSkills, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("/nonexistent", filepath.Join(claudeSkills, "openspec-explore")); err != nil {
		t.Fatal(err)
	}

	// Create a corresponding .littlefactory/skills/openspec-explore/ (not orphaned)
	if err := os.MkdirAll(filepath.Join(tmpDir, ".littlefactory", "skills", "openspec-explore"), 0o755); err != nil {
		t.Fatal(err)
	}

	removed, err := CleanupOrphanedSymlinks(tmpDir, "openspec-")
	if err != nil {
		t.Fatalf("CleanupOrphanedSymlinks failed: %v", err)
	}

	if len(removed) != 0 {
		t.Fatalf("expected 0 removed (skill exists in .littlefactory/skills/), got %d", len(removed))
	}
}

func TestCleanupOrphanedSymlinks_SkipsNonSymlinks(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a real directory (not a symlink) named openspec-something
	claudeSkills := filepath.Join(tmpDir, ".claude", "skills")
	if err := os.MkdirAll(filepath.Join(claudeSkills, "openspec-real-dir"), 0o755); err != nil {
		t.Fatal(err)
	}

	removed, err := CleanupOrphanedSymlinks(tmpDir, "openspec-")
	if err != nil {
		t.Fatalf("CleanupOrphanedSymlinks failed: %v", err)
	}

	if len(removed) != 0 {
		t.Fatalf("expected 0 removed (real directory, not symlink), got %d", len(removed))
	}

	// Directory should still exist
	if _, err := os.Stat(filepath.Join(claudeSkills, "openspec-real-dir")); err != nil {
		t.Fatal("real directory was removed but should not have been")
	}
}

func TestCleanupOrphanedSymlinks_NoClaudeSkillsDir(t *testing.T) {
	tmpDir := t.TempDir()

	removed, err := CleanupOrphanedSymlinks(tmpDir, "openspec-")
	if err != nil {
		t.Fatalf("CleanupOrphanedSymlinks failed: %v", err)
	}

	if len(removed) != 0 {
		t.Fatalf("expected 0 removed, got %d", len(removed))
	}
}

func TestCreateSymlinks_IgnoresFilesInSkillsDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .claude/
	if err := os.MkdirAll(filepath.Join(tmpDir, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Create a skill directory and a loose file in skills/
	skillDir := filepath.Join(tmpDir, ".littlefactory", "skills", "my-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Loose file (not a skill directory) should be ignored
	if err := os.WriteFile(filepath.Join(tmpDir, ".littlefactory", "skills", "README.md"), []byte("info"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := CreateSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("CreateSymlinks failed: %v", err)
	}

	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry (only directories), got %d", len(result.Entries))
	}
	if result.Entries[0].Name != "my-skill" {
		t.Fatalf("expected skill name 'my-skill', got %q", result.Entries[0].Name)
	}
}
