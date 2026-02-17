## Context

Littlefactory orchestrates autonomous coding agents to implement openspec changes. Currently, it runs in a single working directory, which limits parallelization. Users want to run multiple implementations concurrently, each in its own git worktree with a dedicated branch.

The workflow is:
1. Design a change using openspec (creates `openspec/changes/<name>/tasks.json`)
2. Commit the change artifacts to main
3. Create a worktree branching from main
4. Run littlefactory in the worktree to implement the change
5. Review, merge, and clean up

## Goals / Non-Goals

**Goals:**
- Enable parallel implementation of multiple changes via git worktrees
- Make worktree creation explicit and safe (clean working tree required)
- Provide visibility into progress across all active changes
- Support existing single-directory workflows without modification

**Non-Goals:**
- Automatic worktree setup (users set up their repo structure)
- Bare-repo worktree style (wtree) - use standard `git worktree`
- Automatic merging or conflict resolution
- Remote worktree management

## Decisions

### Decision: Use `git rev-parse --git-common-dir` for detection
**Rationale**: This is the future-proof way to find the common git directory. Works for normal repos, bare repos, linked worktrees, and custom GIT_DIR setups.

**Alternative considered**: Check for `.git/worktrees/` directory directly. Rejected because it's less portable and doesn't handle all git configurations.

### Decision: `--worktree` flag is explicit creation only
**Rationale**: Making `-w` mean "create new worktree" (error if exists) prevents accidental state confusion. If you want to run in an existing worktree, cd into it and run without `-w`.

**Alternative considered**: "Find or create" semantics where `-w` reuses existing worktrees. Rejected because it's less predictable and masks potential user errors.

### Decision: Require clean working tree for worktree creation
**Rationale**: Ensures the openspec change is committed before branching. The worktree branches from HEAD, so uncommitted changes would be lost or cause confusion.

**Implementation**: Check `git status --porcelain` before creating worktree.

### Decision: Default worktrees_dir is sibling to repo (`..`)
**Rationale**: This matches git's default behavior for `git worktree add`. Users who want a different layout can configure `worktrees_dir` in Factoryfile.

**Alternative considered**: Create worktrees inside the repo. Rejected because it requires bare-repo setup and is more opinionated.

### Decision: Change name determines branch name
**Rationale**: Simple 1:1 mapping between openspec change name and git branch. Reduces configuration and mental overhead.

**Alternative considered**: Separate `--branch` flag. Rejected as unnecessary complexity for v1.

## Risks / Trade-offs

**[Risk]** User forgets to commit openspec change before `-w`
- **Mitigation**: Clear error message prompting to commit first

**[Risk]** Worktree path conflicts with existing directory
- **Mitigation**: Let `git worktree add` handle this - it errors appropriately

**[Risk]** Users confused about which directory they're in
- **Mitigation**: Status command shows all worktrees; clear output about workspace location

**[Trade-off]** Not supporting "find or create" for `-w`
- Simpler mental model but requires explicit cd for existing worktrees
- Acceptable because running in existing worktree is natural (`cd feature-a && lf run`)

## Package Structure

```
internal/
├── worktree/           # NEW: Git worktree operations
│   ├── detect.go       # HasWorktrees(), GetCommonDir()
│   ├── create.go       # Create(), validates clean state
│   ├── list.go         # List() - parse git worktree list
│   └── worktree_test.go
├── config/
│   └── config.go       # Add WorktreesDir field
└── driver/
    └── driver.go       # Add ChangeName, WorktreePath options
```
