package main

import (
	"context"
	"dagger/anthropic-cli/internal/dagger"
	"dagger/anthropic-cli/dagger/pipelines"
)

type AnthropicCli struct {
	Source *dagger.Directory
}

func New(source *dagger.Directory) *AnthropicCli {
	return &AnthropicCli{Source: source}
}

func (m *AnthropicCli) Lint(ctx context.Context) (string, error) {
	return pipelines.Lint(ctx, m.Source)
}

func (m *AnthropicCli) Test(ctx context.Context) (string, error) {
	return pipelines.Test(ctx, m.Source)
}

func (m *AnthropicCli) TestFast(ctx context.Context) (string, error) {
	return pipelines.TestFast(ctx, m.Source)
}

func (m *AnthropicCli) Build(ctx context.Context, goos string, goarch string) (*dagger.File, error) {
	return dag.Container().
		From("golang:1.24-alpine").
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithEnvVariable("GOOS", goos).
		WithEnvVariable("GOARCH", goarch).
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithExec([]string{"go", "build", "-o", "bin/ant", "./cmd/ant"}).
		File("bin/ant"), nil
}
