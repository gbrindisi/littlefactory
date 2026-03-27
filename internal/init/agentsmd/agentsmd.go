// Package agentsmd handles AGENTS.md creation and management,
// including symlink setup for CLAUDE.md compatibility.
package agentsmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// DefaultContent is the default content for a new AGENTS.md file.
const DefaultContent = `# Agent Instructions

This project uses **littlefactory** for task management. Tasks are stored in ` + "`" + `.littlefactory/tasks.json` + "`" + `.

## Quick Reference

Tasks are managed automatically by the littlefactory driver. Manual task management is not typically needed, but the JSON format is:
- ` + "`" + `status: "todo"` + "`" + ` - Available for work
- ` + "`" + `status: "in_progress"` + "`" + ` - Currently being worked on
- ` + "`" + `status: "done"` + "`" + ` - Completed
`

// MergeSeparator is the separator used when merging AGENTS.md and CLAUDE.md content.
const MergeSeparator = "\n\n---\n<!-- Merged from CLAUDE.md -->\n\n"

// Action describes what the Setup function did.
type Action string

const (
	// ActionCreated indicates a new AGENTS.md was created with default content.
	ActionCreated Action = "created"
	// ActionMigrated indicates CLAUDE.md was renamed to AGENTS.md and symlinked.
	ActionMigrated Action = "migrated"
	// ActionMerged indicates both files were merged into AGENTS.md and CLAUDE.md was symlinked.
	ActionMerged Action = "merged"
	// ActionSkipped indicates CLAUDE.md is already a symlink to AGENTS.md.
	ActionSkipped Action = "skipped"
)

// Result describes the outcome of a Setup call.
type Result struct {
	Action  Action
	Message string
}

// Setup handles AGENTS.md creation and CLAUDE.md symlink management.
// It detects the current state and applies the appropriate action:
//   - No files exist: creates AGENTS.md with default content and symlinks CLAUDE.md
//   - Only CLAUDE.md exists: renames it to AGENTS.md and symlinks CLAUDE.md
//   - Both files exist: merges content into AGENTS.md and symlinks CLAUDE.md
//   - CLAUDE.md is already a symlink to AGENTS.md: skips (already configured)
func Setup(projectRoot string) (Result, error) {
	agentsPath := filepath.Join(projectRoot, "AGENTS.md")
	claudePath := filepath.Join(projectRoot, "CLAUDE.md")

	agentsExists := fileExists(agentsPath)
	claudeExists := fileExists(claudePath)
	claudeIsSymlink := isSymlink(claudePath)

	// Already configured: CLAUDE.md is a symlink (pointing to AGENTS.md)
	if claudeIsSymlink {
		return Result{
			Action:  ActionSkipped,
			Message: "already configured, CLAUDE.md is a symlink",
		}, nil
	}

	switch {
	case !agentsExists && !claudeExists:
		return handleCreate(agentsPath, claudePath)
	case !agentsExists && claudeExists:
		return handleMigrate(agentsPath, claudePath)
	case agentsExists && claudeExists:
		return handleMerge(agentsPath, claudePath)
	default:
		// AGENTS.md exists, no CLAUDE.md: just create the symlink
		return handleSymlinkOnly(agentsPath, claudePath)
	}
}

// handleCreate creates a new AGENTS.md with default content and symlinks CLAUDE.md.
func handleCreate(agentsPath, claudePath string) (Result, error) {
	if err := os.WriteFile(agentsPath, []byte(DefaultContent), 0o644); err != nil {
		return Result{}, fmt.Errorf("creating AGENTS.md: %w", err)
	}

	if err := os.Symlink("AGENTS.md", claudePath); err != nil {
		return Result{}, fmt.Errorf("creating CLAUDE.md symlink: %w", err)
	}

	return Result{
		Action:  ActionCreated,
		Message: "created AGENTS.md with default content",
	}, nil
}

// handleMigrate renames CLAUDE.md to AGENTS.md and creates a symlink.
func handleMigrate(agentsPath, claudePath string) (Result, error) {
	if err := os.Rename(claudePath, agentsPath); err != nil {
		return Result{}, fmt.Errorf("renaming CLAUDE.md to AGENTS.md: %w", err)
	}

	if err := os.Symlink("AGENTS.md", claudePath); err != nil {
		return Result{}, fmt.Errorf("creating CLAUDE.md symlink: %w", err)
	}

	return Result{
		Action:  ActionMigrated,
		Message: "migrated CLAUDE.md to AGENTS.md",
	}, nil
}

// handleMerge merges AGENTS.md and CLAUDE.md content, then symlinks CLAUDE.md.
func handleMerge(agentsPath, claudePath string) (Result, error) {
	agentsContent, err := os.ReadFile(agentsPath)
	if err != nil {
		return Result{}, fmt.Errorf("reading AGENTS.md: %w", err)
	}

	claudeContent, err := os.ReadFile(claudePath)
	if err != nil {
		return Result{}, fmt.Errorf("reading CLAUDE.md: %w", err)
	}

	merged := string(agentsContent) + MergeSeparator + string(claudeContent)

	if err := os.WriteFile(agentsPath, []byte(merged), 0o644); err != nil { // #nosec G703 -- agentsPath is constructed from projectRoot + constant filename
		return Result{}, fmt.Errorf("writing merged AGENTS.md: %w", err)
	}

	if err := os.Remove(claudePath); err != nil {
		return Result{}, fmt.Errorf("removing CLAUDE.md: %w", err)
	}

	if err := os.Symlink("AGENTS.md", claudePath); err != nil {
		return Result{}, fmt.Errorf("creating CLAUDE.md symlink: %w", err)
	}

	return Result{
		Action:  ActionMerged,
		Message: "merged CLAUDE.md into AGENTS.md",
	}, nil
}

// handleSymlinkOnly creates a CLAUDE.md symlink when only AGENTS.md exists.
func handleSymlinkOnly(agentsPath, claudePath string) (Result, error) {
	if err := os.Symlink("AGENTS.md", claudePath); err != nil {
		return Result{}, fmt.Errorf("creating CLAUDE.md symlink: %w", err)
	}

	return Result{
		Action:  ActionCreated,
		Message: "created CLAUDE.md symlink to existing AGENTS.md",
	}, nil
}

// fileExists returns true if the path exists and is not a directory.
func fileExists(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// isSymlink returns true if the path exists and is a symbolic link.
func isSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}
