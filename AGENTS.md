# Agent Instructions

This project uses **littlefactory** for task management. Tasks are stored in `.littlefactory/tasks.json`.

## Quick Reference

Tasks are managed automatically by the littlefactory driver. Manual task management is not typically needed, but the JSON format is:
- `status: "todo"` - Available for work
- `status: "in_progress"` - Currently being worked on
- `status: "done"` - Completed

## Codebase Patterns

- **Init sub-packages** (`internal/init/`): Each sub-package (agentsmd, gitignore, skills) is self-contained with its own types, functions, and tests. Follow the pattern in `internal/init/skills/embed.go` for reference.
- **Symlink handling**: Use `os.Lstat` (not `os.Stat`) when checking for symlinks. `os.Stat` follows symlinks and returns the target's info. Use relative symlink targets for same-directory files (e.g., `os.Symlink("AGENTS.md", claudePath)`). For cross-directory symlinks, use `filepath.Join("..", "..", ...)` relative paths. On macOS, `/var` resolves to `/private/var` via `filepath.EvalSymlinks` -- tests comparing resolved paths must eval both sides.
- **Embedded files**: Use `//go:embed all:embedded/...` for directory trees with `embed.FS`, and `fs.Sub` to strip the embed prefix before walking.
- **Git command packages** (`internal/worktree/`): Packages that shell out to git use `exec.Command` with `cmd.Dir` set to the target directory. Tests create real git repos in `t.TempDir()` with helper functions (`initGitRepo`, `run`). Always resolve paths with `filepath.EvalSymlinks` on both sides when comparing paths in tests.
