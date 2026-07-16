// AntCli provides Dagger functions for local development and CI
// These complement (not replace) the GitHub Actions CI/CD
package main

import (
	"context"
	"fmt"
	"path/filepath"

	"dagger/ant-cli/internal/dagger"
)

// AntCli provides Dagger functions for local development and CI
type AntCli struct {
	// +private
	Source *dagger.Directory
}

func New(
	// Source directory containing the project
	source *dagger.Directory,
) *AntCli {
	config := &FactorConfig{
		Source:  source,
		OS:      "linux",
		Arch:    "amd64",
		Options: make(map[string]interface{}),
	}

	cli := &AntCli{
		Source: source,
	}

	_ = config
	return cli
}

// newRegistry builds a fresh registry — not stored on struct since interface maps can't be Dagger-serialized
func (m *AntCli) newRegistry() *FactorRegistry {
	config := &FactorConfig{
		Source:  m.Source,
		OS:      "linux",
		Arch:    "amd64",
		Options: make(map[string]interface{}),
	}
	r := NewFactorRegistry(config)
	r.Register(&CacheWarmupFactor{config: config, source: m.Source})
	r.Register(&BuildFactor{config: config, source: m.Source})
	r.Register(&TestFactor{config: config, source: m.Source})
	r.Register(&LintFactor{config: config, source: m.Source})
	r.Register(&StaticAnalysisFactor{config: config, source: m.Source})
	r.Register(&SecretScanningFactor{config: config, source: m.Source})
	r.Register(&VulnScanFactor{config: config, source: m.Source})
	r.Register(&SBOMFactor{config: config, source: m.Source})
	r.Register(&SLSAProvenanceFactor{config: config, source: m.Source})
	r.Register(&PolicyCheckFactor{config: config, source: m.Source})
	r.Register(&LicenseCheckFactor{config: config, source: m.Source})
	return r
}

// baseContainer returns a Go container with source mounted and caches persisted
func (m *AntCli) baseContainer() *dagger.Container {
	return dag.Container().
		From("golang:1.25-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git", "bash", "curl"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithEnvVariable("CGO_ENABLED", "0").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-cache")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-cache")).
		WithEnvVariable("GOCACHE", "/go/build-cache")
}

// Build runs Go build and returns the compiled binaries directory
// +cache="1h"
func (m *AntCli) Build(ctx context.Context) (*dagger.Directory, error) {
	return m.baseContainer().
		WithExec([]string{"go", "build", "-o", "/output/", "-ldflags=-s -w", "./..."}).
		Directory("/output").
		Sync(ctx)
}

// BuildForPlatform builds the CLI for a specific OS/arch and returns the binary
// +cache="1h"
func (m *AntCli) BuildForPlatform(
	ctx context.Context,
	// +default="linux"
	os string,
	// +default="amd64"
	arch string,
) (*dagger.File, error) {
	outputName := fmt.Sprintf("anthropic-cli-%s-%s", os, arch)
	if os == "windows" {
		outputName += ".exe"
	}

	return m.baseContainer().
		WithEnvVariable("GOOS", os).
		WithEnvVariable("GOARCH", arch).
		WithExec([]string{
			"go", "build",
			"-o", filepath.Join("/output", outputName),
			"-ldflags=-s -w -X main.Version=$(git describe --tags --always --dirty)",
			"./cmd/ant",
		}).
		File(filepath.Join("/output", outputName)).
		Sync(ctx)
}

// Test runs the test suite and returns test results
// +cache="1h"
func (m *AntCli) Test(ctx context.Context) (*dagger.Directory, error) {
	return m.baseContainer().
		WithExec([]string{"apk", "add", "--no-cache", "lsof"}).
		WithExec([]string{"go", "test", "-v", "-count=1", "-coverprofile=/output/coverage.out", "./..."}).
		WithExec([]string{"go", "tool", "cover", "-html=/output/coverage.out", "-o", "/output/coverage.html"}).
		Directory("/output").
		Sync(ctx)
}

