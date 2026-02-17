## 1. Add Tasks Validation

- [x] 1.1 Create `ValidateTasks(tasks []Task, filePath string) error` function in `internal/tasks/json.go`
- [x] 1.2 Implement required field validation (id, title, status non-empty)
- [x] 1.3 Implement status value validation (todo, in_progress, done)
- [x] 1.4 Implement unique ID validation with duplicate detection
- [x] 1.5 Implement sequential blocker chain validation (one root, single blockers, all reachable)
- [x] 1.6 Implement multi-error collection and formatted error message
- [x] 1.7 Add unit tests for all validation scenarios

## 2. Add Tasks Flag

- [x] 2.1 Add `--tasks/-t` flag to run command in `cmd/littlefactory/main.go`
- [x] 2.2 Implement flag priority resolution (-t > -c > default)
- [x] 2.3 Add file existence validation for explicit --tasks path
- [x] 2.4 Update `validateChangeFlags` to handle new flag combinations
- [x] 2.5 Add unit tests for flag parsing and priority

## 3. Integrate Validation on Load

- [x] 3.1 Update `NewJSONTaskSource` to call validation after parsing
- [x] 3.2 Update `NewJSONTaskSourceWithPath` to call validation after parsing
- [x] 3.3 Add file existence check to `NewJSONTaskSourceWithPath`
- [x] 3.4 Update constructors to return error (change signature if needed)
- [x] 3.5 Update callers in main.go to handle constructor errors

## 4. Remove Obsolete Skill

- [x] 4.1 Remove `internal/init/skills/embedded/skills/openspec-to-lf/` directory
- [x] 4.2 Remove `.littlefactory/skills/openspec-to-lf/` from project (if exists)
- [x] 4.3 Update any tests that reference the removed skill

## 5. Update Schema

- [x] 5.1 Update `openspec/schemas/littlefactory/schema.yaml` tasks-littlefactory instruction to write only to change directory
- [x] 5.2 Update `internal/init/openspec/embedded/schema/schema.yaml` to match
