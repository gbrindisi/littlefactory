## Context

The `init` command currently runs 4 steps: Factoryfile creation, AGENTS.md setup, .gitignore updates, and skill installation. OpenSpec is a required tool for the littlefactory artifact-driven workflow, but init does not verify its presence or configure it. Users must manually set up `openspec/schemas/littlefactory/` and `openspec/config.yaml`.

The codebase already has a pattern for embedding and extracting files at init time (`internal/init/skills/embed.go`), which uses `//go:embed` with `fs.Sub` to walk and copy an embedded directory tree.

## Goals / Non-Goals

**Goals:**
- Fail early with a clear message if `openspec` is not in PATH
- Install the littlefactory schema into `openspec/schemas/littlefactory/` during init
- Create `openspec/config.yaml` with `schema: littlefactory` (only if missing)
- Apply the same setup during `upgrade` for existing projects
- Follow existing sub-package conventions (`internal/init/openspec/`)

**Non-Goals:**
- Installing or downloading the `openspec` binary itself
- Managing multiple schemas or letting users pick a schema at init time
- Validating the openspec version or feature compatibility
- Modifying existing `openspec/config.yaml` content (preserve user customizations)

## Decisions

### 1. New sub-package: `internal/init/openspec/`

**Decision**: Create a self-contained sub-package following the same pattern as `agentsmd/`, `gitignore/`, and `skills/`.

**Rationale**: Each init concern is isolated in its own package with its own types, functions, and tests. This keeps the main `init.go` orchestrator thin.

**Alternative considered**: Adding the logic directly to `init.go`. Rejected because it breaks the established pattern and makes testing harder.

### 2. Binary check via `exec.LookPath`

**Decision**: Use `exec.LookPath("openspec")` to verify the binary is available before any init work begins.

**Rationale**: `exec.LookPath` is the standard Go approach for checking binary availability. It respects PATH and is cross-platform. The check runs as the very first step in `Run()`, before Factoryfile creation, so we fail fast.

**Alternative considered**: Running `openspec --version` to verify the binary works. Rejected as unnecessary complexity -- we only need to know the binary exists, not its version.

### 3. Embed schema files in the binary

**Decision**: Embed the schema directory tree (`openspec/schemas/littlefactory/`) into the new sub-package using `//go:embed all:embedded/schema` and extract it to `openspec/schemas/littlefactory/` at init time.

**Rationale**: This mirrors exactly how skills are embedded and extracted (`skills/embed.go`). The schema files ship with the binary so there is no runtime dependency on the source tree.

**File layout inside the sub-package:**
```
internal/init/openspec/
  openspec.go          # CheckInstalled(), Setup() functions
  embed.go             # //go:embed, ExtractSchema()
  embedded/schema/     # Copy of openspec/schemas/littlefactory/ contents
    schema.yaml
    templates/
      proposal.md
      spec.md
      design.md
      tasks.md
      tasks.json
  openspec_test.go
```

**Alternative considered**: Reading schema files from the source tree at runtime. Rejected because the binary must be self-contained.

### 4. Config file handling

**Decision**: Create `openspec/config.yaml` with `schema: littlefactory` only if the file does not already exist. If it exists, leave it untouched.

**Rationale**: The config file can contain user customizations (`context:`, `rules:`). Overwriting would lose these. Since the schema is installed to `openspec/schemas/littlefactory/`, the user can manually set `schema: littlefactory` if they have a different default. For new projects, we create the minimal config. For existing projects, we respect their configuration.

**Alternative considered**: Parsing existing YAML and patching only the `schema` field. Rejected as over-engineered -- users who upgrade and want the littlefactory schema can edit one line.

### 5. Step ordering in init

**Decision**: The openspec binary check runs as a precondition before step 1. The openspec setup (schema extraction + config) runs as the new step 5 (after skill installation).

**Rationale**: The binary check must happen first to fail fast. The schema setup logically comes after the other scaffolding is in place. Init goes from 4 steps to 5, and upgrade goes from 3 steps to 4.

**Sequence:**
```
CheckInstalled()          # Precondition: fail if openspec not in PATH
[1/5] Creating Factoryfile
[2/5] Setting up AGENTS.md
[3/5] Updating .gitignore
[4/5] Installing skills
[5/5] Setting up OpenSpec    # New step: extract schema + write config if missing
```

## Risks / Trade-offs

- **[openspec not installed]** Users without openspec will be blocked from running init entirely. This is intentional -- openspec is a required dependency. The error message should include installation guidance.

- **[Schema drift]** The embedded schema is a snapshot at build time. If the schema evolves, the binary must be rebuilt. This is acceptable since the schema is part of the littlefactory distribution.

- **[Existing config not updated]** If a user has an existing `openspec/config.yaml` with a different default schema, we won't change it. This is intentional to preserve customizations, but means they must manually set `schema: littlefactory` if desired.
