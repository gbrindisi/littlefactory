## Why

The `littlefactory init` command sets up project scaffolding but does not configure OpenSpec, which is required for the artifact-driven workflow. Users must manually copy the littlefactory schema and create the OpenSpec config after running init. This change makes init handle OpenSpec setup automatically, failing early if the `openspec` binary is not installed.

## What Changes

- `init` checks that the `openspec` binary is available in PATH before proceeding; aborts with a clear error if not found
- `init` copies the embedded littlefactory schema to `openspec/schemas/littlefactory/`
- `init` creates `openspec/config.yaml` with `schema: littlefactory` (only if config doesn't exist, to preserve user customizations)
- Step count increases from 4 to 5 for init, and 3 to 4 for upgrade
- `upgrade` command also gains the OpenSpec setup step for existing projects

## Capabilities

### New Capabilities
- `openspec-setup`: Covers the OpenSpec prerequisite check, schema installation, and config setup during init/upgrade

### Modified Capabilities
- `init-command`: Init now includes an OpenSpec setup step and a prerequisite check for the openspec binary
- `upgrade-command`: Upgrade now includes an OpenSpec setup step and a prerequisite check for the openspec binary

## Impact

- **Code**: `internal/init/init.go` (new step in Run), `internal/init/upgrade.go` (new step in Upgrade), new sub-package `internal/init/openspec/`
- **Dependencies**: Requires `openspec` CLI to be installed on the user's system (runtime dependency, not build dependency)
- **Embedded files**: The littlefactory schema files are embedded in the binary (similar to how skills are embedded)
