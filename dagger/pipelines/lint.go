// Package pipelines provides CI/CD pipeline functions for Anthropic CLI
package pipelines

import (
	"context"
	"dagger/anthropic-cli/internal/dagger"
)

// Lint runs the Go linter (go build ./...)
func Lint(ctx context.Context, source *dagger.Directory) (string, error) {
	return dag.Container().
		From("golang:1.24-alpine").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithExec([]string{"go", "build", "./..."}).
		Stdout(ctx)
}
