package gitignore

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureEntries_NoGitignore(t *testing.T) {
	dir := t.TempDir()

	result, err := EnsureEntries(dir)
	if err != nil {
		t.Fatalf("EnsureEntries failed: %v", err)
	}

	if result.Action != ActionCreated {
		t.Errorf("expected action %q, got %q", ActionCreated, result.Action)
	}

	if len(result.Added) != len(RequiredEntries) {
		t.Errorf("expected %d added entries, got %d", len(RequiredEntries), len(result.Added))
	}

	// Verify .gitignore was created with all required entries
	content, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	for _, entry := range RequiredEntries {
		if !strings.Contains(string(content), entry) {
			t.Errorf(".gitignore missing entry %q", entry)
		}
	}

	// Verify file ends with newline
	if len(content) > 0 && content[len(content)-1] != '\n' {
		t.Error(".gitignore does not end with newline")
	}
}

func TestEnsureEntries_ExistingGitignoreNoOverlap(t *testing.T) {
	dir := t.TempDir()

	existingContent := "node_modules/\n*.log\n"
	gitignorePath := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(existingContent), 0o644); err != nil {
		t.Fatalf("failed to create .gitignore: %v", err)
	}

	result, err := EnsureEntries(dir)
	if err != nil {
		t.Fatalf("EnsureEntries failed: %v", err)
	}

	if result.Action != ActionAdded {
		t.Errorf("expected action %q, got %q", ActionAdded, result.Action)
	}

	if len(result.Added) != len(RequiredEntries) {
		t.Errorf("expected %d added entries, got %d", len(RequiredEntries), len(result.Added))
	}

	// Verify existing content is preserved
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	if !strings.HasPrefix(string(content), existingContent) {
		t.Error("existing .gitignore content was not preserved")
	}

	for _, entry := range RequiredEntries {
		if !strings.Contains(string(content), entry) {
			t.Errorf(".gitignore missing entry %q", entry)
		}
	}
}

func TestEnsureEntries_AllEntriesExist(t *testing.T) {
	dir := t.TempDir()

	var content string
	for _, entry := range RequiredEntries {
		content += entry + "\n"
	}
	gitignorePath := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create .gitignore: %v", err)
	}

	result, err := EnsureEntries(dir)
	if err != nil {
		t.Fatalf("EnsureEntries failed: %v", err)
	}

	if result.Action != ActionSkipped {
		t.Errorf("expected action %q, got %q", ActionSkipped, result.Action)
	}

	if len(result.Skipped) != len(RequiredEntries) {
		t.Errorf("expected %d skipped entries, got %d", len(RequiredEntries), len(result.Skipped))
	}

	// Verify content unchanged
	readContent, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	if string(readContent) != content {
		t.Error(".gitignore content was modified when all entries already existed")
	}
}

func TestEnsureEntries_PartialOverlap(t *testing.T) {
	dir := t.TempDir()

	// Only first entry exists
	content := RequiredEntries[0] + "\n"
	gitignorePath := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create .gitignore: %v", err)
	}

	result, err := EnsureEntries(dir)
	if err != nil {
		t.Fatalf("EnsureEntries failed: %v", err)
	}

	if result.Action != ActionAdded {
		t.Errorf("expected action %q, got %q", ActionAdded, result.Action)
	}

	if len(result.Added) != 1 {
		t.Errorf("expected 1 added entry, got %d", len(result.Added))
	}

	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped entry, got %d", len(result.Skipped))
	}

	// Verify the second entry was added
	readContent, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	for _, entry := range RequiredEntries {
		if !strings.Contains(string(readContent), entry) {
			t.Errorf(".gitignore missing entry %q", entry)
		}
	}
}

