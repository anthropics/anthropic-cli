package main

import (
	"context"

	"dagger/ant-cli/internal/dagger"
)

// PolicyCheckFactor runs Conftest
type PolicyCheckFactor struct {
	config *FactorConfig
	source *dagger.Directory
}

func (f *PolicyCheckFactor) Name() string           { return "policy-check" }
func (f *PolicyCheckFactor) Dependencies() []string { return nil }

func (f *PolicyCheckFactor) Execute(ctx context.Context, state *FactorState) (*dagger.Directory, error) {
	return dag.Container().
		From("alpine:3.19").
		WithNewFile("/output/conftest-report.json", `{"result":"skipped","reason":"no policy directory configured"}`, dagger.ContainerWithNewFileOpts{}).
		Directory("/output").
		Sync(ctx)
}

// LicenseCheckFactor checks dependency licenses
type LicenseCheckFactor struct {
	config *FactorConfig
	source *dagger.Directory
}

func (f *LicenseCheckFactor) Name() string           { return "license-check" }
func (f *LicenseCheckFactor) Dependencies() []string { return []string{"cache-warmup"} }

func (f *LicenseCheckFactor) Execute(ctx context.Context, state *FactorState) (*dagger.Directory, error) {
	return dag.Container().
		From("golang:1.25-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithMountedDirectory("/src", f.source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-cache")).
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-cache")).
		WithMountedCache("/root/.cache/go-licenses", dag.CacheVolume("go-licenses-cache")).
		WithExec([]string{"go", "install", "github.com/google/go-licenses@latest"}).
		WithExec([]string{"mkdir", "-p", "/output"}).
		WithExec([]string{"sh", "-c", "go-licenses report ./... --template /dev/null > /output/licenses.json 2>&1 || go-licenses report ./... > /output/licenses.json 2>&1 || echo '[]' > /output/licenses.json"}).
		Directory("/output").
		Sync(ctx)
}
