## Context

Littlefactory currently writes runtime state to two locations:
- `.littlefactory/tasks.json` - task definitions and status
- `tasks/progress.txt` and `tasks/run_metadata.json` - run logs and metadata

The `tasks/` directory is legacy from the Python "ciccio" version. The Go rewrite moved tasks to `.littlefactory/` but left progress and metadata in the old location. This creates confusion about where state lives and makes configuration harder.

## Goals / Non-Goals

**Goals:**
- Consolidate all runtime state files into a single configurable directory
- Make state directory configurable via Factoryfile `state_dir` option
- Update progress format to proper markdown with "Little Factory" branding
- Clean up function signatures to receive full config for future extensibility

**Non-Goals:**
- Backward compatibility with old `tasks/` location (clean break)
- Automatic migration of existing files
- Changes to the task file format itself

## Decisions

### Decision 1: Add `state_dir` to Config struct

Add new field to `Config` struct with default value `.littlefactory`:

```go
const DefaultStateDir = ".littlefactory"

type Config struct {
    MaxIterations int    `yaml:"max_iterations"`
    Timeout       int    `yaml:"timeout"`
    StateDir      string `yaml:"state_dir"`  // NEW
    DefaultAgent  string `yaml:"default_agent"`
    Agents        map[string]AgentConfig `yaml:"agents"`
}
```

**Rationale**: Follows existing config pattern. YAML tag allows Factoryfile override.

**Alternative considered**: Environment variable - rejected because Factoryfile is the established config mechanism.

### Decision 2: Pass full `*config.Config` to state functions

Change function signatures from:
```go
func SaveMetadata(projectRoot string, metadata *RunMetadata) error
func InitProgressFile(projectRoot string) error
```

To:
```go
func SaveMetadata(projectRoot string, cfg *config.Config, metadata *RunMetadata) error
func InitProgressFile(projectRoot string, cfg *config.Config) error
```

**Rationale**: Future-proofs the API. If these functions ever need other config values (timeouts, formats, etc.), they have access without signature changes.

**Alternative considered**: Pass only `stateDir string` - rejected per user preference for full config access.

### Decision 3: Progress file format

New `progress.md` format with proper markdown:

```markdown
# Little Factory Progress Log

**Started:** 2026-02-03T14:37:00Z

---

## Iteration 1

- **Task:** task-123
- **Status:** completed

---
```

**Rationale**: Renders nicely in editors/GitHub, maintains append-only semantics, clean branding update.

### Decision 4: JSONTaskSource receives config at construction

```go
func NewJSONTaskSource(projectRoot string, cfg *config.Config) *JSONTaskSource {
    tasksPath := filepath.Join(projectRoot, cfg.StateDir, "tasks.json")
    // ...
}
```

**Rationale**: Consistent with other components. Stores config reference for potential future use.

## Risks / Trade-offs

**[Breaking change]** Existing `tasks/` files will be orphaned.
- Mitigation: Document in release notes. Users can manually delete `tasks/` directory.

**[Config validation]** Empty `state_dir` would cause issues.
- Mitigation: Validate `state_dir` is non-empty in config validation, default to `.littlefactory` if not specified.