func TestEnsureEntries_Idempotent(t *testing.T) {
	dir := t.TempDir()

	// First run: creates .gitignore
	result1, err := EnsureEntries(dir)
	if err != nil {
		t.Fatalf("first EnsureEntries failed: %v", err)
	}
	if result1.Action != ActionCreated {
		t.Errorf("first run: expected action %q, got %q", ActionCreated, result1.Action)
	}

	content1, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatalf("failed to read .gitignore after first run: %v", err)
	}

	// Second run: should skip
	result2, err := EnsureEntries(dir)
	if err != nil {
		t.Fatalf("second EnsureEntries failed: %v", err)
	}
	if result2.Action != ActionSkipped {
		t.Errorf("second run: expected action %q, got %q", ActionSkipped, result2.Action)
	}

	content2, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatalf("failed to read .gitignore after second run: %v", err)
	}

	// Content should be identical after both runs
	if string(content1) != string(content2) {
		t.Errorf("content changed after idempotent second run.\nAfter first:\n%s\nAfter second:\n%s", content1, content2)
	}
}

func TestEnsureEntries_ExistingFileNoTrailingNewline(t *testing.T) {
	dir := t.TempDir()

	// Existing file without trailing newline
	existingContent := "node_modules/"
	gitignorePath := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(existingContent), 0o644); err != nil {
		t.Fatalf("failed to create .gitignore: %v", err)
	}

	result, err := EnsureEntries(dir)
	if err != nil {
		t.Fatalf("EnsureEntries failed: %v", err)
	}

	if result.Action != ActionAdded {
		t.Errorf("expected action %q, got %q", ActionAdded, result.Action)
	}

	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	// Should not have entries on the same line as existing content
	lines := strings.Split(string(content), "\n")
	if lines[0] != "node_modules/" {
		t.Errorf("first line should be original content, got %q", lines[0])
	}

	// Verify all required entries are on their own lines
	for _, entry := range RequiredEntries {
		found := false
		for _, line := range lines {
			if strings.TrimSpace(line) == entry {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("entry %q not found on its own line", entry)
		}
	}
}

func TestEnsureEntries_CommentsAndBlankLinesPreserved(t *testing.T) {
	dir := t.TempDir()

	existingContent := "# Build artifacts\n*.o\n\n# Dependencies\nvendor/\n"
	gitignorePath := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(existingContent), 0o644); err != nil {
		t.Fatalf("failed to create .gitignore: %v", err)
	}

	_, err := EnsureEntries(dir)
	if err != nil {
		t.Fatalf("EnsureEntries failed: %v", err)
	}

	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	// Existing content should be preserved verbatim at the start
	if !strings.HasPrefix(string(content), existingContent) {
		t.Error("existing content with comments and blank lines was not preserved")
	}
}

func TestEnsureEntries_EntryWithWhitespace(t *testing.T) {
	dir := t.TempDir()

	// Entry exists with surrounding whitespace
	content := "  .littlefactory/run_metadata.json  \n.littlefactory/tasks.json\n"
	gitignorePath := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create .gitignore: %v", err)
	}

	result, err := EnsureEntries(dir)
	if err != nil {
		t.Fatalf("EnsureEntries failed: %v", err)
	}

	// Should recognize entries even with surrounding whitespace
	if result.Action != ActionSkipped {
		t.Errorf("expected action %q, got %q (entries with whitespace not recognized)", ActionSkipped, result.Action)
	}
}

func TestRequiredEntries(t *testing.T) {
	if len(RequiredEntries) == 0 {
		t.Fatal("RequiredEntries should not be empty")
	}

	for _, entry := range RequiredEntries {
		if !strings.HasPrefix(entry, ".littlefactory/") {
			t.Errorf("entry %q does not start with .littlefactory/", entry)
		}
	}

	expected := map[string]bool{
		".littlefactory/run_metadata.json": true,
		".littlefactory/tasks.json":        true,
	}
	for _, entry := range RequiredEntries {
		if !expected[entry] {
			t.Errorf("unexpected entry %q in RequiredEntries", entry)
		}
	}
}
