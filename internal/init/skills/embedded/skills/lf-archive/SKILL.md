---
name: lf-archive
description: Archive a completed change - merge delta specs into project specs, enrich with boundaries and gotchas, clean up. Use to close out a change after verify.
---

Archive a completed change. Merge delta specs into the project's long-lived specs directory, enrich them with boundaries and gotchas from implementation experience, and optionally clean up the change directory.

**Input**: Optionally specify a change name. If omitted, infer from context or prompt for selection.

---

## Step 1: Determine the change name

Check which changes exist:

```bash
ls .littlefactory/changes/
```

- **If exactly one change directory exists:** use it automatically. Tell the user which change you selected.
- **If multiple change directories exist:** list them and use the **AskUserQuestion tool** to let the user select which change to archive.
- **If no change directories exist:** tell the user there are no changes to archive. Stop here.

The selected directory name is `<name>` for all subsequent steps.

---

## Step 2: Load change artifacts and detect specs_dir

Read all available artifacts from `.littlefactory/changes/<name>/`:

- `proposal.md` -- change proposal (motivation, capabilities)
- `design.md` -- design decisions and rationale
- `specs/*/spec.md` -- delta specs (one per affected capability)
- `tasks.json` -- implementation tasks and their statuses

Read each artifact that exists. Skip any that do not exist.

Also check the conversation history for any `/lf-verify` report that was generated -- its warnings and suggestions are prime candidates for gotchas.

**Detect specs_dir**: Read the `Factoryfile` (or `Factoryfile.yaml`) in the project root. Look for a `specs_dir` configuration key.

- If `specs_dir` is configured: use that path for merging specs.
- If `specs_dir` is NOT configured: ask the user using the **AskUserQuestion tool**:
  > "No `specs_dir` configured in Factoryfile. Where should I merge specs to? (Enter a path, or type 'skip' to skip spec merging)"
  - If the user provides a path, use it.
  - If the user says 'skip', skip Steps 3-6 entirely and go to Step 7.

---

## Step 3: Build merge plan

For each delta spec in `.littlefactory/changes/<name>/specs/`:

1. Read the delta spec file (`specs/<capability>/spec.md`)
2. Check if a corresponding spec exists in `<specs_dir>/<capability>/spec.md`
3. Classify each spec operation:

| Delta Header | Existing Spec? | Action |
|---|---|---|
| ADDED Requirements | No | CREATE new spec file |
| ADDED Requirements | Yes | APPEND requirements to existing spec |
| MODIFIED Requirements | Yes | REPLACE matching requirement blocks |
| REMOVED Requirements | Yes | DELETE matching requirement blocks |

Present the merge plan to the user as a table:

```
## Merge Plan

| Capability | Action | Details |
|---|---|---|
| <capability> | CREATE | New spec with N requirements |
| <capability> | MODIFY | Update N requirements, add M new |
| <capability> | REMOVE | Delete N requirements |
```

Use the **AskUserQuestion tool** to confirm:
> "Proceed with this merge plan? (yes/no/edit)"

- If **yes**: proceed to Step 4.
- If **no**: stop the merge and go to Step 7.
- If **edit**: ask what to change, adjust the plan, and re-confirm.

---

## Step 4: Gotcha mining

Surface candidate gotchas from multiple sources:

1. **Verify findings**: If a `/lf-verify` report exists in conversation history, extract warnings and suggestions as gotcha candidates.
2. **Design surprises**: Review `design.md` for trade-offs, risks, and rejected alternatives that future developers should know about.
3. **Implementation experience**: Based on the conversation context, identify edge cases, workarounds, or unexpected behaviors encountered during implementation.

Present proposed gotchas grouped by capability:

```
## Proposed Gotchas

### <capability>
1. <gotcha description>
   (source: verify warning / design trade-off / implementation experience)
2. <gotcha description>
   (source: ...)
```

Use the **AskUserQuestion tool** to let the user prune:
> "Which gotchas should I persist? (all / none / list numbers to keep, e.g. '1,3')"

Also ask:
> "Any additional gotchas from your experience that I missed?"

