package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/ant-cli/internal/dagger"
)

// ImageCatalogEntry represents a container image used by factors
type ImageCatalogEntry struct {
	Name    string
	Version string
	Source  string
}

// ImageCatalog generates a catalog of all container images used by factors
func (m *AntCli) ImageCatalog(ctx context.Context) (*dagger.File, error) {
	images := []ImageCatalogEntry{
		{Name: "golang", Version: "1.23-alpine", Source: "BuildFactor, TestFactor, LintFactor, StaticAnalysisFactor, VulnScanFactor, SBOMFactor, SLSAProvenanceFactor, LicenseCheckFactor, PrivateRepoAccessFactor, CacheWarmupFactor"},
		{Name: "zricethezav/gitleaks", Version: "latest", Source: "SecretScanningFactor"},
		{Name: "instrumenta/conftest", Version: "latest", Source: "PolicyCheckFactor"},
		{Name: "goreleaser/goreleaser", Version: "v2.15.2", Source: "GoReleaserFactor"},
		{Name: "alpine", Version: "3.19", Source: "ReleaseVerificationFactor"},
		{Name: "ghcr.io/sigstore/cosign/cosign", Version: "v2.2.3", Source: "CosignSign"},
	}

	var catalog strings.Builder
	catalog.WriteString("# Dagger Container Image Catalog\n\n")
	catalog.WriteString("## Images\n\n")

	for _, img := range images {
		catalog.WriteString(fmt.Sprintf("- **%s:%s** - Used by: %s\n", img.Name, img.Version, img.Source))
	}

	catalog.WriteString("\n## Version Policy\n\n")
	catalog.WriteString("Versions must be pinned. Current exceptions to fix:\n")
	catalog.WriteString("- zricethezav/gitleaks:latest\n")
	catalog.WriteString("- instrumenta/conftest:latest\n")

	return dag.Directory().
		WithNewFile("/catalog/images.md", catalog.String(), dagger.DirectoryWithNewFileOpts{}).
		File("/catalog/images.md").
		Sync(ctx)
}
