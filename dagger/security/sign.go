package security

import (
	"context"
	"dagger/anthropic-cli/internal/dagger"
)

// SignArtifacts signs all artifacts in a directory using Cosign
func SignArtifacts(
	ctx context.Context,
	artifacts *dagger.Directory,
	key *dagger.Secret,
) (*dagger.Directory, error) {
	signed := dag.Container().
		From("cgr.dev/chainguard/cosign:latest").
		WithMountedDirectory("/artifacts", artifacts).
		WithMountedSecret("/key/cosign.key", key).
		WithWorkdir("/artifacts")

	// Sign all files in the artifacts directory
	return signed.
		WithExec([]string{"sh", "-c", "for f in *; do cosign sign-blob --key=/key/cosign.key --output-signature=$f.sig $f; done"}).
		Directory("/artifacts"), nil
}

// VerifyArtifacts verifies signatures using Cosign
func VerifyArtifacts(
	ctx context.Context,
	artifacts *dagger.Directory,
	key *dagger.Secret,
) error {
	verify := dag.Container().
		From("cgr.dev/chainguard/cosign:latest").
		WithMountedDirectory("/artifacts", artifacts).
		WithMountedSecret("/key/cosign.pub", key).
		WithWorkdir("/artifacts").
		WithExec([]string{"sh", "-c", "for f in *.sig; do cosign verify-blob --key=/key/cosign.pub --signature=$f ${f%.sig}; done"})

	_, err := verify.Stdout(ctx)
	return err
}
