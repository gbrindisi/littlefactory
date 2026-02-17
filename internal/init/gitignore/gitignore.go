// Package gitignore manages .gitignore entries for littlefactory projects,
// ensuring required patterns are present with idempotent updates.
package gitignore

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RequiredEntries lists the gitignore patterns that must be present for
// littlefactory runtime files.
var RequiredEntries = []string{
	".littlefactory/run_metadata.json",
	".littlefactory/tasks.json",
}

// Action describes what the EnsureEntries function did.
type Action string

const (
	// ActionAdded indicates new entries were appended to .gitignore.
	ActionAdded Action = "added"
	// ActionCreated indicates .gitignore was created with required entries.
	ActionCreated Action = "created"
	// ActionSkipped indicates all entries were already present.
	ActionSkipped Action = "skipped"
)

// Result describes the outcome of an EnsureEntries call.
type Result struct {
	Action  Action
	Message string
	Added   []string
	Skipped []string
}

// EnsureEntries ensures all RequiredEntries are present in the project's
// .gitignore file. It creates the file if it doesn't exist, and appends
// missing entries without modifying existing content. The operation is
// idempotent: running it multiple times produces the same result.
func EnsureEntries(projectRoot string) (Result, error) {
	gitignorePath := filepath.Join(projectRoot, ".gitignore")

	existing, err := readEntries(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return Result{}, fmt.Errorf("reading .gitignore: %w", err)
	}

	fileExists := err == nil

	var added, skipped []string
	for _, entry := range RequiredEntries {
		if existing[entry] {
			skipped = append(skipped, entry)
		} else {
			added = append(added, entry)
		}
	}

	if len(added) == 0 {
		return Result{
			Action:  ActionSkipped,
			Message: "all entries already present",
			Skipped: skipped,
		}, nil
	}

	if err := appendEntries(gitignorePath, fileExists, added); err != nil {
		return Result{}, err
	}

	action := ActionAdded
	msg := fmt.Sprintf("added %d entries", len(added))
	if !fileExists {
		action = ActionCreated
		msg = fmt.Sprintf("created .gitignore with %d entries", len(added))
	}

	return Result{
		Action:  action,
		Message: msg,
		Added:   added,
		Skipped: skipped,
	}, nil
}

// readEntries reads a .gitignore file and returns a set of its non-empty,
// non-comment lines (trimmed of whitespace). Returns os.ErrNotExist if the
// file does not exist.
func readEntries(path string) (map[string]bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	entries := make(map[string]bool)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		entries[line] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning .gitignore: %w", err)
	}

	return entries, nil
}

// appendEntries appends the given entries to the .gitignore file. If the file
// already exists and does not end with a newline, a newline is prepended to
// ensure clean formatting.
func appendEntries(path string, fileExists bool, entries []string) error {
	var content string

	if fileExists {
		// Ensure existing file ends with a newline before appending.
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading .gitignore for newline check: %w", err)
		}
		if len(data) > 0 && data[len(data)-1] != '\n' {
			content = "\n"
		}
	}

	for _, entry := range entries {
		content += entry + "\n"
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("opening .gitignore for append: %w", err)
	}
	defer func() { _ = f.Close() }()

	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("writing to .gitignore: %w", err)
	}

	return nil
}
