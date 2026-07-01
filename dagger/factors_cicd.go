package main

import (
	"context"
	"fmt"

	"dagger/ant-cli/internal/dagger"
)

// GoReleaserFactor runs GoReleaser
type GoReleaserFactor struct {
	config *FactorConfig
	source *dagger.Directory
	mode   string // "snapshot" or "release"
}

func (f *GoReleaserFactor) Name() string           { return "goreleaser" }
func (f *GoReleaserFactor) Dependencies() []string { return []string{"build"} }

func (f *GoReleaserFactor) Execute(ctx context.Context, state *FactorState) (*dagger.Directory, error) {
	args := "release --clean"
	if f.mode == "snapshot" {
		args = "release --snapshot --clean --skip=publish"
	}

	return dag.Container().
		From("goreleaser/goreleaser:v2.15.2").
		WithMountedDirectory("/src", f.source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-cache")).
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-cache")).
		WithMountedCache("/root/.cache/goreleaser", dag.CacheVolume("goreleaser-cache")).
		WithExec([]string{"/goreleaser", args}).
		Directory("/dist").
		Sync(ctx)
}

// CrossPlatformBuildFactor builds for all target platforms in parallel
type CrossPlatformBuildFactor struct {
	config *FactorConfig
	source *dagger.Directory
}

func (f *CrossPlatformBuildFactor) Name() string           { return "cross-platform-build" }
func (f *CrossPlatformBuildFactor) Dependencies() []string { return []string{"cache-warmup"} }

func (f *CrossPlatformBuildFactor) Execute(ctx context.Context, state *FactorState) (*dagger.Directory, error) {
	output := dag.Directory()
	platforms := []struct{ os, arch string }{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
		{"windows", "arm64"},
	}

	type buildResult struct {
		platform string
		binary   *dagger.File
		err      error
	}
	results := make(chan buildResult, len(platforms))

	for _, plat := range platforms {
		go func(p struct{ os, arch string }) {
			binary, err := dag.Container().
				From("golang:1.25-alpine").
				WithExec([]string{"apk", "add", "--no-cache", "git"}).
				WithMountedDirectory("/src", f.source).
				WithWorkdir("/src").
				WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
				WithEnvVariable("CGO_ENABLED", "0").
				WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-cache")).
				WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-cache")).
				WithEnvVariable("GOOS", p.os).
				WithEnvVariable("GOARCH", p.arch).
				WithExec([]string{
					"go", "build",
					"-o", fmt.Sprintf("/output/anthropic-cli-%s-%s", p.os, p.arch),
					"-ldflags=-s -w -X main.Version=$(git describe --tags --always --dirty)",
					"./cmd/ant",
				}).
				File(fmt.Sprintf("/output/anthropic-cli-%s-%s", p.os, p.arch)).
				Sync(ctx)
			results <- buildResult{
				platform: fmt.Sprintf("%s-%s", p.os, p.arch),
				binary:   binary,
				err:      err,
			}
		}(plat)
	}

	for i := 0; i < len(platforms); i++ {
		result := <-results
		if result.err != nil {
			return nil, fmt.Errorf("build for %s failed: %w", result.platform, result.err)
		}
		output = output.WithFile(result.platform, result.binary)
	}

	return output.Sync(ctx)
}

// ReleaseVerificationFactor verifies tag is on main's first-parent history
type ReleaseVerificationFactor struct {
	config *FactorConfig
	source *dagger.Directory
	tag    string
	sha    string
}

func (f *ReleaseVerificationFactor) Name() string           { return "release-verification" }
func (f *ReleaseVerificationFactor) Dependencies() []string { return nil }

func (f *ReleaseVerificationFactor) Execute(ctx context.Context, state *FactorState) (*dagger.Directory, error) {
	return dag.Container().
		From("alpine:3.19").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithMountedDirectory("/src", f.source).
		WithWorkdir("/src").
		WithExec([]string{"git", "fetch", "origin", "main"}).
		WithExec([]string{"sh", "-c", fmt.Sprintf(
			"if ! git rev-list --first-parent origin/main | grep -qxF '%s'; then echo 'FAIL: %s not on main first-parent' > /output/verification.txt; exit 1; fi",
			f.sha, f.sha,
		)}).
		WithNewFile("/output/verification.txt", "PASS: Tag is on main's first-parent history", dagger.ContainerWithNewFileOpts{}).
		Directory("/output").
		Sync(ctx)
}

// PrivateRepoAccessFactor configures access to private Go modules
type PrivateRepoAccessFactor struct {
	config       *FactorConfig
	source       *dagger.Directory
	accessToken  string
	stainlessKey string
}

func (f *PrivateRepoAccessFactor) Name() string           { return "private-repo-access" }
func (f *PrivateRepoAccessFactor) Dependencies() []string { return nil }

func (f *PrivateRepoAccessFactor) Execute(ctx context.Context, state *FactorState) (*dagger.Directory, error) {
	container := dag.Container().
		From("golang:1.25-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithMountedDirectory("/src", f.source).
		WithWorkdir("/src").
		WithEnvVariable("GOPRIVATE", "github.com/anthropics/anthropic-sdk-go,github.com/stainless-sdks/anthropic-go").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-cache")).
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-cache"))

	if f.accessToken != "" {
		container = container.WithExec([]string{"sh", "-c", fmt.Sprintf(
			"git config --global url.\"https://x-access-token:%s@github.com/stainless-sdks/anthropic-go\".insteadOf \"https://github.com/stainless-sdks/anthropic-go\"",
			f.accessToken,
		)})
	}

	return container.
		WithNewFile("/output/git-config.txt", "Private repo access configured", dagger.ContainerWithNewFileOpts{}).
		Directory("/output").
		Sync(ctx)
}
