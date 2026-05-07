package main

import (
	"context"

	"dagger/ant-cli/internal/dagger"
)

// BuildFactor compiles the project
type BuildFactor struct {
	config *FactorConfig
	source *dagger.Directory
}

func (f *BuildFactor) Name() string           { return "build" }
func (f *BuildFactor) Dependencies() []string { return []string{"cache-warmup"} }

func (f *BuildFactor) Execute(ctx context.Context, state *FactorState) (*dagger.Directory, error) {
	return dag.Container().
		From("golang:1.25-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithMountedDirectory("/src", f.source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithEnvVariable("CGO_ENABLED", "0").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-cache")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-cache")).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithExec([]string{"go", "build", "-a", "-o", "/output/", "-ldflags=-s -w", "./..."}).
		Directory("/output").
		Sync(ctx)
}

// TestFactor runs the test suite
type TestFactor struct {
	config *FactorConfig
	source *dagger.Directory
}

func (f *TestFactor) Name() string           { return "test" }
func (f *TestFactor) Dependencies() []string { return []string{"cache-warmup"} }

func (f *TestFactor) Execute(ctx context.Context, state *FactorState) (*dagger.Directory, error) {
	return dag.Container().
		From("golang:1.25-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git", "lsof"}).
		WithMountedDirectory("/src", f.source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-cache")).
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-cache")).
		WithExec([]string{"mkdir", "-p", "/output"}).
		WithExec([]string{"sh", "-c", "go test -v -count=1 -p=4 -coverprofile=/output/coverage.out ./... > /output/test.log 2>&1; echo $? > /output/exit-code.txt; exit 0"}).
		WithExec([]string{"sh", "-c", "go tool cover -html=/output/coverage.out -o /output/coverage.html 2>/dev/null || true"}).
		Directory("/output").
		Sync(ctx)
}

// LintFactor runs golangci-lint
type LintFactor struct {
	config *FactorConfig
	source *dagger.Directory
}

func (f *LintFactor) Name() string           { return "lint" }
func (f *LintFactor) Dependencies() []string { return []string{"cache-warmup"} }

func (f *LintFactor) Execute(ctx context.Context, state *FactorState) (*dagger.Directory, error) {
	return dag.Container().
		From("golang:1.25-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithMountedDirectory("/src", f.source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-cache")).
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-cache")).
		WithMountedCache("/root/.cache/golangci-lint", dag.CacheVolume("golangci-lint-cache")).
		WithExec([]string{"go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"}).
		WithExec([]string{"mkdir", "-p", "/output"}).
		WithExec([]string{"golangci-lint", "run", "./...", "--out-format=json", "--issues-exit-code=0"}).
		WithNewFile("/output/lint-report.json", "", dagger.ContainerWithNewFileOpts{}).
		Directory("/output").
		Sync(ctx)
}
