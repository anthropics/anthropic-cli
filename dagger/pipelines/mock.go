package pipelines

import (
	"dagger/anthropic-cli/internal/dagger"
)

// MockServer returns a container running the Steady mock server
func MockServer(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("node:20-alpine").
		WithExec([]string{"npm", "install", "-g", "@stainless/cli"}).
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExposedPort(4010).
		WithExec([]string{"steady", "mock", "--port", "4010"})
}
