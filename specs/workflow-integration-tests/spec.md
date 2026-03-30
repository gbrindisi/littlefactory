# workflow-integration-tests

## What It Does
Provides binary-level integration tests that exercise the full worktree workflow (run, verify, merge) as subprocesses. Tests build the littlefactory binary, create real git repos with mock changes, and assert on exit codes, filesystem state, and stdout output. Uses `echo done` as a deterministic agent.

## Requirements

### Requirement: Binary build helper for integration tests
The test suite SHALL build the littlefactory binary once per test run using `TestMain` and make it available to all integration tests.

#### Scenario: Binary built before tests
- **WHEN** integration tests run
- **THEN** the binary is built to a temp directory and its path is available to all test functions

#### Scenario: Binary cleaned up after tests
- **WHEN** all integration tests complete
- **THEN** the temp directory containing the binary is removed

### Requirement: Repo scaffolding helper
The test suite SHALL provide a helper that creates a git repo with a Factoryfile using `echo done` as the agent and a change directory with tasks.json.

#### Scenario: Scaffold creates usable repo
- **WHEN** the scaffold helper is called with a change name and task list
- **THEN** a git repo exists with initial commit, Factoryfile, and `.littlefactory/changes/<name>/tasks.json`

### Requirement: Run workflow integration test
The test suite SHALL include a test that exercises `littlefactory run -c <name> -w` as a subprocess and verifies worktree creation, task completion in the worktree, and worktree reuse.

#### Scenario: Run creates worktree and completes tasks
- **WHEN** `littlefactory run -c test-change -w` executes against a repo with 2 todo tasks
- **THEN** exit code is 0, a worktree exists, and tasks.json in the worktree shows all tasks done

#### Scenario: Run reuses existing worktree
- **WHEN** `littlefactory run -c test-change -w` executes a second time
- **THEN** exit code is 0, stdout contains "Reusing existing worktree", and no new worktree is created

### Requirement: Status workflow integration test
The test suite SHALL include a test that exercises `littlefactory status --all` after a completed run and verifies the output includes run state and merge readiness.

#### Scenario: Status shows ready to merge
- **WHEN** `littlefactory status --all` runs after a completed worktree run
- **THEN** stdout contains the change name, task counts, and `[ready to merge]`

### Requirement: Verify workflow integration test
The test suite SHALL include a test that exercises `littlefactory verify -c <name>` as a subprocess and verifies it runs in the worktree context.

#### Scenario: Verify passes with echo agent
- **WHEN** `littlefactory verify -c test-change` executes and the echo agent exits 0
- **THEN** exit code is 0 and stdout contains "Running verification in worktree"

### Requirement: Merge workflow integration test
The test suite SHALL include a test that exercises the full merge lifecycle: verify pass, merge to main, cleanup.

#### Scenario: Merge completes full cycle
- **WHEN** `littlefactory merge -c test-change` executes after a successful run
- **THEN** exit code is 0, the worktree is removed, the branch is deleted, and main contains the worktree's commits

#### Scenario: Merge with advanced main rebases
- **WHEN** main has advanced with non-conflicting commits since the worktree branched
- **THEN** merge rebases successfully, merges to main, and both sets of changes exist on main

### Requirement: Integration tests behind build tag
Integration tests SHALL use a `//go:build integration` build tag so they don't run with `go test ./...` (they require a binary build step).

#### Scenario: Regular test run skips integration tests
- **WHEN** `go test ./...` runs without build tags
- **THEN** integration tests are not executed

#### Scenario: Integration tests run with tag
- **WHEN** `go test -tags=integration ./cmd/littlefactory/` runs
- **THEN** integration tests execute

## Boundaries
- ALWAYS: Use `//go:build integration` tag on integration test files to avoid slowing down `go test ./...`
- ALWAYS: Each integration test must be self-contained (own temp dir, own git repo)
- NEVER: Use a real LLM agent in integration tests -- always use `echo done` or similar deterministic command

## Gotchas
- `TestMain` with `//go:build integration` means there are TWO TestMain functions in the package -- the integration one and potentially any other. Go handles this via build tags, but adding a second TestMain without a build tag will break compilation.
  (learned: add-workflow-integration-tests, 2026-03-31)
- The `echo done` agent runs via PTY in the real binary. On some CI environments PTY allocation may fail. Tests may need `TERM=dumb` or similar.
  (learned: add-workflow-integration-tests, 2026-03-31)
