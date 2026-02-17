## 1. Package Structure and Embedded Files

- [x] 1.1 Create `internal/init/openspec/` directory
- [x] 1.2 Copy `openspec/schemas/littlefactory/` contents into `internal/init/openspec/embedded/schema/`
- [x] 1.3 Create `internal/init/openspec/embed.go` with `//go:embed all:embedded/schema` and `ExtractSchema(projectRoot string) error` function using `fs.Sub` and `fs.WalkDir` (mirror `skills/embed.go` pattern)

## 2. Core OpenSpec Setup Logic

- [x] 2.1 Create `internal/init/openspec/openspec.go` with `CheckInstalled() error` using `exec.LookPath("openspec")`
- [x] 2.2 Add `Setup(projectRoot string) error` function that calls `ExtractSchema` and conditionally writes `openspec/config.yaml`
- [x] 2.3 Config logic: create `openspec/config.yaml` with `schema: littlefactory` only if file does not exist

## 3. Integrate into Init and Upgrade

- [x] 3.1 Update `internal/init/init.go` `Run()` to call `openspec.CheckInstalled()` before any steps, returning error if not found
- [x] 3.2 Add new step 5 in `Run()` calling `openspec.Setup(projectRoot)` with logger output
- [x] 3.3 Update logger total from 4 to 5 in `Run()`
- [x] 3.4 Update `internal/init/upgrade.go` `Upgrade()` to call `openspec.CheckInstalled()` before any steps
- [x] 3.5 Add OpenSpec setup step to `Upgrade()` and update logger total from 3 to 4

## 4. Tests

- [x] 4.1 Write tests for `CheckInstalled()` (binary found and not-found cases)
- [x] 4.2 Write tests for `ExtractSchema()` verifying schema.yaml and template files are extracted to correct paths
- [x] 4.3 Write tests for `Setup()` verifying config creation when missing and config preservation when existing
- [x] 4.4 Update existing init integration tests to account for the openspec prerequisite check and new step count
- [x] 4.5 Update existing upgrade integration tests to account for the openspec prerequisite check and new step count
