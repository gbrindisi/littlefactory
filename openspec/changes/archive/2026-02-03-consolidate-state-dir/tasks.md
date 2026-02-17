## 1. Config Changes

- [x] 1.1 Add `DefaultStateDir` constant to `internal/config/config.go`
- [x] 1.2 Add `StateDir` field to `Config` struct with yaml tag `state_dir`
- [x] 1.3 Initialize `StateDir` to default in `LoadConfig`
- [x] 1.4 Add validation for empty `StateDir` in `validate()`

## 2. Progress File Changes

- [x] 2.1 Update `InitProgressFile` signature to accept `*config.Config`
- [x] 2.2 Change progress file path from `tasks/progress.txt` to `cfg.StateDir/progress.md`
- [x] 2.3 Update header format to "# Little Factory Progress Log" with markdown formatting
- [x] 2.4 Update `AppendSessionToProgress` signature to accept `*config.Config`
- [x] 2.5 Update iteration block format to proper markdown with bold labels
- [x] 2.6 Update `ProgressFilePath` signature and implementation
- [x] 2.7 Update `progress_test.go` to use new signatures and paths

## 3. Metadata Changes

- [x] 3.1 Update `SaveMetadata` signature to accept `*config.Config`
- [x] 3.2 Change metadata file path from `tasks/run_metadata.json` to `cfg.StateDir/run_metadata.json`
- [x] 3.3 Update `metadata_test.go` to use new signature and paths

## 4. JSONTaskSource Changes

- [x] 4.1 Update `NewJSONTaskSource` signature to accept `*config.Config`
- [x] 4.2 Store config in `JSONTaskSource` struct
- [x] 4.3 Use `cfg.StateDir` instead of hardcoded `.littlefactory` for task path
- [x] 4.4 Update `json_test.go` to use new signature

## 5. Driver Integration

- [x] 5.1 Update `Driver.Run` calls to pass config to `InitProgressFile`
- [x] 5.2 Update `Driver.Run` calls to pass config to `SaveMetadata`
- [x] 5.3 Update `Driver.RunIteration` call to pass config to `AppendSessionToProgress`
- [x] 5.4 Update driver instantiation to pass config to `NewJSONTaskSource`

## 6. Cleanup

- [x] 6.1 Update `ProgressFileName` constant from `progress.txt` to `progress.md`
- [x] 6.2 Run all tests to verify changes
