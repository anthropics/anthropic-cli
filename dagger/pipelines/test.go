package pipelines

import (
	"context"
	"dagger/anthropic-cli/internal/dagger"

	"dagger.io/dagger/dag"
)

// Test runs the full test suite including mock server
func Test(ctx context.Context, source *dagger.Directory) (string, error) {
	// Start mock server
	mockServer := dag.Container().
		From("node:20-alpine").
		WithExec([]string{"npm", "install", "-g", "@stainless/cli"}).
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExposedPort(4010).
		WithExec([]string{"steady", "mock", "--port", "4010"})

	// Run tests with mock server
	return dag.Container().
		From("golang:1.24-alpine").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithServiceBinding("steady", mockServer).
		WithEnvVariable("TEST_API_BASE_URL", "http://steady:4010").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithExec([]string{"go", "test", "-v", "./..."}).
		Stdout(ctx)
}

// TestFast runs unit tests only (no mock server)
func TestFast(ctx context.Context, source *dagger.Directory) (string, error) {
	return dag.Container().
		From("golang:1.24-alpine").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithExec([]string{"go", "test", "-short", "./..."}).
		Stdout(ctx)
}
