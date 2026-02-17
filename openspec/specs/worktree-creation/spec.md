## ADDED Requirements

### Requirement: Worktree flag on run command
The system SHALL accept `--worktree` / `-w` flag on the `run` command to create a new git worktree for the change.

#### Scenario: Run with worktree flag creates worktree
- **WHEN** User runs `littlefactory run claude -c feature-a -w`
- **THEN** System creates a git worktree for branch `feature-a` and runs the agent in that worktree

#### Scenario: Worktree flag requires change flag
- **WHEN** User runs `littlefactory run claude -w` without `--change`
- **THEN** System returns error "The --worktree flag requires --change to specify the branch name"

### Requirement: Worktree creation requires clean working tree
The system SHALL refuse to create a worktree if there are uncommitted changes.

#### Scenario: Dirty working tree blocks worktree creation
- **WHEN** User runs `littlefactory run claude -c feature-a -w` with uncommitted changes
- **THEN** System returns error "Uncommitted changes detected. Commit or stash before creating worktree."

#### Scenario: Clean working tree allows worktree creation
- **WHEN** User runs `littlefactory run claude -c feature-a -w` with clean working tree
- **THEN** System proceeds with worktree creation

### Requirement: Worktree flag errors if worktree exists
The system SHALL refuse to create a worktree if one already exists for the branch.

#### Scenario: Existing worktree blocks creation
- **WHEN** User runs `littlefactory run claude -c feature-a -w` and worktree for `feature-a` exists
- **THEN** System returns error "Worktree for 'feature-a' already exists at <path>. Run without -w to use existing worktree."

### Requirement: Worktree creation uses git worktree add
The system SHALL create worktrees using `git worktree add` command.

#### Scenario: Worktree created with correct branch
- **WHEN** System creates worktree for change `feature-a`
- **THEN** System executes `git worktree add <worktrees_dir>/feature-a -b feature-a`

#### Scenario: Worktree branches from current HEAD
- **WHEN** System creates worktree for change `feature-a`
- **THEN** System creates branch from current HEAD

### Requirement: Worktree detection via git common dir
The system SHALL detect worktree-enabled repos using `git rev-parse --git-common-dir`.

#### Scenario: Detect worktrees exist
- **WHEN** System checks for worktree support
- **THEN** System runs `git rev-parse --git-common-dir` and checks if `<common-dir>/worktrees` exists and is non-empty

#### Scenario: No worktrees directory
- **WHEN** `<common-dir>/worktrees` does not exist
- **THEN** System reports no existing worktrees (but can still create one)
