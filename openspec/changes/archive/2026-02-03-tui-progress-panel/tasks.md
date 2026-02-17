## 1. Add fsnotify dependency

- [x] 1.1 Add `github.com/fsnotify/fsnotify` to go.mod

## 2. Update TUI constructor

- [x] 2.1 Update `tui.New()` signature to accept `*config.Config` and `projectRoot string`
- [x] 2.2 Store config and projectRoot in Model struct
- [x] 2.3 Add `progressFilePath` field computed from `cfg.StateDir`
- [x] 2.4 Update `cmd/littlefactory/main.go` to pass config and projectRoot to TUI

## 3. Implement file watching

- [x] 3.1 Create `FileChangedMsg` type for file change notifications
- [x] 3.2 Implement `watchProgressFile()` function that returns tea.Cmd
- [x] 3.3 Start file watcher in `Init()` method
- [x] 3.4 Handle `FileChangedMsg` in `Update()` to re-read file and refresh viewport

## 4. Replace output buffer with file content

- [x] 4.1 Remove `outputBuf bytes.Buffer` field from Model
- [x] 4.2 Add `progressContent string` field to hold file content
- [x] 4.3 Implement `loadProgressFile()` method to read file content
- [x] 4.4 Call `loadProgressFile()` on startup in `Init()` for initial content
- [x] 4.5 Update `View()` to use `progressContent` instead of `outputBuf`

## 5. Remove cursor navigation

- [x] 5.1 Remove `cursor int` field from Model
- [x] 5.2 Remove `j` and `k` key handlers from `Update()`
- [x] 5.3 Update `renderTasksPanel()` to remove cursor parameter and highlighting
- [x] 5.4 Update status bar to remove j/k from keyboard hints

## 6. Handle missing file gracefully

- [x] 6.1 Show placeholder message when progress.md does not exist
- [x] 6.2 Handle file deletion by reverting to placeholder
- [x] 6.3 Handle file creation by loading content and starting normal display
