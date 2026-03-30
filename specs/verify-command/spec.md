# verify-command

## What It Does
Provides a CLI command that invokes the configured agent with a verification prompt to check implementation against change specs. The agent runs non-interactively, exits 0 on pass or non-zero on drift (with remediation tasks appended to tasks.json). This makes verification composable by other commands like `merge`.

## Requirements

### Requirement: Verify command runs agent-driven spec validation
The system SHALL provide a `verify` subcommand that invokes the configured agent with a verification prompt to check implementation against change specs.

#### Scenario: Verify invokes agent with verifier template
- **WHEN** the user runs `littlefactory verify -c <name>`
- **THEN** the system loads the VERIFIER.md template, renders it with change context (proposal, specs, design), and executes the configured agent with that prompt

#### Scenario: Verify exits 0 on pass
- **WHEN** the verification agent exits with code 0
- **THEN** the verify command exits with code 0 indicating all specs are satisfied

#### Scenario: Verify exits 1 on drift
- **WHEN** the verification agent exits with non-zero code
- **THEN** the verify command exits with code 1 indicating drift was found

### Requirement: Verify requires change flag
The system SHALL require the `--change` flag on the verify command.

#### Scenario: Verify without change flag errors
- **WHEN** the user runs `littlefactory verify` without `--change`
- **THEN** the system exits with an error indicating that `--change` is required

### Requirement: Verify resolves worktree context
The system SHALL run verification in the worktree directory if a worktree exists for the change, otherwise in the project root.

#### Scenario: Verify in worktree
- **WHEN** a worktree exists for the change name
- **THEN** the verify agent executes in the worktree directory

#### Scenario: Verify in project root
- **WHEN** no worktree exists for the change name
- **THEN** the verify agent executes in the project root directory

### Requirement: Verify agent appends remediation tasks
The VERIFIER.md template SHALL instruct the agent to append remediation tasks to `tasks.json` when drift is detected. Remediation tasks MUST follow the existing blocker chain convention (each new task blocked by the previous last task).

#### Scenario: Remediation tasks appended on drift
- **WHEN** the verification agent detects drift
- **THEN** it appends new tasks to `tasks.json` with status `todo`, wired into the blocker chain after the last existing task

#### Scenario: No tasks modified on pass
- **WHEN** the verification agent finds no drift
- **THEN** `tasks.json` is not modified

### Requirement: Verify uses configured agent
The system SHALL use the same agent resolution as the `run` command (default agent or positional arg).

#### Scenario: Verify with default agent
- **WHEN** no agent is specified
- **THEN** the system uses the `default_agent` from config

#### Scenario: Verify with explicit agent
- **WHEN** the user runs `littlefactory verify -c <name> claude`
- **THEN** the system uses the specified agent

## Boundaries
- ALWAYS: Map any non-zero agent exit to verify exit 1 (drift)

## Gotchas
- The verify agent writes remediation tasks directly to `tasks.json`. If the agent produces malformed JSON, the next `littlefactory run` will fail fast on task validation -- this is by design, not a bug.
  (learned: parallel-worktree-workflow, 2026-03-28)
