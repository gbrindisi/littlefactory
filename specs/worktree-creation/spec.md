# worktree-creation

## What It Does
Manages the creation of git worktrees for change isolation. Worktrees allow littlefactory to execute changes in isolated branches without affecting the user's working tree. The system validates preconditions, creates worktrees via git, and detects existing worktrees.

## Requirements

### Requirement: Worktree flag on run command
The system SHALL support a worktree flag on the run command that creates a git worktree for the change before executing tasks.

#### Scenario: Run with worktree flag creates worktree
- **WHEN** the user runs `littlefactory run -c <name> --worktree`
- **THEN** the system creates a git worktree in the configured worktrees directory and executes tasks within it

#### Scenario: Worktree flag requires change flag
- **WHEN** the user runs `littlefactory run --worktree` without specifying a change
- **THEN** the system exits with an error indicating that the worktree flag requires a change name

### Requirement: Worktree creation requires clean working tree
The system SHALL require a clean git working tree before creating a worktree to prevent data loss.

#### Scenario: Dirty working tree blocks worktree creation
- **WHEN** the user runs with `--worktree` and the working tree has uncommitted changes
- **THEN** the system exits with an error indicating the working tree must be clean

#### Scenario: Clean working tree allows worktree creation
- **WHEN** the user runs with `--worktree` and the working tree is clean
- **THEN** the system proceeds with worktree creation

### Requirement: Worktree flag reuses existing worktree
The system SHALL reuse an existing worktree for the given change name instead of erroring, logging that it is reusing the existing worktree.

#### Scenario: Existing worktree is reused
- **WHEN** the user runs with `--worktree` and a worktree for that change name already exists
- **THEN** the system logs "Reusing existing worktree at <path>" and proceeds to execute tasks in the existing worktree

#### Scenario: New worktree created when none exists
- **WHEN** the user runs with `--worktree` and no worktree for that change name exists
- **THEN** the system creates a new worktree as before

### Requirement: Worktree creation uses git worktree add
The system SHALL use `git worktree add` to create worktrees with a branch named after the change.

#### Scenario: Correct branch naming
- **WHEN** a worktree is created for change `add-user-auth`
- **THEN** the git branch is named `add-user-auth` and the worktree directory is `<worktrees_dir>/add-user-auth`

#### Scenario: Worktree branches from HEAD
- **WHEN** a worktree is created
- **THEN** the new branch is based on the current HEAD of the main working tree

### Requirement: Worktree detection via git common dir
The system SHALL detect existing worktrees by inspecting the git common directory, not by checking the filesystem directly.

#### Scenario: Detect existing worktree
- **WHEN** the system checks whether a worktree exists for a change
- **THEN** it uses `git worktree list` or the git common directory to determine existence

#### Scenario: No worktrees directory does not error
- **WHEN** the worktrees directory does not yet exist on disk
- **THEN** the detection reports no worktrees (does not error)

## Boundaries
- ALWAYS: Check worktree-exists before IsClean -- reuse path must not be blocked by dirty tree check
- NEVER: Error when a worktree already exists for a change -- always reuse

## Gotchas
- `IsClean` check must run AFTER worktree-exists check, not before. Otherwise the verify-fix loop fails when `tasks.json` has uncommitted changes in an existing worktree.
  (learned: parallel-worktree-workflow, 2026-03-28)
