## Context

Agents in littlefactory are configured via Factoryfile with a `command` string. When agents need environment variables (e.g., API keys), users must currently wrap commands in shell invocations with complex quoting. This is fragile and hard to maintain.

Current workaround:
```yaml
ab:
  command: "sh -c \"ANTHROPIC_API_KEY=$(security find-generic-password -w -s 'Claude Code') agent-box run --\""
```

## Goals / Non-Goals

**Goals:**
- Allow declarative env var configuration per agent
- Support static string values
- Support dynamic values evaluated via shell command
- Maintain backward compatibility (env is optional)

**Non-Goals:**
- Inheriting specific vars from parent (just inherit all by default)
- Shell expansion for the command itself (separate feature)
- Encrypted/secret management beyond shell commands

## Decisions

### Decision 1: EnvValue as union type

Env values support two forms in YAML:
```yaml
env:
  STATIC_VAR: "literal"        # string -> static
  DYNAMIC_VAR:
    shell: "command here"      # object with shell -> dynamic
```

**Rationale**: Clean syntax for common case (static), explicit for dynamic. Alternative considered: always use object form (`value:` + optional `shell: true`) - rejected as more verbose for static values.

### Decision 2: Shell evaluation at agent start time

Dynamic env values are evaluated once when the agent starts, not on each iteration.

**Rationale**: Credentials typically don't change mid-run. Evaluating once avoids repeated subprocess overhead. If per-iteration refresh is needed, users can wrap the whole command.

### Decision 3: Inherit parent environment

`cmd.Env` is set to `os.Environ()` plus overrides from config.

**Rationale**: Standard Unix behavior. Agents often need PATH, HOME, etc. Explicit opt-out would be more complex.

### Decision 4: Shell errors are fatal

If a shell command in `env` fails (non-zero exit), the agent fails to start.

**Rationale**: Missing credentials usually mean the agent can't function. Fail fast rather than run with incomplete env.

## Risks / Trade-offs

- **Shell command stdout trimming**: Must trim trailing newlines from shell output (common for command substitution). Forgetting this would add `\n` to values.
- **No secret masking**: Env values may appear in logs/errors. Users should be aware. Mitigation: document this limitation.
- **YAML quoting complexity**: Users must understand YAML string quoting for shell commands with special chars. Mitigation: provide examples in docs.