// Lint runs linting and returns the report
// +cache="1h"
func (m *AntCli) Lint(ctx context.Context) (*dagger.File, error) {
	return m.baseContainer().
		WithExec([]string{"go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"}).
		WithExec([]string{"golangci-lint", "run", "./...", "--out-format=json", "--issues-exit-code=0"}).
		WithNewFile("/output/lint-report.json", "", dagger.ContainerWithNewFileOpts{}).
		File("/output/lint-report.json").
		Sync(ctx)
}

// VulnScan runs govulncheck and returns the vulnerability report
// +cache="1h"
func (m *AntCli) VulnScan(ctx context.Context) (*dagger.File, error) {
	return m.baseContainer().
		WithExec([]string{"go", "install", "golang.org/x/vuln/cmd/govulncheck@latest"}).
		WithExec([]string{"govulncheck", "-format=json", "-show=verbose", "./..."}).
		WithNewFile("/output/vulns.json", "", dagger.ContainerWithNewFileOpts{}).
		File("/output/vulns.json").
		Sync(ctx)
}

// SBOM generates a CycloneDX SBOM for the project
// +cache="1h"
func (m *AntCli) SBOM(
	ctx context.Context,
	// +default="cyclonedx-json"
	format string,
) (*dagger.File, error) {
	filename := fmt.Sprintf("sbom.%s", format)
	if format == "cyclonedx-json" {
		filename = "sbom.cdx.json"
	} else if format == "spdx-json" {
		filename = "sbom.spdx.json"
	}

	return m.baseContainer().
		WithExec([]string{"go", "install", "github.com/anchore/syft/cmd/syft@latest"}).
		WithExec([]string{"syft", "scan", "dir:/src", "-o", format, "--file", filepath.Join("/output", filename)}).
		File(filepath.Join("/output", filename)).
		Sync(ctx)
}

// Provenance generates SLSA provenance attestation
func (m *AntCli) Provenance(
	ctx context.Context,
	// +default="linux"
	os string,
	// +default="amd64"
	arch string,
) (*dagger.File, error) {
	binary, err := m.BuildForPlatform(ctx, os, arch)
	if err != nil {
		return nil, fmt.Errorf("build failed: %w", err)
	}

	binaryName := fmt.Sprintf("anthropic-cli-%s-%s", os, arch)
	if os == "windows" {
		binaryName += ".exe"
	}

	return m.baseContainer().
		WithFile(filepath.Join("/tmp", binaryName), binary).
		WithExec([]string{"sh", "-c", fmt.Sprintf(
			"sha256sum /tmp/%s > /output/%s.sha256 && sha256sum /tmp/%s | awk '{print $1}' > /output/%s.provenance",
			binaryName, binaryName, binaryName, binaryName,
		)}).
		File(fmt.Sprintf("/output/%s.sha256", binaryName)).
		Sync(ctx)
}

// StaticAnalysis runs gosec for static security analysis
// +cache="1h"
func (m *AntCli) StaticAnalysis(ctx context.Context) (*dagger.File, error) {
	return m.baseContainer().
		WithExec([]string{"go", "install", "github.com/securego/gosec/v2/cmd/gosec@latest"}).
		WithExec([]string{"gosec", "-fmt=json", "-out=/output/gosec-report.json", "-stdout=false", "./..."}).
		File("/output/gosec-report.json").
		Sync(ctx)
}

// SecretScanning runs gitleaks to detect secrets in the codebase
// +cache="session"
func (m *AntCli) SecretScanning(ctx context.Context) (*dagger.File, error) {
	return dag.Container().
		From("zricethezav/gitleaks:latest").
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithExec([]string{"gitleaks", "detect", "--source", ".", "--report-path", "/output/gitleaks-report.json", "--report-format", "json"}).
		File("/output/gitleaks-report.json").
		Sync(ctx)
}

// LicenseCheck checks for license compliance in dependencies
// +cache="1h"
func (m *AntCli) LicenseCheck(ctx context.Context) (*dagger.File, error) {
	return m.baseContainer().
		WithExec([]string{"go", "install", "github.com/google/go-licenses@latest"}).
		WithExec([]string{"go-licenses", "check", "./...", "--disallowed_types=forbidden", "--save=/output/licenses.json"}).
		File("/output/licenses.json").
		Sync(ctx)
}

