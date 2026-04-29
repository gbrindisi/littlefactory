package worktree

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Seed copies change artifacts and project specs from the source repo into a
// freshly created worktree. This is necessary because `.littlefactory/` is
// commonly gitignored, so `git worktree add` does not carry these files over
// and the agent would have no proposal/specs/tasks.json to work from.
//
// It copies (when present):
//   - <projectRoot>/.littlefactory/changes/<changeName>/ → <worktreePath>/.littlefactory/changes/<changeName>/
//   - <projectRoot>/.littlefactory/specs/                → <worktreePath>/.littlefactory/specs/
//
// Per-run state (run_metadata.json, progress.md) is intentionally not copied;
// it should regenerate fresh in the worktree.
//
// Missing source directories are not an error — Seed silently skips them so
// callers running without a formalized change still work.
func Seed(projectRoot, worktreePath, changeName string) error {
	sources := []struct {
		rel string
	}{
		{filepath.Join(".littlefactory", "changes", changeName)},
		{filepath.Join(".littlefactory", "specs")},
	}

	for _, s := range sources {
		src := filepath.Join(projectRoot, s.rel)
		info, err := os.Stat(src)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return fmt.Errorf("stat %s: %w", src, err)
		}
		if !info.IsDir() {
			continue
		}

		dst := filepath.Join(worktreePath, s.rel)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(dst), err)
		}
		if err := os.CopyFS(dst, os.DirFS(src)); err != nil {
			return fmt.Errorf("copy %s → %s: %w", src, dst, err)
		}
	}

	return nil
}
