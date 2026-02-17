## Why

Littlefactory currently runs in a single working directory, limiting work to one change at a time. To parallelize implementation of multiple openspec changes, we need git worktree support so each change can run in an isolated workspace with its own branch.

## What Changes

- Add `--change` / `-c` flag to `run` command to specify which openspec change's tasks.json to use
- Add `--worktree` / `-w` flag to `run` command to create a new git worktree for the change
- Add `worktrees_dir` configuration option in Factoryfile for custom worktree locations
- Add `status` command to show progress of changes (tasks done/total) across worktrees
- Detect worktree-enabled repos via `git rev-parse --git-common-dir`

## Capabilities

### New Capabilities

- `change-flag`: The `--change` / `-c` flag for specifying which openspec change to use as the task source
- `worktree-creation`: The `--worktree` / `-w` flag for creating isolated git worktrees per change
- `worktree-config`: Configuration option `worktrees_dir` in Factoryfile for worktree location
- `status-command`: Command to display task progress across all changes/worktrees

### Modified Capabilities

- `config-management`: Add `worktrees_dir` option to Factoryfile schema
- `loop-driver`: Support running in a worktree directory with change-specific task source

## Impact

- `cmd/littlefactory/`: New flags on `run` command, new `status` subcommand
- `internal/config/`: Extended Factoryfile parsing for `worktrees_dir`
- `internal/driver/`: Modified to accept change name and resolve task source path
- `internal/worktree/`: New package for git worktree detection and creation
- Users can now run multiple implementations in parallel using separate worktrees