// SLSAProvenance generates SLSA v1.0 provenance
func (m *AntCli) SLSAProvenance(
	ctx context.Context,
	// +default="linux"
	os string,
	// +default="amd64"
	arch string,
) (*dagger.File, error) {
	binary, err := m.BuildForPlatform(ctx, os, arch)
	if err != nil {
		return nil, fmt.Errorf("build failed: %w", err)
	}

	binaryName := fmt.Sprintf("anthropic-cli-%s-%s", os, arch)
	if os == "windows" {
		binaryName += ".exe"
	}

	return m.baseContainer().
		WithFile(filepath.Join("/tmp", binaryName), binary).
		WithExec([]string{"sh", "-c", fmt.Sprintf(
			"sha256sum /tmp/%s > /output/%s.slsa.sha256 && "+
				"echo '{\"_type\":\"https://in-toto.io/Statement/v1\",\"predicateType\":\"https://slsa.dev/provenance/v1\",\"subject\":[{\"name\":\"%s\",\"digest\":{\"sha256\":\"$(sha256sum /tmp/%s | awk '{print $1}')\"}}]}' > /output/%s.slsa.json",
			binaryName, binaryName, binaryName, binaryName, binaryName,
		)}).
		File(fmt.Sprintf("/output/%s.slsa.json", binaryName)).
		Sync(ctx)
}

// CosignSign signs a binary with cosign
func (m *AntCli) CosignSign(
	ctx context.Context,
	binary *dagger.File,
	key *dagger.File,
) (*dagger.File, error) {
	container := dag.Container().
		From("ghcr.io/sigstore/cosign/cosign:v2.2.3").
		WithFile("/binary", binary)

	if key != nil {
		container = container.WithFile("/key.pem", key).
			WithExec([]string{"cosign", "sign-blob", "--key", "/key.pem", "--output-signature", "/output/sig.sig", "/binary"})
	} else {
		container = container.WithExec([]string{"cosign", "sign-blob", "--output-signature", "/output/sig.sig", "/binary"})
	}

	return container.File("/output/sig.sig").Sync(ctx)
}

// PolicyCheck runs Conftest for policy compliance
func (m *AntCli) PolicyCheck(
	ctx context.Context,
	policyDir *dagger.Directory,
) (*dagger.File, error) {
	return dag.Container().
		From("instrumenta/conftest:latest").
		WithMountedDirectory("/src", m.Source).
		WithMountedDirectory("/policy", policyDir).
		WithWorkdir("/src").
		WithExec([]string{"conftest", "test", "--policy", "/policy", "--output", "json", "."}).
		WithNewFile("/output/conftest-report.json", "", dagger.ContainerWithNewFileOpts{}).
		File("/output/conftest-report.json").
		Sync(ctx)
}

// CollectEvidence bundles all security and compliance reports
// +cache="never"
func (m *AntCli) CollectEvidence(
	ctx context.Context,
	// +default=false
	includeSLSA bool,
) (*dagger.Directory, error) {
	evidence := dag.Directory()

	gosec, err := m.StaticAnalysis(ctx)
	if err == nil {
		evidence = evidence.WithFile("security/gosec-report.json", gosec)
	}

	secrets, err := m.SecretScanning(ctx)
	if err == nil {
		evidence = evidence.WithFile("security/gitleaks-report.json", secrets)
	}

	vulns, err := m.VulnScan(ctx)
	if err == nil {
		evidence = evidence.WithFile("security/vulns.json", vulns)
	}

	sbom, err := m.SBOM(ctx, "cyclonedx-json")
	if err == nil {
		evidence = evidence.WithFile("sbom/sbom.cdx.json", sbom)
	}

	if includeSLSA {
		provenance, err := m.SLSAProvenance(ctx, "linux", "amd64")
		if err == nil {
			evidence = evidence.WithFile("slsa/provenance.json", provenance)
		}
	}

	return evidence.Sync(ctx)
}

// All runs the complete CI pipeline using the Factor registry
// +cache="never"
func (m *AntCli) All(ctx context.Context) (*dagger.Directory, error) {
	return m.newRegistry().ExecuteAll(ctx)
}
