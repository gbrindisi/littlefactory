## Context

littlefactory currently hardcodes the `claude` executable and requires manual Factoryfile creation. The agent interface includes session ID management that adds complexity without clear benefit for the current use case. Users want to:
1. Bootstrap new projects quickly with `littlefactory init`
2. Configure different agents (claude, aider, custom wrappers) without code changes
3. Run specific agents via `littlefactory run <agent-name>`

Current config structure:
```yaml
max_iterations: 10
timeout: 600
agent:
  type: "..."  # unused placeholder
```

## Goals / Non-Goals

**Goals:**
- Add `init` command to scaffold Factoryfile with sensible defaults
- Support multiple named agents via `agents:` map in config
- Make agent command fully configurable (binary + arguments)
- Simplify agent interface by removing session ID management

**Non-Goals:**
- Template variables in command (e.g., `{{session_id}}`) - deferred for future
- Per-agent timeout/iteration settings - keep global for simplicity
- Interactive init wizard - just dump defaults
- Migration tooling for existing Factoryfiles - document breaking change

## Decisions

### Decision: Command as full string vs structured fields
**Choice**: Single `command` string field containing executable and all arguments.

**Alternatives considered**:
- Structured: `{ executable: "claude", args: ["--print", "--flag"] }` - More validation possible but verbose
- Template: `command: "claude {{session_id}}"` - Flexible but adds parsing complexity

**Rationale**: A single command string is simple, familiar (like shell), and sufficient. Users can include any arguments they need. If structured validation becomes important later, we can add it.

### Decision: Agent selection via positional argument
**Choice**: `littlefactory run [agent-name]` where agent-name is optional positional arg.

**Alternatives considered**:
- Flag: `littlefactory run --agent claude` - More explicit but verbose for common case
- Subcommand per agent: `littlefactory run-claude` - Doesn't scale

**Rationale**: Positional is concise. Default agent fallback handles the common case of single-agent setups.

### Decision: Remove session ID from interface
**Choice**: Remove `sessionID` parameter from `Agent.Run()` interface entirely.

**Alternatives considered**:
- Keep optional: Makes interface messier for agents that don't use it
- Move to config: Could add `session_flag: "--session-id"` but adds complexity

**Rationale**: Session ID was for debugging (finding conversation jsonl). This can be re-added later with template variables if needed. Simplifying now reduces scope.

### Decision: Fail on existing Factoryfile
**Choice**: `init` command fails if Factoryfile already exists (no overwrite).

**Alternatives considered**:
- Overwrite with `--force` flag - Risk of accidental data loss
- Merge with existing - Complex to implement correctly

**Rationale**: Safe default. Users can delete and re-init if needed.

## Risks / Trade-offs

- **Breaking change to config format** - Existing Factoryfiles will fail to parse. Mitigation: Document migration in changelog, provide example of new format.
- **No session continuity** - Removing session ID means each iteration may lose context. Mitigation: Most agents maintain their own context; can re-add later if needed.
- **Command injection risk** - User controls full command string. Mitigation: This is intentional - users configure their own agents. Document that command runs with shell semantics.
