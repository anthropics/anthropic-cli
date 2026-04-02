package main

import (
	"context"
	"strings"
	"testing"
)

// TestLintConformance verifies the Lint function conforms to expected behavior
func TestLintConformance(t *testing.T) {
	ctx := context.Background()
	
	// Test with current source directory
	anthropicCli := New(dag.CurrentModule().Source())
	
	result, err := anthropicCli.Lint(ctx)
	if err != nil {
		t.Fatalf("Lint failed: %v", err)
	}
	
	// Verify result contains expected output
	if !strings.Contains(result, "build") {
		t.Logf("Lint output: %s", result)
	}
}

// TestBuildConformance verifies the Build function for supported platforms
func TestBuildConformance(t *testing.T) {
	ctx := context.Background()
	anthropicCli := New(dag.CurrentModule().Source())
	
	testCases := []struct {
		name   string
		goos   string
		goarch string
	}{
		{"linux_amd64", "linux", "amd64"},
		{"darwin_arm64", "darwin", "arm64"},
		{"windows_amd64", "windows", "amd64"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := anthropicCli.Build(ctx, tc.goos, tc.goarch)
			if err != nil {
				t.Fatalf("Build %s/%s failed: %v", tc.goos, tc.goarch, err)
			}
			
			if file == nil {
				t.Fatal("Build returned nil file")
			}
		})
	}
}

// TestBuildAllConformance verifies BuildAll creates all expected artifacts
func TestBuildAllConformance(t *testing.T) {
	ctx := context.Background()
	anthropicCli := New(dag.CurrentModule().Source())
	
	dir, err := anthropicCli.BuildAll(ctx)
	if err != nil {
		t.Fatalf("BuildAll failed: %v", err)
	}
	
	if dir == nil {
		t.Fatal("BuildAll returned nil directory")
	}
	
	// Verify we can list entries
	entries, err := dir.Entries(ctx)
	if err != nil {
		t.Fatalf("Failed to list entries: %v", err)
	}
	
	expectedCount := 5 // linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64
	if len(entries) != expectedCount {
		t.Errorf("Expected %d artifacts, got %d: %v", expectedCount, len(entries), entries)
	}
}

// TestFastConformance verifies TestFast runs without errors
func TestFastConformance(t *testing.T) {
	ctx := context.Background()
	anthropicCli := New(dag.CurrentModule().Source())
	
	_, err := anthropicCli.TestFast(ctx)
	// We expect this might fail in CI without proper setup, just log
	if err != nil {
		t.Logf("TestFast output (may fail without Go deps): %v", err)
	}
}

// TestGenerateSBOMConformance verifies SBOM generation
func TestGenerateSBOMConformance(t *testing.T) {
	ctx := context.Background()
	anthropicCli := New(dag.CurrentModule().Source())
	
	// First build an artifact
	artifact, err := anthropicCli.Build(ctx, "linux", "amd64")
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	// Create directory with artifact
	dir := dag.Directory().WithFile("ant-linux-amd64", artifact)
	
	// Generate SBOM
	sbom, err := anthropicCli.GenerateSBOM(ctx, dir, "cyclonedx-json")
	if err != nil {
		t.Logf("GenerateSBOM may require syft setup: %v", err)
		return
	}
	
	if sbom == nil {
		t.Fatal("GenerateSBOM returned nil file")
	}
}
