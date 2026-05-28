package main

import (
	"context"

	"dagger/ant-cli/internal/dagger"
)

// CacheWarmupFactor pre-warms all caches before pipeline execution
// +cache="session"
type CacheWarmupFactor struct {
	config *FactorConfig
	source *dagger.Directory
}

func (f *CacheWarmupFactor) Name() string           { return "cache-warmup" }
func (f *CacheWarmupFactor) Dependencies() []string { return nil }

func (f *CacheWarmupFactor) Execute(ctx context.Context, state *FactorState) (*dagger.Directory, error) {
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
		WithExec([]string{"go", "mod", "download"}).
		WithExec([]string{"go", "build", "-o", "/dev/null", "./..."}).
		Directory("/src").
		Sync(ctx)
}
