# Worktree Package

## Codebase Patterns

- **Git command packages** (`internal/worktree/`): Packages that shell out to git use `exec.Command` with `cmd.Dir` set to the target directory. Tests create real git repos in `t.TempDir()` with helper functions (`initGitRepo`, `run`). Always resolve paths with `filepath.EvalSymlinks` on both sides when comparing paths in tests.
