## Why

Ciccio is a Python-based autonomous agent loop driver that works well but suffers from Python packaging complexity and deployment friction. A Go rewrite as "littlefactory" provides a single binary distribution while establishing a foundation for a lightweight orchestration platform for autonomous agents.

## What Changes

- Rewrite core loop driver (Ralph) from Python to Go
- Remove Rich-based dashboard component (not needed)
- Introduce flexible Agent and TaskSource interfaces for future extensibility
- Add Factoryfile configuration support (YAML-based)
- Maintain full compatibility with beads task system integration
- Preserve metadata tracking and progress logging behavior
- Package as single Go binary (no Python/venv required)

## Capabilities

### New Capabilities
- `agent-interface`: Abstract interface for running autonomous agents (initially Claude Code, extensible to others)
- `task-source-interface`: Abstract interface for task management systems (initially beads, extensible to others)
- `loop-driver`: Core sequential agent execution loop with iteration limits and timeouts
- `metadata-tracking`: Run and iteration metadata capture with JSON serialization
- `progress-logging`: Append-only progress file for session traceability
- `template-system`: Embedded templates with local override support and task injection
- `config-management`: Factoryfile loading with flag override support
- `project-detection`: Automatic project root discovery via .beads directory

### Modified Capabilities
<!-- No existing capabilities are being modified - this is a rewrite -->

## Impact

- Complete rewrite: New Go codebase replaces existing Python implementation
- Dependencies: Removes Python dependencies (click, rich), adds minimal Go dependencies (cobra, uuid)
- Build system: Replaces Python packaging (pyproject.toml, uv) with Go modules
- Distribution: Single binary instead of Python package
- Beads integration: Maintains shell-based integration with bd CLI (no changes to beads)
- Metadata format: Preserves JSON structure for backward compatibility
- Template format: Preserves CLAUDE.md structure and injection variables
