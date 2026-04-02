package main

import (
	"context"
	"fmt"
	"path/filepath"

	"dagger/anthropic-cli/internal/dagger"
)

// DocsGenerateCLIReference generates CLI documentation from source code
// This function parses the Go source files and generates markdown documentation
// for all CLI commands, flags, and usage examples.
func (m *AnthropicCli) DocsGenerateCLIReference(ctx context.Context) (*dagger.Directory, error) {
	// Run the docs generation script
	generated := dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithExec([]string{"mkdir", "-p", "docs/cli/commands"}).
		WithExec([]string{"sh", "scripts/generate-docs.sh"}).
		Directory("/src/docs/cli/commands")

	return generated, nil
}

// DocsBuild builds the static documentation site using MkDocs
func (m *AnthropicCli) DocsBuild(ctx context.Context) (*dagger.Directory, error) {
	container := dag.Container().
		From("python:3.11-alpine").
		WithExec([]string{"pip", "install", "mkdocs", "mkdocs-material"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithExec([]string{"sh", "scripts/generate-docs.sh"}).
		WithExec([]string{"mkdocs", "build", "--site-dir", "site"})

	return container.Directory("/src/site"), nil
}

// DocsBuildFast builds docs without regenerating CLI reference (faster)
func (m *AnthropicCli) DocsBuildFast(ctx context.Context) (*dagger.Directory, error) {
	container := dag.Container().
		From("python:3.11-alpine").
		WithExec([]string{"pip", "install", "mkdocs", "mkdocs-material"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithExec([]string{"mkdocs", "build", "--site-dir", "site"})

	return container.Directory("/src/site"), nil
}

// DocsServe runs the MkDocs development server and returns a service
// This starts a local server accessible at http://localhost:8000
func (m *AnthropicCli) DocsServe() *dagger.Service {
	return dag.Container().
		From("python:3.11-alpine").
		WithExec([]string{"pip", "install", "mkdocs", "mkdocs-material"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithExposedPort(8000).
		WithExec([]string{"mkdocs", "serve", "-a", "0.0.0.0:8000"}).
		AsService()
}

// DocsDeploy builds and deploys documentation to GitHub Pages
// This function handles the full pipeline: build -> deploy
func (m *AnthropicCli) DocsDeploy(
	ctx context.Context,
	// GitHub token for Pages deployment (requires repo scope)
	token *dagger.Secret,
	// Repository in format owner/repo (e.g., "anthropics/anthropic-cli")
	repo string,
	// Branch to deploy from (default: gh-pages)
	// +optional
	deployBranch string,
) error {
	if deployBranch == "" {
		deployBranch = "gh-pages"
	}

	// Build the docs
	docsDir, err := m.DocsBuild(ctx)
	if err != nil {
		return fmt.Errorf("failed to build docs: %w", err)
	}

	// Deploy to GitHub Pages using git
	_, err = dag.Container().
		From("alpine/git:latest").
		WithMountedDirectory("/docs", docsDir).
		WithMountedSecret("/token", token).
		WithEnvVariable("GITHUB_TOKEN", "${token}").
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf(`
				git config --global user.email "docs-bot@anthropic.com" &&
				git config --global user.name "Docs Bot" &&
				git init /tmp/gh-pages &&
				cd /tmp/gh-pages &&
				cp -r /docs/* . &&
				git add -A &&
				git commit -m "Deploy documentation" &&
				git push --force "https://${GITHUB_TOKEN}@github.com/%s.git" HEAD:%s
			`, repo, deployBranch),
		}).
		Sync(ctx)

	if err != nil {
		return fmt.Errorf("failed to deploy docs: %w", err)
	}

	return nil
}

// DocsDeployDryRun builds docs and validates but doesn't deploy
// Use this to verify the build works before actual deployment
func (m *AnthropicCli) DocsDeployDryRun(ctx context.Context) (*dagger.Directory, error) {
	return m.DocsBuild(ctx)
}

// DocsLint validates the documentation structure and markdown files
// Checks for broken links, missing files, and markdown syntax errors
func (m *AnthropicCli) DocsLint(ctx context.Context) (string, error) {
	// Run markdown linting and validation
	output, err := dag.Container().
		From("node:20-alpine").
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithExec([]string{"npm", "install", "-g", "markdownlint-cli"}).
		WithExec([]string{"markdownlint", "docs/**/*.md", "--ignore", "docs/cli/commands/*.md"}).
		Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("docs lint failed: %w", err)
	}

	return output, nil
}

// DocsPackage creates a distributable archive of the documentation
// Useful for offline distribution or archiving
func (m *AnthropicCli) DocsPackage(
	ctx context.Context,
	// Archive format: tar.gz or zip
	// +optional
	format string,
) (*dagger.File, error) {
	if format == "" {
		format = "tar.gz"
	}

	// Build docs
	docsDir, err := m.DocsBuild(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build docs: %w", err)
	}

	// Create archive
	var archiveCmd []string
	var archiveName string

	switch format {
	case "zip":
		archiveName = "docs.zip"
		archiveCmd = []string{"zip", "-r", archiveName, "."}
	default: // tar.gz
		archiveName = "docs.tar.gz"
		archiveCmd = []string{"tar", "-czf", archiveName, "."}
	}

	archivePath := filepath.Join("/src", archiveName)

	file := dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/src/docs", docsDir).
		WithWorkdir("/src/docs").
		WithExec(archiveCmd).
		File(archivePath)

	return file, nil
}
