## 1. JSON Task Source Implementation

- [x] 1.1 Create `internal/tasks/json.go` with JSONTaskSource struct and JSON file parsing
- [x] 1.2 Implement Ready() - return first task with status "todo"
- [x] 1.3 Implement List() - return all tasks from JSON
- [x] 1.4 Implement Show(id) - find and return task by ID
- [x] 1.5 Implement Claim(id) - set status to "in_progress" and write file
- [x] 1.6 Implement Close(id, reason) - set status to "done" and write file
- [x] 1.7 Implement Reset(id) - set status to "todo" and write file
- [x] 1.8 Add directory creation logic for `.littlefactory/` on write

## 2. Update TaskSource Interface

- [x] 2.1 Update `internal/tasks/source.go` - add Claim() and Reset() methods to interface
- [x] 2.2 Remove Sync() method from interface (no longer needed)
- [x] 2.3 Update Task struct status field documentation (todo/in_progress/done)

## 3. Update Driver for State Management

- [x] 3.1 Update `internal/driver/driver.go` - call Claim() before running agent iteration
- [x] 3.2 Update driver - call Close() on successful iteration (exit code 0)
- [x] 3.3 Update driver - call Reset() on failed iteration or timeout

## 4. Update Project Detection

- [x] 4.1 Update `internal/config/project.go` - detect project root via Factoryfile instead of .beads

## 5. Update Main Entry Point

- [x] 5.1 Update `cmd/littlefactory/main.go` - remove bd CLI availability check
- [x] 5.2 Update main.go - initialize JSONTaskSource instead of BeadsClient
- [x] 5.3 Update main.go - pass project root to JSONTaskSource for tasks.json path

## 6. Update Embedded Template

- [x] 6.1 Update `internal/template/embedded/CLAUDE.md` - remove all bd command references
- [x] 6.2 Update template - simplify workflow to focus on implementation only
- [x] 6.3 Update template - add note that task completion is automatic

## 7. Cleanup

- [x] 7.1 Remove `internal/tasks/beads.go` (BeadsClient implementation)
- [x] 7.2 Update or remove references to beads in any remaining code

## 8. Claude Skill

- [x] 8.1 Rename `.claude/skills/openspec-to-beads/` to `openspec-to-lf/`
- [x] 8.2 Update skill to output `.littlefactory/tasks.json` instead of calling bd commands
- [x] 8.3 Update skill output format to match JSON schema (id, title, description, status)
