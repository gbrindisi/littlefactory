# verifier-template

## What It Does
Embeds a VERIFIER.md agent template in the littlefactory binary at compile time, analogous to WORKER.md. The template instructs verification agents on how to check implementation against change specs across three dimensions (completeness, correctness, coherence) and how to emit remediation tasks when drift is detected.

## Requirements

### Requirement: Embedded verifier template
The system SHALL embed a VERIFIER.md template in the binary at compile time, analogous to the existing WORKER.md template.

#### Scenario: Verifier template available at runtime
- **WHEN** the verify command executes
- **THEN** the embedded VERIFIER.md template is available for rendering

#### Scenario: Local override supported
- **WHEN** a file exists at `<state_dir>/agents/VERIFIER.md`
- **THEN** the system uses the local override instead of the embedded template

### Requirement: Verifier template instructs three-dimension check
The template SHALL instruct the agent to check completeness, correctness, and coherence -- matching the dimensions from the existing `/lf-verify` skill.

#### Scenario: Template references change artifacts
- **WHEN** the template is rendered
- **THEN** it includes paths to proposal.md, specs/, design.md (if present), and tasks.json for the change

#### Scenario: Template instructs remediation task generation
- **WHEN** the template is rendered
- **THEN** it instructs the agent to append remediation tasks to tasks.json if drift is found, following the blocker chain convention

#### Scenario: Template instructs exit code convention
- **WHEN** the template is rendered
- **THEN** it instructs the agent to exit 0 on pass and exit non-zero on drift

### Requirement: Verifier template receives change context
The system SHALL render the verifier template with change name, change path, and paths to all change artifacts.

#### Scenario: Template placeholders rendered
- **WHEN** the verify command renders the template
- **THEN** placeholders for change name, change path, proposal path, specs paths, and design path are replaced with actual values

## Boundaries

## Gotchas
