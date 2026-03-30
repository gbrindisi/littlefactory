# merge-command

## What It Does
Provides a CLI command that orchestrates a verify-fix loop, rebases the worktree branch onto main if needed, merges with `git merge --no-ff`, and cleans up the worktree and branch. This is the final step in the parallel worktree workflow, allowing completed changes to be merged back to main.

## Requirements

### Requirement: Merge command orchestrates verify-fix-merge cycle
The system SHALL provide a `merge` subcommand that runs a verify-fix loop, rebases if needed, merges the worktree branch into main, and cleans up.

#### Scenario: Merge with all tasks done and specs satisfied
- **WHEN** the user runs `littlefactory merge -c <name>` and verify passes
- **THEN** the system merges the branch into main and cleans up the worktree

#### Scenario: Merge with drift triggers fix loop
- **WHEN** the user runs `littlefactory merge -c <name>` and verify detects drift
- **THEN** the system runs `littlefactory run -c <name>` in the worktree to address remediation tasks, then re-verifies

### Requirement: Merge requires change flag
The system SHALL require the `--change` flag on the merge command.

#### Scenario: Merge without change flag errors
- **WHEN** the user runs `littlefactory merge` without `--change`
- **THEN** the system exits with an error indicating that `--change` is required

### Requirement: Merge requires existing worktree
The system SHALL require an existing worktree for the change. The merge command operates on worktree branches only.

#### Scenario: Merge with no worktree errors
- **WHEN** no worktree exists for the given change name
- **THEN** the system exits with an error indicating no worktree found

### Requirement: Merge validates task completion before verify
The system SHALL check that all tasks in `tasks.json` have status `done` before running verification. If incomplete tasks remain, merge SHALL exit with an error unless `--force` is provided.

#### Scenario: Incomplete tasks block merge
- **WHEN** `tasks.json` has tasks not in `done` status and `--force` is not set
- **THEN** the system exits with an error listing the incomplete tasks

#### Scenario: Force flag bypasses task check
- **WHEN** `tasks.json` has incomplete tasks and `--force` is set
- **THEN** the system proceeds to verification despite incomplete tasks

### Requirement: Merge rebases before merging
The system SHALL rebase the worktree branch onto main if main has advanced since the branch was created.

#### Scenario: Branch is up to date with main
- **WHEN** main has not advanced past the branch point
- **THEN** the system proceeds directly to merge without rebasing

#### Scenario: Branch needs rebase
- **WHEN** main has advanced past the branch point
- **THEN** the system runs `git rebase main` in the worktree before merging

#### Scenario: Rebase conflict aborts merge
- **WHEN** the rebase encounters conflicts
- **THEN** the system aborts the rebase (`git rebase --abort`), exits with an error, and tells the user to resolve conflicts manually

### Requirement: Merge uses no-fast-forward merge
The system SHALL merge the branch into main using `git merge --no-ff` to preserve branch history.

#### Scenario: No-fast-forward merge
- **WHEN** the merge proceeds
- **THEN** the system runs `git merge --no-ff <branch>` on the main branch, creating a merge commit

### Requirement: Merge cleans up worktree and branch
The system SHALL remove the worktree and delete the branch after a successful merge.

#### Scenario: Cleanup after merge
- **WHEN** the merge completes successfully
- **THEN** the system runs `git worktree remove <path>` and `git branch -d <branch>`

#### Scenario: Cleanup failure does not revert merge
- **WHEN** the merge succeeds but cleanup fails (e.g., worktree remove errors)
- **THEN** the merge stands; the system warns about cleanup failure but exits 0

### Requirement: Merge verify-fix loop has retry limit
The system SHALL limit the verify-fix loop to a configurable maximum number of retries (default 3).

#### Scenario: Max retries exhausted
- **WHEN** the verify-fix loop exceeds `--max-verify-retries` (default 3) without passing verification
- **THEN** the system exits with an error indicating verification failed after N attempts

#### Scenario: Custom retry limit
- **WHEN** the user runs `littlefactory merge -c <name> --max-verify-retries 5`
- **THEN** the system allows up to 5 verify-fix cycles

## Boundaries
- ALWAYS: Abort rebase on conflict rather than leaving worktree in broken state
- ALWAYS: Use `--no-ff` for merges to preserve branch history
- NEVER: Revert a successful merge due to cleanup failures

## Gotchas
- The verify-fix loop shells out to `littlefactory run -c <name> -w` which triggers the worktree reuse path. The `-w` flag is needed so the driver `chdir`s into the worktree.
  (learned: parallel-worktree-workflow, 2026-03-28)
- `git rebase --abort` is essential after conflict detection. Without it, the worktree is left in a broken rebase state.
  (learned: parallel-worktree-workflow, 2026-03-28)
