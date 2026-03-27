// Package skills handles skill extraction from embedded files
// and symlinking into the .claude/skills/ directory.
package skills

import (
	"fmt"
	"os"
	"path/filepath"
)

// SymlinkAction describes what happened to an individual skill symlink.
type SymlinkAction string

const (
	// SymlinkCreated indicates a new symlink was created.
	SymlinkCreated SymlinkAction = "created"
	// SymlinkSkipped indicates the target already existed.
	SymlinkSkipped SymlinkAction = "skipped"
)

// SymlinkEntry records the action taken for a single skill.
type SymlinkEntry struct {
	Name   string
	Action SymlinkAction
}

// SymlinkResult describes the outcome of a CreateSymlinks call.
type SymlinkResult struct {
	// ClaudeDirExists is false when .claude/ was not found (no symlinks created).
	ClaudeDirExists bool
	Entries         []SymlinkEntry
}

// Created returns the entries that were newly created.
func (r SymlinkResult) Created() []SymlinkEntry {
	var out []SymlinkEntry
	for _, e := range r.Entries {
		if e.Action == SymlinkCreated {
			out = append(out, e)
		}
	}
	return out
}

// Skipped returns the entries that were skipped because they already existed.
func (r SymlinkResult) Skipped() []SymlinkEntry {
	var out []SymlinkEntry
	for _, e := range r.Entries {
		if e.Action == SymlinkSkipped {
			out = append(out, e)
		}
	}
	return out
}

// CreateSymlinks creates symlinks in .claude/skills/ pointing to .littlefactory/skills/<name>
// for each skill directory found in .littlefactory/skills/. The symlink direction is:
//
//	.claude/skills/<name> -> ../../.littlefactory/skills/<name>
//
// If .claude/ does not exist, the function returns immediately with ClaudeDirExists=false.
// If .claude/skills/ does not exist, it is created.
// Existing files or symlinks in .claude/skills/ are not overwritten.
func CreateSymlinks(projectRoot string) (SymlinkResult, error) {
	claudeDir := filepath.Join(projectRoot, ".claude")
	if !dirExists(claudeDir) {
		return SymlinkResult{ClaudeDirExists: false}, nil
	}

	skillsSrc := filepath.Join(projectRoot, ".littlefactory", "skills")
	if !dirExists(skillsSrc) {
		return SymlinkResult{ClaudeDirExists: true}, nil
	}

	claudeSkillsDir := filepath.Join(claudeDir, "skills")
	if err := os.MkdirAll(claudeSkillsDir, 0o755); err != nil {
		return SymlinkResult{}, fmt.Errorf("creating .claude/skills/ directory: %w", err)
	}

	entries, err := os.ReadDir(skillsSrc)
	if err != nil {
		return SymlinkResult{}, fmt.Errorf("reading .littlefactory/skills/: %w", err)
	}

	var result SymlinkResult
	result.ClaudeDirExists = true

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		linkPath := filepath.Join(claudeSkillsDir, name)

		// Skip if anything already exists at the link path (file or symlink).
		if pathExists(linkPath) {
			result.Entries = append(result.Entries, SymlinkEntry{
				Name:   name,
				Action: SymlinkSkipped,
			})
			continue
		}

		// Relative symlink: .claude/skills/<name> -> ../../.littlefactory/skills/<name>
		target := filepath.Join("..", "..", ".littlefactory", "skills", name)
		if err := os.Symlink(target, linkPath); err != nil {
			return SymlinkResult{}, fmt.Errorf("creating symlink for skill %s: %w", name, err)
		}

		result.Entries = append(result.Entries, SymlinkEntry{
			Name:   name,
			Action: SymlinkCreated,
		})
	}

	return result, nil
}

// CleanupOrphanedSymlinks removes symlinks from .claude/skills/ whose names
// match the given prefix but do not correspond to any skill in .littlefactory/skills/.
// This is used during upgrade to remove stale openspec-* symlinks.
func CleanupOrphanedSymlinks(projectRoot, prefix string) ([]string, error) {
	claudeSkillsDir := filepath.Join(projectRoot, ".claude", "skills")
	if !dirExists(claudeSkillsDir) {
		return nil, nil
	}

	lfSkillsDir := filepath.Join(projectRoot, ".littlefactory", "skills")

	entries, err := os.ReadDir(claudeSkillsDir)
	if err != nil {
		return nil, fmt.Errorf("reading .claude/skills/: %w", err)
	}

	var removed []string
	for _, entry := range entries {
		name := entry.Name()
		if len(name) <= len(prefix) || name[:len(prefix)] != prefix {
			continue
		}

		linkPath := filepath.Join(claudeSkillsDir, name)

		// Only remove symlinks, not real directories/files.
		fi, err := os.Lstat(linkPath)
		if err != nil || fi.Mode()&os.ModeSymlink == 0 {
			continue
		}

		// Keep if a corresponding skill exists in .littlefactory/skills/.
		if dirExists(filepath.Join(lfSkillsDir, name)) {
			continue
		}

		if err := os.Remove(linkPath); err != nil {
			return removed, fmt.Errorf("removing orphaned symlink %s: %w", name, err)
		}
		removed = append(removed, name)
	}

	return removed, nil
}

// dirExists returns true if the path exists and is a directory.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// pathExists returns true if anything exists at the path (file, directory, or symlink).
func pathExists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}
