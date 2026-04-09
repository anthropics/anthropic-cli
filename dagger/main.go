package main

import (
	"context"
	"dagger/anthropic-cli/internal/dagger"
)

// AnthropicCli provides Dagger functions for local development
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

// Build runs Go build in a clean container
func (m *AnthropicCli) Build(ctx context.Context) (string, error) {
	return dag.Container().
		From("golang:1.26-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithExec([]string{"go", "build", "./..."}).
		Stdout(ctx)
}

// Test runs the test suite in a clean container
func (m *AnthropicCli) Test(ctx context.Context) (string, error) {
	return dag.Container().
		From("golang:1.26-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git", "curl", "lsof"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithExec([]string{"go", "test", "./..."}).
		Stdout(ctx)
}

// Lint runs the lint script in a clean container
func (m *AnthropicCli) Lint(ctx context.Context) (string, error) {
	return dag.Container().
		From("golang:1.26-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git", "bash"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithExec([]string{"./scripts/lint"}).
		Stdout(ctx)
}
