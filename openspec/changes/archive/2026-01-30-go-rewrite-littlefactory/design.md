## Context

Ciccio is a ~750 LOC Python application with three main components: CLI (click), Ralph loop driver, and Rich-based dashboard. The core value is the autonomous agent loop pattern, but Python packaging creates friction for distribution. The dashboard component is unused and adds complexity.

The goal is to rewrite as "littlefactory" in Go, preserving the core loop behavior while establishing extensibility for a future lightweight orchestration platform. The rewrite must maintain backward compatibility with beads task system and metadata JSON format.

**Constraints:**
- Must preserve exact beads integration behavior (bd CLI shell commands)
- Must maintain JSON metadata format for backward compatibility
- Must preserve template injection mechanism and CLAUDE.md format
- Must keep "keep-going" error handling strategy
- Must maintain visual output style (banners, summaries)

**Stakeholders:**
- Single user migrating from Python ciccio to Go littlefactory

## Goals / Non-Goals

**Goals:**
- Single Go binary with no external dependencies (besides bd CLI)
- Interface-based design allowing future agent/task source implementations
- Full behavioral compatibility with ciccio Ralph loop
- Embedded template with local override support
- Factoryfile configuration with CLI flag overrides
- Sub-500 LOC implementation (vs 750 LOC Python with dashboard)

**Non-Goals:**
- Dashboard/TUI (explicitly removed - not needed)
- Distribution tooling (homebrew, releases) - out of scope
- Integration tests with real beads - unit tests with mocks only
- Agent implementations beyond Claude Code
- Task source implementations beyond beads
- Config keys beyond max_iterations and timeout

## Decisions

### 1. Directory Structure

**Decision:** Organize by domain concern under `internal/` with cmd/ for CLI entry.

```
littlefactory/
├── cmd/littlefactory/main.go          # Cobra CLI entry
├── internal/
│   ├── agent/                         # Agent interface + Claude impl
│   ├── tasks/                         # TaskSource interface + Beads impl
│   ├── driver/                        # Loop orchestrator + metadata
│   ├── template/                      # Template loading/rendering
│   └── config/                        # Factoryfile loading
├── templates/CLAUDE.md                # Embedded via go:embed
└── go.mod
```

**Rationale:** Clean separation of concerns, `internal/` prevents external imports, domain-based packages (agent, tasks, driver) rather than technical layers (models, services).

**Alternative considered:** Flat structure under cmd/ - rejected due to lack of organization for ~500 LOC.

### 2. Agent Interface Design

**Decision:** Interface with single `Run(ctx, prompt, sessionID)` method returning result struct.

```go
type Agent interface {
    Run(ctx context.Context, prompt string, sessionID string) (AgentResult, error)
}

type AgentResult struct {
    ExitCode    int
    Output      string
    OutputLines int
    OutputBytes int
}
```

**Rationale:**
- Context-aware for timeout enforcement
- Session ID passed explicitly for tracking (Claude Code specific but extensible)
- Result struct captures all observability data
- Error return for execution failures distinct from non-zero exit codes

**Alternative considered:** Streaming interface with channel-based output - rejected as overkill for current needs, can add later if needed.

### 3. TaskSource Interface Design

**Decision:** Interface mirroring beads operations: Ready(), Show(id), Close(id, reason), Sync().

```go
type TaskSource interface {
    Ready() ([]Task, error)
    Show(id string) (*Task, error)
    Close(id, reason string) error
    Sync() error
}
```

**Rationale:**
- Maps 1:1 to beads CLI operations (bd ready, bd show, bd close, bd sync)
- Simple enough to implement for other task systems (GitHub issues, todo.txt)
- Sync() explicit rather than automatic for control

**Alternative considered:** Active Record pattern with task.Close() methods - rejected as mixing concerns, harder to mock.

### 4. Template System Approach

**Decision:** Simple string replacement (ReplaceAll) rather than text/template.

```go
func Render(tmpl string, task *Task) string {
    if task == nil {
        return tmpl
    }
    result := strings.ReplaceAll(tmpl, "{task_id}", task.ID)
    result = strings.ReplaceAll(result, "{task_title}", task.Title)
    result = strings.ReplaceAll(result, "{task_description}", task.Description)
    return result
}
```

**Rationale:**
- Matches Python behavior exactly (simple replace)
- No template parsing overhead
- Only 3 injection points, no complex logic needed
- Can upgrade to text/template later if complexity grows

**Alternative considered:** text/template with {{.TaskID}} syntax - rejected as unnecessary complexity for current needs.

### 5. Configuration Management

**Decision:** gopkg.in/yaml.v3 for Factoryfile parsing with struct-based config.

```go
type Config struct {
    MaxIterations int           `yaml:"max_iterations"`
    Timeout       int           `yaml:"timeout"`
    Agent         *AgentConfig  `yaml:"agent,omitempty"`  // Future use
}
```

**Loading order:**
1. Hardcoded defaults in code
2. Factoryfile (if exists) - unmarshal into struct
3. CLI flags - override struct fields

**Rationale:**
- YAML widely understood, good Go library support
- Struct tags provide clear schema
- Optional fields (omitempty) support future expansion
- Simple precedence: last one wins

