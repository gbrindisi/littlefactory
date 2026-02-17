## 1. Package Structure

- [x] 1.1 Create internal/init/ package directory structure
- [x] 1.2 Create internal/init/agentsmd/ package for AGENTS.md handling
- [x] 1.3 Create internal/init/gitignore/ package for .gitignore management
- [x] 1.4 Create internal/init/skills/ package for skill installation

## 2. Embedded Skills

- [x] 2.1 Create embedded/skills/ directory with openspec-to-lf skill
- [x] 2.2 Add go:embed directive for skills in internal/init/skills/embed.go
- [x] 2.3 Implement ExtractSkills function to copy embedded skills to .littlefactory/skills/

## 3. AGENTS.md Setup

- [x] 3.1 Implement Setup function in agentsmd package
- [x] 3.2 Handle case: no AGENTS.md or CLAUDE.md (create AGENTS.md with default content)
- [x] 3.3 Handle case: CLAUDE.md exists without AGENTS.md (rename and symlink)
- [x] 3.4 Handle case: both files exist (merge and symlink)
- [x] 3.5 Handle case: already configured (CLAUDE.md is symlink to AGENTS.md)
- [x] 3.6 Add default AGENTS.md content constant

## 4. Gitignore Management

- [x] 4.1 Implement EnsureEntries function in gitignore package
- [x] 4.2 Add idempotent check for existing entries
- [x] 4.3 Create .gitignore if it does not exist
- [x] 4.4 Append new entries preserving existing content

## 5. Skill Symlinking

- [x] 5.1 Implement CreateSymlinks function for .claude/skills/ integration
- [x] 5.2 Check if .claude/ directory exists
- [x] 5.3 Create .claude/skills/ directory if needed
- [x] 5.4 Create symlinks for each skill (skip if already exists)

## 6. Init Orchestration

- [x] 6.1 Create internal/init/init.go with Run function
- [x] 6.2 Implement step logging with numbered steps and indented sub-operations
- [x] 6.3 Wire up Factoryfile creation, AGENTS.md setup, gitignore, and skills
- [x] 6.4 Update cmd/littlefactory/main.go runInit to use new init package

## 7. Upgrade Command

- [x] 7.1 Add upgradeCmd cobra command in main.go
- [x] 7.2 Implement runUpgrade function that checks for existing Factoryfile
- [x] 7.3 Create internal/init/upgrade.go with Upgrade function
- [x] 7.4 Reuse agentsmd, gitignore, skills packages with idempotent behavior

## 8. Testing

- [x] 8.1 Add unit tests for agentsmd package (all scenarios)
- [x] 8.2 Add unit tests for gitignore package (idempotency, creation)
- [x] 8.3 Add unit tests for skills package (extraction, symlinking)
- [x] 8.4 Add integration test for full init flow
- [x] 8.5 Add integration test for upgrade flow
