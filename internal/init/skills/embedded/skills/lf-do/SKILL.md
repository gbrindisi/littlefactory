---
name: lf-do
description: Run a littlefactory change in the background and monitor progress. Use to execute a change, kick off implementation, or run littlefactory on a change.
---

Run a littlefactory change as a background process and monitor it to completion.

---

## Step 1: Determine the change name

Check which changes exist:

```bash
ls .littlefactory/changes/
```

- **If exactly one change directory exists:** use it automatically. Tell the user which change you selected.
- **If multiple change directories exist:** list them and ask the user which change to run. Wait for their response before proceeding.
- **If no change directories exist:** tell the user there are no changes to run. Suggest they create one with `/lf-formalize`.

The selected directory name is `<name>` for all subsequent commands.

---

## Step 2: Verify tasks exist

Confirm that `.littlefactory/changes/<name>/tasks.json` exists and contains tasks. If it does not exist, tell the user the change has no tasks and suggest running `/lf-formalize` first.

---

## Step 3: Run littlefactory in the background

Execute the following using the Bash tool with `run_in_background: true`:

```bash
littlefactory run -c <name>
```

Tell the user that littlefactory is now running change `<name>` in the background.

---

## Step 4: Monitor progress

After launching the background process, periodically check progress using:

```bash
littlefactory status -c <name>
```

**Monitoring rules:**
- Wait for the background process notification before checking status -- do NOT poll in a loop or use sleep
- When the background process completes, run the status command one final time to capture the outcome
- If the user asks about progress at any point, run the status command to give them an update

---

## Step 5: Report outcome

When the background process finishes, report the result to the user:

- **All tasks completed:** Tell the user the change finished successfully. Suggest running `/lf-verify` to validate the implementation.
- **Some tasks failed:** Report which tasks failed and any error output. Suggest the user investigate the failures.
- **Process was cancelled:** Acknowledge the cancellation and report how far it got.

Include a brief summary of what was accomplished (number of tasks completed, any notable output).
