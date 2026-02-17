package worktree

import (
	"fmt"
	"os/exec"
	"strings"
)

// Worktree represents a single git worktree entry.
type Worktree struct {
	// Path is the absolute filesystem path of the worktree.
	Path string
	// Commit is the HEAD commit hash of the worktree.
	Commit string
	// Branch is the branch name (e.g., "refs/heads/main"), or empty if detached.
	Branch string
	// IsBare is true if this is a bare repository entry.
	IsBare bool
	// IsDetached is true if HEAD is detached.
	IsDetached bool
}

// BranchShort returns the short branch name (e.g., "main" from "refs/heads/main").
// Returns an empty string if detached or bare.
func (w Worktree) BranchShort() string {
	const prefix = "refs/heads/"
	if strings.HasPrefix(w.Branch, prefix) {
		return w.Branch[len(prefix):]
	}
	return w.Branch
}

// List returns all worktrees for the repository at the given path.
// It parses the output of `git worktree list --porcelain`.
func List(repoDir string) ([]Worktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("listing worktrees: %w", err)
	}

	return parseWorktreeList(string(out)), nil
}

// parseWorktreeList parses the porcelain output of `git worktree list`.
// Each worktree block is separated by a blank line and has the format:
//
//	worktree /path/to/worktree
//	HEAD <commit>
//	branch refs/heads/main
//
// or for detached HEAD:
//
//	worktree /path/to/worktree
//	HEAD <commit>
//	detached
//
// or for bare repos:
//
//	worktree /path/to/repo
//	bare
func parseWorktreeList(output string) []Worktree {
	var worktrees []Worktree
	var current *Worktree

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimRight(line, "\r")

		if line == "" {
			if current != nil {
				worktrees = append(worktrees, *current)
				current = nil
			}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			current = &Worktree{
				Path: line[len("worktree "):],
			}
		} else if current != nil {
			switch {
			case strings.HasPrefix(line, "HEAD "):
				current.Commit = line[len("HEAD "):]
			case strings.HasPrefix(line, "branch "):
				current.Branch = line[len("branch "):]
			case line == "detached":
				current.IsDetached = true
			case line == "bare":
				current.IsBare = true
			}
		}
	}

	// Handle last entry if output doesn't end with a blank line.
	if current != nil {
		worktrees = append(worktrees, *current)
	}

	return worktrees
}
