---
name: lf:formalize
description: Formalize a change - derive name from conversation, generate all artifacts (proposal, specs, design, tasks.json) in one shot. Use to turn a discussion into a structured change with implementation tasks.
---

Formalize a change. Derive the change name from conversation context. Generate all artifacts sequentially in one shot.

**Input**: No arguments required. The change name and all content are derived from conversation context.

---

## Step 1: Derive Change Name

From the conversation context, derive a kebab-case change name that captures the essence of what's being built or fixed (e.g., "add user authentication" becomes `add-user-auth`).

If the conversation context is insufficient to determine what the change is about, use the **AskUserQuestion tool** (open-ended) to ask:
> "What change do you want to formalize? Describe what you want to build or fix."

**IMPORTANT**: Do NOT proceed without a clear understanding of what the change is.

## Step 2: Create Change Directory

```bash
mkdir -p .littlefactory/changes/<name>/specs
```

If a change with that name already exists in `.littlefactory/changes/`, inform the user and ask whether to overwrite or pick a different name.

## Step 3: Generate proposal.md

**Output path**: `.littlefactory/changes/<name>/proposal.md`

Use this template as the structure for the file:

```markdown
## Why

<!-- Explain the motivation for this change. What problem does this solve? Why now? -->

## What Changes

<!-- Describe what will change. Be specific about new capabilities, modifications, or removals. -->

## Capabilities

### New Capabilities
<!-- Capabilities being introduced. Replace <name> with kebab-case identifier (e.g., user-auth, data-export). Each creates specs/<name>/spec.md -->
- `<name>`: <brief description of what this capability covers>

### Modified Capabilities
<!-- Existing capabilities whose REQUIREMENTS are changing (not just implementation).
     Only list here if spec-level behavior changes. Each needs a delta spec file.
     Leave empty if no requirement changes. -->
- `<existing-name>`: <what requirement is changing>

## Impact

<!-- Affected code, APIs, dependencies, systems -->
```

**Generation instructions**:

Create the proposal document that establishes WHY this change is needed.

Sections:

- **Why**: 1-2 sentences on the problem or opportunity. What problem does this solve? Why now?
- **What Changes**: Bullet list of changes. Be specific about new capabilities, modifications, or removals. Mark breaking changes with **BREAKING**.
- **Capabilities**: Identify which specs will be created or modified:
  - **New Capabilities**: List capabilities being introduced. Each becomes a new `specs/<name>/spec.md`. Use kebab-case names.
  - **Modified Capabilities**: List existing capabilities whose REQUIREMENTS are changing. Only include if spec-level behavior changes (not just implementation details). Each needs a delta spec file. Leave empty if no requirement changes.
- **Impact**: Affected code, APIs, dependencies, or systems.

IMPORTANT: The Capabilities section is critical. It creates the contract between proposal and specs phases. Research existing specs before filling this in. Each capability listed here will need a corresponding spec file.

Keep it concise (1-2 pages). Focus on the "why" not the "how" -- implementation details belong in design.md.

This is the foundation -- specs, design, and tasks all build on this.

**After writing**: Verify the file exists, then proceed.

## Step 4: Generate specs

**Output path**: `.littlefactory/changes/<name>/specs/<capability>/spec.md` (one per capability from the proposal)

Use this template as the structure for each spec file:

```markdown
## ADDED Requirements

### Requirement: <!-- requirement name -->
<!-- requirement text -->

#### Scenario: <!-- scenario name -->
- **WHEN** <!-- condition -->
- **THEN** <!-- expected outcome -->
```

**Generation instructions**:

Create specification files that define WHAT the system should do.

Create one spec file per capability listed in the proposal's Capabilities section.
- New capabilities: use the exact kebab-case name from the proposal (`specs/<capability>/spec.md`).
- Modified capabilities: use the existing spec folder name when creating the delta spec at `specs/<capability>/spec.md`.

