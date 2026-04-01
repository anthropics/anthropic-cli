package pipelines

import (
	"context"
	"dagger/anthropic-cli/internal/dagger"
	"fmt"
)

// BuildConfig contains build parameters
type BuildConfig struct {
	GoOS   string
	GoArch string
}

// Build creates a cross-compiled binary
func Build(ctx context.Context, source *dagger.Directory, cfg BuildConfig) (*dagger.File, error) {
	output := fmt.Sprintf("bin/ant-%s-%s", cfg.GoOS, cfg.GoArch)
	return dag.Container().
		From("golang:1.24-alpine").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithEnvVariable("GOOS", cfg.GoOS).
		WithEnvVariable("GOARCH", cfg.GoArch).
		WithEnvVariable("CGO_ENABLED", "0").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithExec([]string{"go", "build", "-ldflags", "-s -w", "-o", output, "./cmd/ant"}).
		File(output), nil
}

// BuildAll builds for all target platforms
func BuildAll(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	targets := []BuildConfig{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
	}

	binaries := dag.Directory()
	for _, cfg := range targets {
		binary, err := Build(ctx, source, cfg)
		if err != nil {
			return nil, fmt.Errorf("build %s/%s: %w", cfg.GoOS, cfg.GoArch, err)
		}
		binaries = binaries.WithFile(fmt.Sprintf("ant-%s-%s", cfg.GoOS, cfg.GoArch), binary)
	}
	return binaries, nil
}
