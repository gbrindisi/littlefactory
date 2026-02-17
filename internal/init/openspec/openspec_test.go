package openspec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckInstalled(t *testing.T) {
	t.Run("binary found in PATH", func(t *testing.T) {
		// Create a temp dir with a fake "openspec" binary and prepend it to PATH.
		binDir := t.TempDir()
		fakeBin := filepath.Join(binDir, "openspec")
		if err := os.WriteFile(fakeBin, []byte("#!/bin/sh\n"), 0o755); err != nil {
			t.Fatal(err)
		}

		t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

		if err := CheckInstalled(); err != nil {
			t.Fatalf("expected no error when openspec is in PATH, got: %v", err)
		}
	})

	t.Run("binary not found returns error", func(t *testing.T) {
		// Set PATH to an empty directory so openspec cannot be found.
		emptyDir := t.TempDir()
		t.Setenv("PATH", emptyDir)

		err := CheckInstalled()
		if err == nil {
			t.Fatal("expected error when openspec is not in PATH")
		}
		if got := err.Error(); !strings.Contains(got, "openspec is not installed") {
			t.Errorf("expected descriptive error message, got: %s", got)
		}
	})
}

func TestExtractSchema(t *testing.T) {
	tmpDir := t.TempDir()

	if err := ExtractSchema(tmpDir); err != nil {
		t.Fatalf("ExtractSchema failed: %v", err)
	}

	destRoot := filepath.Join(tmpDir, "openspec", "schemas", "littlefactory")

	// Verify schema.yaml exists
	schemaPath := filepath.Join(destRoot, "schema.yaml")
	info, err := os.Stat(schemaPath)
	if err != nil {
		t.Fatalf("expected schema.yaml to exist: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("expected schema.yaml to have content")
	}

	// Verify template files exist
	expectedTemplates := []string{
		"tasks.md",
		"spec.md",
		"design.md",
		"proposal.md",
		"tasks.json",
	}
	for _, tmpl := range expectedTemplates {
		tmplPath := filepath.Join(destRoot, "templates", tmpl)
		info, err := os.Stat(tmplPath)
		if err != nil {
			t.Errorf("expected template %s to exist: %v", tmpl, err)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("expected template %s to have content", tmpl)
		}
	}
}

func TestExtractSchema_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()

	if err := ExtractSchema(tmpDir); err != nil {
		t.Fatalf("first ExtractSchema failed: %v", err)
	}
	if err := ExtractSchema(tmpDir); err != nil {
		t.Fatalf("second ExtractSchema failed: %v", err)
	}
}

func TestSetup(t *testing.T) {
	t.Run("config created when missing", func(t *testing.T) {
		tmpDir := t.TempDir()

		result, err := Setup(tmpDir)
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}
		if !result.ConfigCreated {
			t.Error("expected ConfigCreated to be true when config was missing")
		}

		// Verify config was created with expected content.
		configPath := filepath.Join(tmpDir, "openspec", "config.yaml")
		content, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("expected config.yaml to exist: %v", err)
		}
		if string(content) != defaultConfig {
			t.Errorf("expected config content %q, got %q", defaultConfig, string(content))
		}

		// Verify schema was also extracted.
		schemaPath := filepath.Join(tmpDir, "openspec", "schemas", "littlefactory", "schema.yaml")
		if _, err := os.Stat(schemaPath); err != nil {
			t.Fatalf("expected schema.yaml to exist after Setup: %v", err)
		}
	})

	t.Run("config preserved when existing", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Pre-create a custom config.
		configDir := filepath.Join(tmpDir, "openspec")
		if err := os.MkdirAll(configDir, 0o755); err != nil {
			t.Fatal(err)
		}
		customConfig := "schema: custom\nfoo: bar\n"
		configPath := filepath.Join(configDir, "config.yaml")
		if err := os.WriteFile(configPath, []byte(customConfig), 0o644); err != nil {
			t.Fatal(err)
		}

		result, err := Setup(tmpDir)
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}
		if result.ConfigCreated {
			t.Error("expected ConfigCreated to be false when config already existed")
		}

		// Verify the existing config was NOT overwritten.
		content, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("expected config.yaml to still exist: %v", err)
		}
		if string(content) != customConfig {
			t.Errorf("expected config to be preserved as %q, got %q", customConfig, string(content))
		}

		// Verify schema was still extracted (schemas should be overwritten).
		schemaPath := filepath.Join(tmpDir, "openspec", "schemas", "littlefactory", "schema.yaml")
		if _, err := os.Stat(schemaPath); err != nil {
			t.Fatalf("expected schema.yaml to exist after Setup: %v", err)
		}
	})
}