If the user provides additional gotchas, add them to the list.

---

## Step 5: Boundary discovery

Review the change's design decisions and implementation to propose boundaries (rules) for affected capabilities:

- **ALWAYS**: Rules that must always be followed (invariants, security constraints)
- **ASK**: Rules where the user should be consulted before deviating
- **NEVER**: Rules that must never be violated (anti-patterns, known-bad approaches)

Present proposed boundaries:

```
## Proposed Boundaries

### <capability>
- ALWAYS: <rule>
- ASK: <rule>
- NEVER: <rule>
```

Use the **AskUserQuestion tool** to let the user prune:
> "Which boundaries should I persist? (all / none / list numbers to keep)"

Also ask:
> "Any additional boundaries I should add?"

---

## Step 6: Write enriched specs

For each capability in the merge plan, write the enriched spec to `<specs_dir>/<capability>/spec.md` using this canonical format:

```markdown
# <capability-name>

## What It Does
One paragraph summary describing the current state of this capability.

## Requirements
### Requirement: <name>
<requirement text with SHALL/MUST statements>

#### Scenario: <name>
- **WHEN** <condition>
- **THEN** <expected outcome>

## Boundaries
- ALWAYS: <rule>
- ASK: <rule>
- NEVER: <rule>

## Gotchas
- <gotcha description>
  (learned: <change-name>, <date>)
```

**Merge rules**:

- **CREATE** (new spec): Write full enriched format with all four sections.
- **MODIFY** (existing spec):
  - Rewrite the "What It Does" paragraph to reflect changes.
  - Merge requirements: replace MODIFIED blocks, append ADDED blocks, remove REMOVED blocks.
  - Strip ADDED/MODIFIED/REMOVED headers -- merged specs use plain `### Requirement:` format.
  - Append new boundaries to existing Boundaries section (create if missing).
  - Append new gotchas to existing Gotchas section (create if missing).
  - If existing spec uses old format (requirements only), add the new sections without breaking existing content.
- **REMOVE**: Delete the matching requirement blocks. If all requirements are removed, capture the removal reason as a gotcha in a related spec or note it for the user.

**Gotcha provenance**: Each gotcha entry MUST include the change name and today's date:
```
- <gotcha description>
  (learned: <change-name>, <date>)
```

After writing, show the user a summary of what was written and where.

---

## Step 7: Optional cleanup

Use the **AskUserQuestion tool** to ask:
> "Spec merge complete. Delete the change directory `.littlefactory/changes/<name>/`? (yes/no)"

- If **yes**: delete the change directory.
- If **no**: leave it intact.

If spec merging was skipped (no specs_dir), still offer cleanup.

---

## Step 8: Offer to commit

Use the **AskUserQuestion tool** to ask:
> "Commit the spec merge and archival? (yes/no)"

- If **yes**: stage all changed/added files under `<specs_dir>/` (and the deletion of `.littlefactory/changes/<name>/` if cleaned up), then create a commit with message: `docs: archive <name> - merge specs and capture learnings`
- If **no**: leave changes unstaged.

---

## Output

Summarize what was done:
- Change name
- Specs merged (list of capabilities and actions taken)
- Gotchas captured (count)
- Boundaries captured (count)
- Whether change directory was cleaned up
- Whether changes were committed
- Suggest next steps if applicable

---

## Guardrails

- **Always interactive** -- propose changes and let the user confirm. Never auto-write specs without confirmation.
- **No references to openspec CLI** -- this skill is self-contained.
- **All change paths under `.littlefactory/changes/<name>/`** -- read from here.
- **All spec output paths under `<specs_dir>/`** -- write enriched specs here.
- **Preserve existing spec content** -- when modifying, merge carefully. Do not lose existing requirements, boundaries, or gotchas.
- **Gotcha provenance is mandatory** -- every gotcha must have `(learned: <change-name>, <date>)`.
- **Handle missing sections gracefully** -- if an existing spec lacks Boundaries or Gotchas sections, create them rather than failing.
- **Do not force enrichment** -- if the user says "none" for gotchas or boundaries, respect that and write specs without those sections.
