package main

import (
	"context"

	"dagger/ant-cli/internal/dagger"
)

// StaticAnalysisFactor runs gosec
type StaticAnalysisFactor struct {
	config *FactorConfig
	source *dagger.Directory
}

func (f *StaticAnalysisFactor) Name() string           { return "static-analysis" }
func (f *StaticAnalysisFactor) Dependencies() []string { return []string{"cache-warmup"} }

func (f *StaticAnalysisFactor) Execute(ctx context.Context, state *FactorState) (*dagger.Directory, error) {
	return dag.Container().
		From("golang:1.25-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithMountedDirectory("/src", f.source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-cache")).
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-cache")).
		WithMountedCache("/root/.cache/gosec", dag.CacheVolume("gosec-cache")).
		WithExec([]string{"go", "install", "github.com/securego/gosec/v2/cmd/gosec@latest"}).
		WithExec([]string{"mkdir", "-p", "/output"}).
		WithExec([]string{"sh", "-c", "gosec -fmt=json -out=/output/gosec-report.json -stdout=false ./...; exit 0"}).
		Directory("/output").
		Sync(ctx)
}

// SecretScanningFactor runs gitleaks
type SecretScanningFactor struct {
	config *FactorConfig
	source *dagger.Directory
}

func (f *SecretScanningFactor) Name() string           { return "secret-scanning" }
func (f *SecretScanningFactor) Dependencies() []string { return nil }

func (f *SecretScanningFactor) Execute(ctx context.Context, state *FactorState) (*dagger.Directory, error) {
	return dag.Container().
		From("zricethezav/gitleaks:latest").
		WithMountedDirectory("/src", f.source).
		WithWorkdir("/src").
		WithMountedCache("/root/.cache/gitleaks", dag.CacheVolume("gitleaks-cache")).
		WithExec([]string{"mkdir", "-p", "/output"}).
		WithExec([]string{"gitleaks", "detect", "--source", ".", "--report-path", "/output/gitleaks-report.json", "--report-format", "json", "--exit-code", "0"}).
		Directory("/output").
		Sync(ctx)
}

// VulnScanFactor runs govulncheck
type VulnScanFactor struct {
	config *FactorConfig
	source *dagger.Directory
}

func (f *VulnScanFactor) Name() string           { return "vuln-scan" }
func (f *VulnScanFactor) Dependencies() []string { return []string{"cache-warmup"} }

func (f *VulnScanFactor) Execute(ctx context.Context, state *FactorState) (*dagger.Directory, error) {
	return dag.Container().
		From("golang:1.25-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithMountedDirectory("/src", f.source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-cache")).
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-cache")).
		WithMountedCache("/root/.cache/govulncheck", dag.CacheVolume("govulncheck-cache")).
		WithExec([]string{"go", "install", "golang.org/x/vuln/cmd/govulncheck@latest"}).
		WithExec([]string{"mkdir", "-p", "/output"}).
		WithExec([]string{"sh", "-c", "govulncheck -format=json ./... > /output/vulns.json 2>&1; exit 0"}).
		Directory("/output").
		Sync(ctx)
}
