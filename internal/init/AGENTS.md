# Init Package

## Codebase Patterns

- **Init sub-packages** (`internal/init/`): Each sub-package (agentsmd, gitignore, skills) is self-contained with its own types, functions, and tests. Follow the pattern in `internal/init/skills/embed.go` for reference.
- **Symlink handling**: Use `os.Lstat` (not `os.Stat`) when checking for symlinks. `os.Stat` follows symlinks and returns the target's info. Use relative symlink targets for same-directory files (e.g., `os.Symlink("AGENTS.md", claudePath)`). For cross-directory symlinks, use `filepath.Join("..", "..", ...)` relative paths. On macOS, `/var` resolves to `/private/var` via `filepath.EvalSymlinks` -- tests comparing resolved paths must eval both sides.
- **Embedded files**: Use `//go:embed all:embedded/...` for directory trees with `embed.FS`, and `fs.Sub` to strip the embed prefix before walking.
