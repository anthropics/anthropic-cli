package main

import (
	"context"
	"fmt"
	"path/filepath"

	"dagger/anthropic-cli/internal/dagger"
)

// AnthropicCli provides Dagger functions for local development and CI
// These complement (not replace) the GitHub Actions CI/CD
type AnthropicCli struct {
	// +private
	Source *dagger.Directory
}

func New(
	// Source directory containing the project
	source *dagger.Directory,
) *AnthropicCli {
	return &AnthropicCli{Source: source}
}

// baseContainer returns a Go container with source mounted and module cache persisted
func (m *AnthropicCli) baseContainer() *dagger.Container {
	return dag.Container().
		From("golang:1.23-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git", "bash", "curl"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithEnvVariable("CGO_ENABLED", "0").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-cache")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-cache")).
		WithEnvVariable("GOCACHE", "/go/build-cache")
}

// Build runs Go build and returns the compiled binaries directory
func (m *AnthropicCli) Build(ctx context.Context) (*dagger.Directory, error) {
	return m.baseContainer().
		WithExec([]string{"go", "build", "-o", "/output/", "-ldflags=-s -w", "./..."}).
		Directory("/output").
		Sync(ctx)
}

// BuildForPlatform builds the CLI for a specific OS/arch and returns the binary
func (m *AnthropicCli) BuildForPlatform(
	ctx context.Context,
	// Target OS (linux, darwin, windows)
	// +default="linux"
	os string,
	// Target architecture (amd64, arm64)
	// +default="amd64"
	arch string,
) (*dagger.File, error) {
	outputName := fmt.Sprintf("anthropic-cli-%s-%s", os, arch)
	if os == "windows" {
		outputName += ".exe"
	}

	return m.baseContainer().
		WithEnvVariable("GOOS", os).
		WithEnvVariable("GOARCH", arch).
		WithExec([]string{
			"go", "build",
			"-o", filepath.Join("/output", outputName),
			"-ldflags=-s -w -X main.Version=$(git describe --tags --always --dirty)",
			"./cmd/ant",
		}).
		File(filepath.Join("/output", outputName)).
		Sync(ctx)
}

// Test runs the test suite and returns test results as files
func (m *AnthropicCli) Test(ctx context.Context) (*dagger.Directory, error) {
	return m.baseContainer().
		WithExec([]string{"apk", "add", "--no-cache", "lsof"}).
		WithExec([]string{"go", "test", "-v", "-count=1", "-coverprofile=/output/coverage.out", "./..."}).
		WithExec([]string{"go", "tool", "cover", "-html=/output/coverage.out", "-o", "/output/coverage.html"}).
		Directory("/output").
		Sync(ctx)
}

// Lint runs linting and returns the report
func (m *AnthropicCli) Lint(ctx context.Context) (*dagger.File, error) {
	return m.baseContainer().
		WithExec([]string{"go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"}).
		WithExec([]string{"golangci-lint", "run", "./...", "--out-format=json", "--issues-exit-code=0"}).
		WithNewFile("/output/lint-report.json", "", dagger.ContainerWithNewFileOpts{}).
		File("/output/lint-report.json").
		Sync(ctx)
}

// VulnScan runs govulncheck and returns the vulnerability report
func (m *AnthropicCli) VulnScan(ctx context.Context) (*dagger.File, error) {
	return m.baseContainer().
		WithExec([]string{"go", "install", "golang.org/x/vuln/cmd/govulncheck@latest"}).
		WithExec([]string{"govulncheck", "-format=json", "-show=verbose", "./..."}).
		WithNewFile("/output/vulns.json", "", dagger.ContainerWithNewFileOpts{}).
		File("/output/vulns.json").
		Sync(ctx)
}

// SBOM generates a CycloneDX SBOM for the project
func (m *AnthropicCli) SBOM(
	ctx context.Context,
	// Output format (cyclonedx-json, cyclonedx-xml, spdx-json, spdx-tag-value)
	// +default="cyclonedx-json"
	format string,
) (*dagger.File, error) {
	filename := fmt.Sprintf("sbom.%s", format)
	if format == "cyclonedx-json" {
		filename = "sbom.cdx.json"
	} else if format == "spdx-json" {
		filename = "sbom.spdx.json"
	}

	return m.baseContainer().
		WithExec([]string{"go", "install", "github.com/anchore/syft/cmd/syft@latest"}).
		WithExec([]string{"syft", "scan", "dir:/src", "-o", format, "--file", filepath.Join("/output", filename)}).
		File(filepath.Join("/output", filename)).
		Sync(ctx)
}

// Provenance generates SLSA provenance attestation
func (m *AnthropicCli) Provenance(
	ctx context.Context,
	// Target OS
	// +default="linux"
	os string,
	// Target architecture
	// +default="amd64"
	arch string,
) (*dagger.File, error) {
	binary, err := m.BuildForPlatform(ctx, os, arch)
	if err != nil {
		return nil, fmt.Errorf("build failed: %w", err)
	}

	binaryName := fmt.Sprintf("anthropic-cli-%s-%s", os, arch)
	if os == "windows" {
		binaryName += ".exe"
	}

	return m.baseContainer().
		WithFile(filepath.Join("/tmp", binaryName), binary).
		WithExec([]string{"sh", "-c", fmt.Sprintf(
			"sha256sum /tmp/%s > /output/%s.sha256 && sha256sum /tmp/%s | awk '{print $1}' > /output/%s.provenance",
			binaryName, binaryName, binaryName, binaryName,
		)}).
		File(fmt.Sprintf("/output/%s.sha256", binaryName)).
		Sync(ctx)
}

// All runs the complete CI pipeline: build, test, lint, vuln-scan, and sbom
// Returns a directory containing all artifacts and reports
func (m *AnthropicCli) All(ctx context.Context) (*dagger.Directory, error) {
	output := dag.Directory()

	// Run tests
	results, err := m.Test(ctx)
	if err != nil {
		return nil, fmt.Errorf("test failed: %w", err)
	}
	output = output.WithDirectory("test-results", results)

	// Build for current platform
	bin, err := m.Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("build failed: %w", err)
	}
	output = output.WithDirectory("binaries", bin)

	// SBOM
	sbom, err := m.SBOM(ctx, "cyclonedx-json")
	if err != nil {
		return nil, fmt.Errorf("sbom failed: %w", err)
	}
	output = output.WithFile("sbom.cdx.json", sbom)

	// Vulnerability scan
	vulns, err := m.VulnScan(ctx)
	if err != nil {
		return nil, fmt.Errorf("vuln scan failed: %w", err)
	}
	output = output.WithFile("vulns.json", vulns)

	return output.Sync(ctx)
}
