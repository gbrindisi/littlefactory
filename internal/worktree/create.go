package worktree

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Create creates a new git worktree for the given branch name.
// The worktree is placed at <worktreesDir>/<branchName> and a new branch
// is created from the current HEAD.
//
// It does NOT check for clean working tree or existing worktrees; callers
// should use IsClean() and WorktreeExists() before calling Create().
func Create(repoDir, branchName, worktreesDir string) (string, error) {
	worktreePath := filepath.Join(worktreesDir, branchName)

	cmd := exec.Command("git", "worktree", "add", worktreePath, "-b", branchName)
	cmd.Dir = repoDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("creating worktree: %s: %w", strings.TrimSpace(string(out)), err)
	}

	return worktreePath, nil
}

// IsClean checks if the working tree at the given path has no uncommitted changes.
// It uses `git status --porcelain` which outputs nothing for a clean tree.
func IsClean(repoDir string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("checking git status: %w", err)
	}

	return strings.TrimSpace(string(out)) == "", nil
}

// WorktreeExists checks if a worktree already exists for the given branch name.
// If found, it returns true and the path to the existing worktree.
func WorktreeExists(repoDir, branchName string) (bool, string, error) {
	worktrees, err := List(repoDir)
	if err != nil {
		return false, "", err
	}

	for _, wt := range worktrees {
		if wt.BranchShort() == branchName {
			return true, wt.Path, nil
		}
	}

	return false, "", nil
}
