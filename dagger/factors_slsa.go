package main

import (
	"context"
	"fmt"

	"dagger/ant-cli/internal/dagger"
)

// SBOMFactor generates CycloneDX SBOM
type SBOMFactor struct {
	config *FactorConfig
	source *dagger.Directory
}

func (f *SBOMFactor) Name() string           { return "sbom" }
func (f *SBOMFactor) Dependencies() []string { return []string{"build"} }

func (f *SBOMFactor) Execute(ctx context.Context, state *FactorState) (*dagger.Directory, error) {
	return dag.Container().
		From("golang:1.25-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithMountedDirectory("/src", f.source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-cache")).
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-cache")).
		WithMountedCache("/root/.cache/syft", dag.CacheVolume("syft-cache")).
		WithExec([]string{"go", "install", "github.com/anchore/syft/cmd/syft@latest"}).
		WithExec([]string{"mkdir", "-p", "/output"}).
		WithExec([]string{"syft", "scan", "dir:/src", "-o", "cyclonedx-json", "--file", "/output/sbom.cdx.json"}).
		Directory("/output").
		Sync(ctx)
}

// SLSAProvenanceFactor generates SLSA provenance
type SLSAProvenanceFactor struct {
	config *FactorConfig
	source *dagger.Directory
}

func (f *SLSAProvenanceFactor) Name() string           { return "slsa-provenance" }
func (f *SLSAProvenanceFactor) Dependencies() []string { return []string{"build"} }

func (f *SLSAProvenanceFactor) Execute(ctx context.Context, state *FactorState) (*dagger.Directory, error) {
	buildArtifacts := state.Artifacts["build"]
	if buildArtifacts == nil {
		return nil, fmt.Errorf("build artifacts not found in state")
	}

	return dag.Container().
		From("golang:1.25-alpine").
		WithMountedDirectory("/build", buildArtifacts).
		WithExec([]string{"mkdir", "-p", "/output"}).
		WithExec([]string{"sh", "-c", "sha256sum /build/* > /output/provenance.sha256 && echo '{\"_type\":\"https://in-toto.io/Statement/v1\",\"predicateType\":\"https://slsa.dev/provenance/v1\"}' > /output/slsa.json"}).
		Directory("/output").
		Sync(ctx)
}
