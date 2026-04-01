package security

import (
	"context"
	"dagger/anthropic-cli/internal/dagger"
	"encoding/json"
	"fmt"
	"time"
)

// ProvenanceMetadata contains build provenance information
type ProvenanceMetadata struct {
	BuildType   string
	BuilderID   string
	SourceRepo  string
	CommitSHA   string
	BuildTime   time.Time
}

// GenerateSLSAProvenance creates SLSA provenance for artifacts
func GenerateSLSAProvenance(
	ctx context.Context,
	artifacts *dagger.Directory,
	meta ProvenanceMetadata,
) (*dagger.File, error) {
	provenance := map[string]interface{}{
		"_type": "https://in-toto.io/Statement/v0.1",
		"subject": []map[string]string{
			{"name": "ant", "digest": {"sha256": "placeholder"}},
		},
		"predicateType": "https://slsa.dev/provenance/v0.2",
		"predicate": map[string]interface{}{
			"buildType":   meta.BuildType,
			"builder":     map[string]string{"id": meta.BuilderID},
			"invocation":  map[string]interface{}{},
			"metadata": map[string]interface{}{
				"buildInvocationId": fmt.Sprintf("%d", time.Now().Unix()),
				"buildFinishedOn":   meta.BuildTime.Format(time.RFC3339),
			},
			"materials": []map[string]string{
				{"uri": meta.SourceRepo, "digest": {"sha1": meta.CommitSHA}},
			},
		},
	}

	data, err := json.MarshalIndent(provenance, "", "  ")
	if err != nil {
		return nil, err
	}

	return dag.Container().
		From("alpine:latest").
		WithNewFile("/provenance.intoto.jsonl", string(data)).
		File("/provenance.intoto.jsonl"), nil
}
