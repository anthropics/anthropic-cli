package main

import (
	"context"
	"fmt"

	"dagger/anthropic-cli/internal/dagger"
)

type AnthropicCli struct {
	// Source is the source code directory
	Source *dagger.Directory

	// Git token for accessing private repositories (optional)
	// +optional
	GitToken *dagger.Secret
}

func New(
	// Source directory containing the Go project (defaults to current module source)
	source *dagger.Directory,
	// Git token for accessing private Go modules
	// +optional
	gitToken *dagger.Secret,
) *AnthropicCli {
	return &AnthropicCli{Source: source, GitToken: gitToken}
}

// Lint runs Go build to check for compilation errors
func (m *AnthropicCli) Lint(ctx context.Context) (string, error) {
	ctr := dag.Container().
		From("golang:1.26-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go")

	if m.GitToken != nil {
		ctr = ctr.WithSecretVariable("GIT_TOKEN", m.GitToken).
			WithExec([]string{"sh", "-c", "git config --global url.\"https://x-access-token:${GIT_TOKEN}@github.com/\".insteadOf \"https://github.com/\""})
	}

	return ctr.WithExec([]string{"go", "build", "./..."}).Stdout(ctx)
}

// Test runs the full test suite
func (m *AnthropicCli) Test(ctx context.Context) (string, error) {
	ctr := dag.Container().
		From("golang:1.26-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go")

	if m.GitToken != nil {
		ctr = ctr.WithSecretVariable("GIT_TOKEN", m.GitToken).
			WithExec([]string{"sh", "-c", "git config --global url.\"https://x-access-token:${GIT_TOKEN}@github.com/\".insteadOf \"https://github.com/\""})
	}

	return ctr.WithExec([]string{"go", "test", "-v", "./..."}).Stdout(ctx)
}

// TestFast runs short unit tests only
func (m *AnthropicCli) TestFast(ctx context.Context) (string, error) {
	ctr := dag.Container().
		From("golang:1.26-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go")

	if m.GitToken != nil {
		ctr = ctr.WithSecretVariable("GIT_TOKEN", m.GitToken).
			WithExec([]string{"sh", "-c", "git config --global url.\"https://x-access-token:${GIT_TOKEN}@github.com/\".insteadOf \"https://github.com/\""})
	}

	return ctr.WithExec([]string{"go", "test", "-short", "./..."}).Stdout(ctx)
}

// Build creates a cross-compiled binary for a specific platform
func (m *AnthropicCli) Build(
	ctx context.Context,
	// Target OS (linux, darwin, windows)
	goos string,
	// Target architecture (amd64, arm64)
	goarch string,
) (*dagger.File, error) {
	output := fmt.Sprintf("bin/ant-%s-%s", goos, goarch)
	ctr := dag.Container().
		From("golang:1.26-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithEnvVariable("GOOS", goos).
		WithEnvVariable("GOARCH", goarch).
		WithEnvVariable("CGO_ENABLED", "0").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go")

	if m.GitToken != nil {
		ctr = ctr.WithSecretVariable("GIT_TOKEN", m.GitToken).
			WithExec([]string{"sh", "-c", "git config --global url.\"https://x-access-token:${GIT_TOKEN}@github.com/\".insteadOf \"https://github.com/\""})
	}

	return ctr.WithExec([]string{"go", "build", "-ldflags", "-s -w", "-o", output, "./cmd/ant"}).File(output), nil
}

// BuildAll builds binaries for all supported platforms
func (m *AnthropicCli) BuildAll(ctx context.Context) (*dagger.Directory, error) {
	targets := []struct {
		os   string
		arch string
	}{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
	}

	binaries := dag.Directory()
	for _, t := range targets {
		binary, err := m.Build(ctx, t.os, t.arch)
		if err != nil {
			return nil, fmt.Errorf("build %s/%s: %w", t.os, t.arch, err)
		}
		binaries = binaries.WithFile(fmt.Sprintf("ant-%s-%s", t.os, t.arch), binary)
	}
	return binaries, nil
}

// GenerateSBOM generates SBOM for built artifacts using Syft
func (m *AnthropicCli) GenerateSBOM(
	ctx context.Context,
	// Directory containing build artifacts
	artifacts *dagger.Directory,
	// SBOM format (cyclonedx-json, spdx-json)
	format string,
) (*dagger.File, error) {
	if format == "" {
		format = "cyclonedx-json"
	}
	return dag.Container().
		From("anchore/syft:latest").
		WithMountedDirectory("/artifacts", artifacts).
		WithWorkdir("/artifacts").
		WithExec([]string{"syft", ".", "-o", format, "--file", "sbom.json"}).
		File("sbom.json"), nil
}

// SignArtifacts signs artifacts using Cosign
func (m *AnthropicCli) SignArtifacts(
	ctx context.Context,
	// Directory containing artifacts to sign
	artifacts *dagger.Directory,
	// Cosign signing key (private key secret)
	key *dagger.Secret,
) (*dagger.Directory, error) {
	return dag.Container().
		From("cgr.dev/chainguard/cosign:latest").
		WithMountedDirectory("/artifacts", artifacts).
		WithMountedSecret("/key/cosign.key", key).
		WithWorkdir("/artifacts").
		WithExec([]string{"sh", "-c", "for f in ant-*; do cosign sign-blob --key=/key/cosign.key --output-signature=${f}.sig ${f}; done"}).
		Directory("/artifacts"), nil
}
