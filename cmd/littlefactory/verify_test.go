package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gbrindisi/littlefactory/internal/template"
)

func TestVerifyCmd_HasChangeFlag(t *testing.T) {
	flag := verifyCmd.Flags().Lookup("change")
	if flag == nil {
		t.Fatal("expected --change flag to be registered")
	}
	if flag.Shorthand != "c" {
		t.Errorf("expected shorthand 'c', got '%s'", flag.Shorthand)
	}
}

func TestVerifyCmd_ChangeFlagRequired(t *testing.T) {
	flag := verifyCmd.Flags().Lookup("change")
	if flag == nil {
		t.Fatal("expected --change flag to be registered")
	}
	annotations := flag.Annotations
	if annotations == nil {
		t.Fatal("expected --change flag to have annotations (required)")
	}
	if _, ok := annotations["cobra_annotation_bash_completion_one_required_flag"]; !ok {
		t.Error("expected --change flag to be marked as required")
	}
}

func TestVerifyCmd_AcceptsOptionalAgentArg(t *testing.T) {
	// cobra.MaximumNArgs(1) means 0 or 1 args
	if verifyCmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}
}

func TestBuildChangeContext_FullChange(t *testing.T) {
	tmpDir := t.TempDir()

	// Create change directory with all artifacts
	changePath := filepath.Join(tmpDir, ".littlefactory", "changes", "my-change")
	specsDir := filepath.Join(changePath, "specs", "feature-a")
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create proposal, design, tasks, spec
	for _, f := range []string{"proposal.md", "design.md", "tasks.json"} {
		if err := os.WriteFile(filepath.Join(changePath, f), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(specsDir, "spec.md"), []byte("spec"), 0644); err != nil {
		t.Fatal(err)
	}

	ctx := buildChangeContext(tmpDir, "my-change")

	if ctx.ChangeName != "my-change" {
		t.Errorf("expected ChangeName 'my-change', got %q", ctx.ChangeName)
	}
	if ctx.ChangePath != filepath.Join(".littlefactory", "changes", "my-change") {
		t.Errorf("unexpected ChangePath: %q", ctx.ChangePath)
	}
	if ctx.ProposalPath == "" {
		t.Error("expected ProposalPath to be set")
	}
	if ctx.DesignPath == "" {
		t.Error("expected DesignPath to be set")
	}
	if ctx.TasksPath == "" {
		t.Error("expected TasksPath to be set")
	}
	if !strings.Contains(ctx.SpecsPaths, "feature-a") {
		t.Errorf("expected SpecsPaths to contain feature-a, got %q", ctx.SpecsPaths)
	}
}

func TestBuildChangeContext_EmptyChange(t *testing.T) {
	tmpDir := t.TempDir()

	// Create change directory with no artifacts
	changePath := filepath.Join(tmpDir, ".littlefactory", "changes", "empty-change")
	if err := os.MkdirAll(changePath, 0755); err != nil {
		t.Fatal(err)
	}

	ctx := buildChangeContext(tmpDir, "empty-change")

	if ctx.ChangeName != "empty-change" {
		t.Errorf("expected ChangeName 'empty-change', got %q", ctx.ChangeName)
	}
	if ctx.ProposalPath != "" {
		t.Error("expected ProposalPath to be empty")
	}
	if ctx.DesignPath != "" {
		t.Error("expected DesignPath to be empty")
	}
	if ctx.TasksPath != "" {
		t.Error("expected TasksPath to be empty")
	}
	if ctx.SpecsPaths != "" {
		t.Error("expected SpecsPaths to be empty")
	}
}

func TestBuildChangeContext_MultipleSpecs(t *testing.T) {
	tmpDir := t.TempDir()

	changePath := filepath.Join(tmpDir, ".littlefactory", "changes", "multi-spec")
	for _, spec := range []string{"spec-a", "spec-b"} {
		specDir := filepath.Join(changePath, "specs", spec)
		if err := os.MkdirAll(specDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(specDir, "spec.md"), []byte("spec"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	ctx := buildChangeContext(tmpDir, "multi-spec")

	if !strings.Contains(ctx.SpecsPaths, "spec-a") || !strings.Contains(ctx.SpecsPaths, "spec-b") {
		t.Errorf("expected SpecsPaths to contain both specs, got %q", ctx.SpecsPaths)
	}
	// Should be newline-separated
	parts := strings.Split(ctx.SpecsPaths, "\n")
	if len(parts) != 2 {
		t.Errorf("expected 2 spec paths separated by newline, got %d", len(parts))
	}
}

func TestBuildChangeContext_RendersIntoTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	changePath := filepath.Join(tmpDir, ".littlefactory", "changes", "render-test")
	if err := os.MkdirAll(changePath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(changePath, "proposal.md"), []byte("p"), 0644); err != nil {
		t.Fatal(err)
	}

	ctx := buildChangeContext(tmpDir, "render-test")
	tmpl := "Change: {change_name}, Path: {change_path}, Proposal: {proposal_path}"
	result := template.RenderVerifier(tmpl, ctx)

	if !strings.Contains(result, "render-test") {
		t.Errorf("expected rendered template to contain change name, got %q", result)
	}
	if strings.Contains(result, "{change_name}") {
		t.Error("expected {change_name} placeholder to be replaced")
	}
}
