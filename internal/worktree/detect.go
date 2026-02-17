// Package worktree provides git worktree detection and management operations.
package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetCommonDir returns the git common directory for the repository at the given path.
// It uses `git rev-parse --git-common-dir` which works for normal repos, bare repos,
// linked worktrees, and custom GIT_DIR setups.
func GetCommonDir(repoDir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-common-dir")
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("getting git common dir: %w", err)
	}

	commonDir := strings.TrimSpace(string(out))

	// git rev-parse may return a relative path; resolve it relative to repoDir.
	if !filepath.IsAbs(commonDir) {
		commonDir = filepath.Join(repoDir, commonDir)
	}

	// Clean the path to resolve any ".." components.
	commonDir = filepath.Clean(commonDir)

	return commonDir, nil
}

// HasWorktrees checks if the repository at the given path has any existing worktrees.
// It looks for a non-empty `worktrees` directory inside the git common directory.
func HasWorktrees(repoDir string) (bool, error) {
	commonDir, err := GetCommonDir(repoDir)
	if err != nil {
		return false, err
	}

	worktreesDir := filepath.Join(commonDir, "worktrees")
	entries, err := os.ReadDir(worktreesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("reading worktrees directory: %w", err)
	}

	return len(entries) > 0, nil
}