**Alternative considered:** TOML or JSON - rejected as YAML more human-friendly for config files.

### 6. Error Handling Strategy

**Decision:** Preserve "keep-going" strategy - failed iterations do not abort run.

```go
func (d *Driver) runIteration(ctx context.Context, num int) {
    // ... setup ...

    result, err := d.agent.Run(ctx, prompt, sessionID)

    switch {
    case ctx.Err() == context.DeadlineExceeded:
        iter.Status = StatusTimeout
    case err != nil:
        iter.Status = StatusFailed
        iter.ErrorMessage = ptr(err.Error())
    case result.ExitCode != 0:
        iter.Status = StatusFailed
    default:
        iter.Status = StatusCompleted
    }

    // Always continue to next iteration
}
```

**Rationale:**
- Matches Python behavior exactly
- Allows agent to fail/timeout without stopping entire run
- User can inspect metadata to see which iterations failed

**Alternative considered:** Configurable abort-on-failure - rejected as over-engineering, can add later if needed.

### 7. Metadata JSON Compatibility

**Decision:** Match Python JSON structure exactly using struct tags and ISO8601 timestamps.

```go
type RunMetadata struct {
    RunID                      string               `json:"run_id"`
    StartedAt                  time.Time            `json:"started_at"`
    EndedAt                    *time.Time           `json:"ended_at,omitempty"`
    // ... other fields with exact JSON keys
}

// Marshal with custom time format
func (r *RunMetadata) MarshalJSON() ([]byte, error) {
    // Convert time.Time to ISO8601 strings matching Python's isoformat()
}
```

**Rationale:**
- Backward compatibility with existing tooling/scripts
- Pointer fields for optional values (matches Python None)
- Custom marshal for timestamp format control

**Alternative considered:** New JSON format with better naming - rejected due to backward compatibility requirement.

### 8. Session Path Computation

**Decision:** Replicate Python's exact encoding logic for Claude Code session paths.

```go
func sessionPath(projectRoot, sessionID string) string {
    encoded := strings.ReplaceAll(projectRoot, "/", "-")
    encoded = strings.ReplaceAll(encoded, ".", "-")
    homeDir, _ := os.UserHomeDir()
    return filepath.Join(homeDir, ".claude", "projects", encoded, sessionID+".jsonl")
}
```

**Rationale:**
- Claude Code specific but needed for session traceability
- Python ciccio computes this, must match exactly
- Future agents may not need this, but interface allows ignoring

**Alternative considered:** Abstract session storage - rejected as premature, this is Claude-specific.

### 9. Logging Approach

**Decision:** Use fmt.Printf() for now, no structured logging.

**Rationale:**
- Simple, matches Python's print() behavior
- Visual banners require formatted output anyway
- Can add log/slog later if needed for structured output

**Alternative considered:** log/slog from the start - rejected as overkill for current needs.

### 10. CLI Framework

**Decision:** Use github.com/spf13/cobra for CLI.

**Rationale:**
- Industry standard Go CLI framework
- Good flag handling and subcommand support
- Easy to add commands later if needed
- Better than stdlib flag for complex CLIs

**Alternative considered:** urfave/cli/v2 - both are good, cobra chosen for consistency with ecosystem.

## Risks / Trade-offs

### Risk: Breaking changes to beads CLI output format
**Impact:** JSON parsing failures if bd ready/show/list change format
**Mitigation:** Error handling around JSON unmarshaling, integration tests (future), version pinning for bd

### Risk: Session path computation diverges from Claude Code
**Impact:** Incorrect session paths in metadata if Claude Code changes storage location
**Mitigation:** Document assumption, add comment with Claude Code version tested against

### Risk: Template embedding increases binary size
**Impact:** Minimal - CLAUDE.md is ~3KB, negligible for binary size
**Mitigation:** None needed, acceptable trade-off for single binary goal

### Trade-off: Simple string replacement limits template flexibility
**Impact:** Cannot add conditionals, loops, or complex logic in templates
**Mitigation:** Acceptable for current needs, can upgrade to text/template if requirements change

### Trade-off: No streaming output during agent execution
**Impact:** User sees no output until iteration completes (vs Python showing real-time)
**Mitigation:** Acceptable, output is captured and printed after completion, matches Python --print behavior

### Trade-off: Beads shell integration vs native library
**Impact:** Slower than native Go library, requires bd CLI installed
**Mitigation:** Acceptable, beads is Go-based so shell overhead minimal, requirement is bd CLI integration

## Migration Plan

### Deployment
1. Build Go binary: `go build -o littlefactory ./cmd/littlefactory`
2. Install to PATH: `cp littlefactory ~/bin/` or similar
3. Test against existing ciccio project with beads tasks
4. Verify metadata JSON format matches
5. Compare progress.txt output format
6. Replace `ciccio start` with `littlefactory start` in workflows

### Rollback
- Keep Python ciccio installed during transition
- If issues arise, `ciccio start` still works
- No database migrations or state changes, just binary swap

### Testing Strategy
1. Unit tests for each package (agent, tasks, driver, template, config)
2. Mock implementations for testing (MockAgent, MockTaskSource)
3. Golden file tests for metadata JSON format
4. Manual integration test with real beads project

## Open Questions

None - design is complete for v1 scope.
