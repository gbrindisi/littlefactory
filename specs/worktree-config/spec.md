# worktree-config

## What It Does
Configures the directory where littlefactory creates git worktrees for change isolation. Supports default, custom, relative, and absolute paths via the Factoryfile.

## Requirements

### Requirement: Worktrees directory configuration
The system SHALL allow configuration of the directory used for git worktrees, with a sensible default and support for custom paths.

#### Scenario: Default worktrees directory
- **WHEN** no `worktrees_dir` is specified in the Factoryfile
- **THEN** the system uses `.littlefactory/worktrees/` as the worktrees directory

#### Scenario: Custom worktrees directory
- **WHEN** the user sets `worktrees_dir` to a custom value in the Factoryfile
- **THEN** the system uses that value as the worktrees directory

#### Scenario: Relative path resolution
- **WHEN** `worktrees_dir` is set to a relative path (e.g., `../worktrees`)
- **THEN** the system resolves the path relative to the project root (where the Factoryfile lives)

#### Scenario: Absolute path resolution
- **WHEN** `worktrees_dir` is set to an absolute path (e.g., `/tmp/lf-worktrees`)
- **THEN** the system uses the absolute path as-is

## Boundaries

## Gotchas
