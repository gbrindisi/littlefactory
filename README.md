# littlefactory

An autonomous coding agent orchestrator that runs iterative loops to complete software engineering tasks using Claude Code (or other agents).

## Overview

littlefactory coordinates autonomous agent execution in a sequential loop:

1. Retrieves the next ready task from a JSON task file
2. Renders a prompt template with task details
3. Executes the configured agent (Claude Code, etc.) with the prompt
4. Tracks iteration metadata (timing, output, status)
5. Repeats until all tasks complete or max iterations reached

It supports multiple agent backends, git worktree isolation, and an embedded spec-driven workflow.

## Requirements

- **Go 1.24+** for building from source
- **[Claude Code](https://claude.ai/claude-code)** (or another compatible agent) installed and authenticated

## Installation

### From source

```bash
go install github.com/gbrindisi/littlefactory/cmd/littlefactory@latest
```

### Build locally

```bash
git clone https://github.com/gbrindisi/littlefactory.git
cd littlefactory
make build
make install
```

## Quick Start

1. Initialize a project:
   ```bash
   littlefactory init
   ```
   This creates a `Factoryfile`, sets up `AGENTS.md`, updates `.gitignore`, installs skills, and creates the `.littlefactory/changes/` directory.

2. The only requirement for littlefactory is a `tasks.json` with explicitly sequential tasks:
   ```json
   {
     "tasks": [
       {
         "id": "feat-a1b",
         "title": "Implement feature X",
         "description": "Details...",
         "status": "todo",
         "blockers": []
       },
       {
         "id": "feat-c2d",
         "title": "Add tests for feature X",
         "description": "Details...",
         "status": "todo",
         "blockers": ["feat-a1b"]
       }
     ]
   }
   ```
   Tasks form a linear chain via `blockers` -- each task (except the first) must reference exactly one predecessor. littlefactory processes them one at a time in order.

3. Run the agent loop:
   ```bash
   littlefactory run
   ```

### Spec-Driven Workflow

littlefactory embeds a lightweight spec-driven workflow you can use via Claude Code skills. Instead of writing `tasks.json` by hand, you can go from idea to implementation in five steps:

1. **`/lf:explore`** -- Think through a problem with a rubber-duck thinking partner. Explore the codebase, draw diagrams, compare options. No code is written -- this is pure thinking time.

2. **`/lf:formalize`** -- Turn the conversation into a structured change. Takes no arguments -- it derives a change name and generates all artifacts automatically:
   - `proposal.md` -- why this change is needed
   - `specs/*/spec.md` -- what the system should do (requirements with scenarios)
   - `design.md` -- how to implement it (key decisions and tradeoffs)
   - `tasks.json` -- sequential tasks with rich context for autonomous execution

   All artifacts are written to `.littlefactory/changes/<name>/`.

3. **`/lf:do`** -- Kick off `littlefactory run -c <name>` in the background and monitor progress. The agent works through tasks autonomously while you get status updates.

4. **`/lf:verify`** -- Verify the implementation against the change artifacts. Checks three dimensions: completeness (all tasks done, all requirements covered), correctness (implementation matches specs), and coherence (code follows design decisions).

5. **`/lf:archive`** -- Merge delta specs from the change into main specs, capture gotchas and boundaries as reusable knowledge, and optionally clean up the change directory.

## Configuration

littlefactory uses a `Factoryfile` (YAML) for configuration. Place it in your project root.

### Factoryfile Format

```yaml
max_iterations: 10
timeout: 600
default_agent: claude
worktrees_dir: ".."
state_dir: ".littlefactory"

agents:
  claude:
    command: "claude --dangerously-skip-permissions --print"
  custom:
    command: "my-agent run --"
    env:
      API_KEY:
        shell: "security find-generic-password -w -s MyKey"
```

### Fields

| Field | Default | Description |
|-------|---------|-------------|
| `max_iterations` | 10 | Maximum iterations before stopping |
| `timeout` | 600 | Timeout in seconds per iteration |
| `default_agent` | (required) | Name of the agent to use from the `agents` map |
| `agents` | (required) | Map of named agent configurations |
| `state_dir` | `.littlefactory` | Directory for state files (progress, metadata, tasks) |
| `worktrees_dir` | `..` | Base directory for git worktrees |

### Agent Configuration

Each agent entry has:

| Field | Description |
|-------|-------------|
| `command` | Shell command to invoke the agent |
| `env` | Optional environment variables (static strings or `{shell: "cmd"}` for dynamic values) |

### Configuration Precedence

1. **Defaults**: Built-in defaults
2. **Factoryfile**: `./Factoryfile` or `./Factoryfile.yaml`
3. **CLI flags**: `--max-iterations`, `--timeout`

## CLI Reference

### `littlefactory init`

Initialize a new project in the current directory. Creates a Factoryfile, sets up AGENTS.md, updates .gitignore, and installs skills. Fails if a Factoryfile already exists.

```bash
littlefactory init
```

### `littlefactory run`

Start the autonomous agent loop.

```bash
littlefactory run [agent] [flags]
```

If `agent` is not specified, uses `default_agent` from the Factoryfile.

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--max-iterations` | from config or 10 | Maximum iterations before stopping |
| `--timeout` | from config or 600 | Timeout in seconds per iteration |
| `-c, --change` | | Change name to use as task source |
| `-t, --tasks` | | Explicit path to a tasks.json file |
| `-w, --worktree` | false | Create a new git worktree for the change |

**Examples:**

```bash
# Run with defaults
littlefactory run

# Use a specific agent
littlefactory run custom

# Run tasks from a change in a worktree
littlefactory run -c my-feature -w

# Point to a specific tasks file
littlefactory run -t path/to/tasks.json

# Override iteration limits
littlefactory run --max-iterations 20 --timeout 300
```

### `littlefactory upgrade`

Upgrade an existing project. Applies AGENTS.md setup, .gitignore updates, and skill installation. All operations are idempotent.

```bash
littlefactory upgrade
```

### `littlefactory version`

Show version information.

```bash
littlefactory version
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success -- all tasks completed or max iterations reached |
| 1 | Error -- execution failed |
| 130 | Interrupted -- SIGINT received (Ctrl+C) |

## Task Management

Tasks are stored in a JSON file (default: `.littlefactory/tasks.json`):

```json
{
  "tasks": [
    {
      "id": "task-1",
      "title": "Task Title",
      "description": "Task description",
      "status": "todo",
      "labels": [],
      "blockers": []
    }
  ]
}
```

Status values: `todo`, `in_progress`, `done`.

Tasks can also come from a change directory via the `--change` flag (resolves to `.littlefactory/changes/<name>/tasks.json`), or from an explicit path via `--tasks`.

## Worktree Support

When running with `-c <change> -w`, littlefactory creates a git worktree for the change, isolating work on a dedicated branch. The worktree directory is determined by `worktrees_dir` in the Factoryfile.

## Template System

littlefactory uses templates to generate prompts for the agent. Templates support placeholder injection:

| Placeholder | Description |
|-------------|-------------|
| `{task_id}` | The task ID |
| `{task_title}` | The task title |
| `{task_description}` | The full task description |

Templates are loaded in order:

1. **Local override**: `<state_dir>/agents/WORKER.md` (if exists, e.g. `.littlefactory/agents/WORKER.md`)
2. **Embedded default**: Built into the binary

## Output Files

State files are stored in the state directory (default `.littlefactory/`):

| File | Description |
|------|-------------|
| `progress.md` | Markdown log of iteration sessions |
| `run_metadata.json` | JSON metadata for the current run |
| `tasks.json` | JSON task list |

## Project Detection

littlefactory detects the project root by searching for a `Factoryfile`, starting from the current working directory and walking up the directory tree.

## License

MIT