Delta operations (use ## headers):
- **ADDED Requirements**: New capabilities
- **MODIFIED Requirements**: Changed behavior -- MUST include full updated content
- **REMOVED Requirements**: Deprecated features -- MUST include **Reason** and **Migration**
- **RENAMED Requirements**: Name changes only -- use FROM:/TO: format

Format requirements:
- Each requirement: `### Requirement: <name>` followed by description
- Use SHALL/MUST for normative requirements (avoid should/may)
- Each scenario: `#### Scenario: <name>` with WHEN/THEN format
- **CRITICAL**: Scenarios MUST use exactly 4 hashtags (`####`). Using 3 hashtags or bullets will fail silently.
- Every requirement MUST have at least one scenario.

MODIFIED requirements workflow:
1. Locate the existing requirement in the relevant spec file
2. Copy the ENTIRE requirement block (from `### Requirement:` through all scenarios)
3. Paste under `## MODIFIED Requirements` and edit to reflect new behavior
4. Ensure header text matches exactly (whitespace-insensitive)

Common pitfall: Using MODIFIED with partial content loses detail. If adding new concerns without changing existing behavior, use ADDED instead.

Specs should be testable -- each scenario is a potential test case.

**After writing**: Read the proposal again to verify all capabilities have spec files, then proceed.

## Step 5: Generate design.md (conditional)

**Output path**: `.littlefactory/changes/<name>/design.md`

**Skip this step** if NONE of the following apply:
- Cross-cutting change (multiple services/modules) or new architectural pattern
- New external dependency or significant data model changes
- Security, performance, or migration complexity
- Ambiguity that benefits from technical decisions before coding

If skipping, move directly to Step 6.

Use this template as the structure for the file:

```markdown
## Context

<!-- Background and current state -->

## Goals / Non-Goals

**Goals:**
<!-- What this design aims to achieve -->

**Non-Goals:**
<!-- What is explicitly out of scope -->

## Decisions

<!-- Key design decisions and rationale -->

## Risks / Trade-offs

<!-- Known risks and trade-offs -->
```

**Generation instructions**:

Create the design document that explains HOW to implement the change.

Sections:

- **Context**: Background, current state, constraints, stakeholders
- **Goals / Non-Goals**: What this design achieves and explicitly excludes
- **Decisions**: Key technical choices with rationale (why X over Y?). Include alternatives considered for each decision.
- **Risks / Trade-offs**: Known limitations, things that could go wrong. Format: [Risk] -> Mitigation
- **Migration Plan**: Steps to deploy, rollback strategy (if applicable)
- **Open Questions**: Outstanding decisions or unknowns to resolve

Focus on architecture and approach, not line-by-line implementation. Reference the proposal for motivation and specs for requirements.

Good design docs explain the "why" behind technical decisions.

**After writing**: Verify the file exists, then proceed.

## Step 6: Generate tasks.json

**Output path**: `.littlefactory/changes/<name>/tasks.json`

This is the final artifact. It goes directly from the proposal, specs, and design (if present) to tasks.json. There is no intermediate tasks.md.

**Read all previously generated artifacts** (proposal.md, all spec files, design.md if it exists) before generating tasks.

### Task Object Structure

```json
{
  "id": "<change-name>-<random-3char>",
  "title": "Section Title",
  "description": "<fat context - see below>",
  "status": "todo",
  "labels": ["change-<change-name>", "<section-label-kebab>"],
  "blockers": []
}
```

### Rules

1. **Group related work into tasks** -- each task should be completable in one focused session
2. **Do NOT create micro-tasks** -- group related checklist items into coherent units of work
3. **Sequential blockers** -- Task N+1 is blocked by Task N (linear chain)
4. **Unique IDs** -- Format: `<change-name>-<random-3char>` (lowercase alphanumeric suffix)
5. **Two labels per task**: `change-<change-name>` and `<section-title-in-kebab-case>`

### Fat Context Description Format

Each task's `description` field MUST contain this structure:

```
## Context
2-5 sentences explaining what this task belongs to and why.
Pull motivation from proposal.md, approach from design.md.

## Scope
- 3-7 bullets summarizing what this task delivers

## Checklist
- [ ] Concrete item 1
- [ ] Concrete item 2
...

## Implementation plan
Break the checklist into 3-8 actionable steps an agent can execute.
Mention concrete files/artifacts to produce.

## Acceptance criteria
- Checkable items (build passes, tests pass, behavior works)
- All checklist items above are implemented and verified

## Key references
**Change:** <change-name>
**Change path:** .littlefactory/changes/<change-name>/

- .littlefactory/changes/<change-name>/proposal.md (section: <heading>)
- .littlefactory/changes/<change-name>/design.md (section: <heading>)
- .littlefactory/changes/<change-name>/specs/<capability>/spec.md (section: <heading>)
```

### Blocker Wiring

- First task: `"blockers": []`
- Task 2: `"blockers": ["<task-1-id>"]`
- Task 3: `"blockers": ["<task-2-id>"]`
- Continue linear chain for all tasks

### Final JSON Format

```json
{
  "tasks": [
    { ... task 1 ... },
    { ... task 2 ... }
  ]
}
```

Ensure valid JSON with 2-space indentation.

**After writing**: Verify the file contains valid JSON, then proceed to output.

---

## Output

After completing all artifacts, summarize:
- Change name and location (`.littlefactory/changes/<name>/`)
- List of artifacts created with brief descriptions
- Number of tasks generated
- "All artifacts created. Ready for implementation."
- Prompt: "Run `/lf:do` or ask me to implement to start working on the tasks."

---

## Guardrails

- **Derive everything from context** -- do not ask for arguments, the change name comes from the conversation
- **Create ALL artifacts** needed for implementation (proposal, specs, optional design, tasks.json)
- **Always read dependency artifacts** before creating the next one
- **No references to openspec CLI** -- this skill is self-contained
- **All paths under `.littlefactory/changes/<name>/`** -- never write to `openspec/`
- **If context is critically unclear, ask** -- but prefer making reasonable decisions to keep momentum
- **Verify each artifact file exists** after writing before proceeding to the next
- **tasks.json must be valid JSON** -- verify before finishing
