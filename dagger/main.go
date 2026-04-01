package main

import (
	"context"
	"dagger/anthropic-cli/internal/dagger"
	"dagger/anthropic-cli/dagger/pipelines"
	"dagger/anthropic-cli/dagger/security"
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
	return pipelines.Build(ctx, m.Source, pipelines.BuildConfig{GoOS: goos, GoArch: goarch})
}

func (m *AnthropicCli) BuildAll(ctx context.Context) (*dagger.Directory, error) {
	return pipelines.BuildAll(ctx, m.Source)
}

func (m *AnthropicCli) Attest(ctx context.Context, artifacts *dagger.Directory) (*dagger.File, error) {
	return security.GenerateSLSAProvenance(ctx, artifacts, security.ProvenanceMetadata{
		BuildType:  "https://github.com/anthropics/anthropic-cli/dagger",
		BuilderID:  "https://github.com/anthropics/anthropic-cli/.github/workflows/dagger.yml",
		SourceRepo: "https://github.com/anthropics/anthropic-cli",
	})
}

func (m *AnthropicCli) Sign(ctx context.Context, artifacts *dagger.Directory, key *dagger.Secret) (*dagger.Directory, error) {
	return security.SignArtifacts(ctx, artifacts, key)
}
