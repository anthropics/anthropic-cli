// Package claude teaches the core reconciler about the Claude Developer
// Platform: which resources exist, what their files look like, which fields
// the server owns, and how to talk to the API.
//
// Everything Claude-specific in `ant apply` lives here. The core package holds
// the reconciler itself and does not import this one, so a hand-written fact
// about the API has exactly one home — and when the spec moves, this is the
// only directory to re-read.
package claude

import (
	"maps"

	"github.com/anthropics/anthropic-cli/v2/internal/apischema"
	"github.com/anthropics/anthropic-cli/v2/internal/declarative/core"
)

// The kinds `ant apply` manages. Each is the `type` of its API resource, by
// which fromAPI finds its rows.
const (
	KindSkill       core.Kind = "skill"
	KindEnvironment core.Kind = "environment"
	KindMemoryStore core.Kind = "memory_store"
	KindVault       core.Kind = "vault"
	KindAgent       core.Kind = "agent"
	KindDeployment  core.Kind = "deployment"
)

// Removal is not uniform: a skill can only be deleted, and every other kind is
// archived. Where the API offers both we archive, because an archived resource
// keeps its data and a deleted one does not.
var (
	byArchiving = core.Destroy{Verb: "archive", Past: "archived"}
	byDeleting  = core.Destroy{Verb: "delete", Past: "deleted"}
)

// Field sets shared across kinds, so each kind's own Fields list only what is
// its own.
var (
	// named kinds show their name in a plan.
	named = core.Fields{
		"name": {Summary: true},
	}
	// archivable kinds report when they were archived; there is no unarchive.
	archivable = core.Fields{
		"archived_at": {Archived: true},
	}
	// metadata is the user key-value bag, patched key by key.
	metadata = core.Fields{
		"metadata": {Metadata: true},
	}
)

// updatedAt is the stamp of every kind whose updated_at moves when it is edited.
func updatedAt(remote map[string]any) string {
	at, _ := remote["updated_at"].(string)
	return at
}

// agentVersion is an agent's stamp. The API mints a new version on every edit
// and none for a request that changes nothing, and lockfiles written before
// stamps were recorded already hold it.
func agentVersion(remote map[string]any) string {
	return core.FormatVersion(remote["version"])
}

// skillVersion is a skill's stamp: the API mints a new version per upload.
func skillVersion(remote map[string]any) string {
	return core.FormatVersion(remote["latest_version_id"])
}

// with merges field sets into one; a later set wins on a shared path.
func with(sets ...core.Fields) core.Fields {
	out := core.Fields{}
	for _, set := range sets {
		maps.Copy(out, set)
	}
	return out
}

// fromAPI overwrites Computed, WriteOnly and Ref on every path apischema has
// for the kind.
func fromAPI(kind core.Kind, byHand core.Fields) core.Fields {
	for _, resource := range apischema.Resources {
		if resource.Type != string(kind) {
			continue
		}
		for _, f := range resource.Fields {
			field := byHand[f.Path]
			field.Computed = f.InResponse && !f.InCreate
			field.WriteOnly = f.InCreate && !f.InResponse
			field.Ref = reference(f)
			byHand[f.Path] = field
		}
	}
	return byHand
}

func reference(f apischema.Field) *core.Ref {
	for _, r := range f.References {
		resource, ok := apischema.Resources[r.Resource]
		// A path names a skill of the user's own, never one of Anthropic's.
		if !ok || r.Type == "anthropic" {
			continue
		}
		as := core.EncodeID
		if r.IDField != "" {
			as = core.EncodeObject(r.Type, r.IDField)
		}
		return &core.Ref{To: []core.Kind{core.Kind(resource.Type)}, List: f.List, As: as}
	}
	return nil
}

// Registry describes every kind to the reconciler.
//
// What each kind passes to fromAPI is written by hand, and nothing checks it
// against the API: a new clearable field not marked Clearable means deleting
// the line from a file silently does nothing.
//
// Order matters: a kind must come after everything it can reference.
// TestKindOrderIsTopological holds that.
func Registry() *core.Registry {
	return core.NewRegistry(schema{},
		core.KindSpec{
			Kind:     KindSkill,
			IDPrefix: "skill",
			// A skill's content is its files; the server derives name and
			// description from SKILL.md and mints a version per upload.
			VersionField:         "latest_version_id",
			Stamp:                skillVersion,
			StampMovesOnlyOnEdit: true,
			StampName:            "latest_version_id",
			Destroy:              byDeleting,
			Fields: fromAPI(KindSkill, core.Fields{
				// Skills have no update endpoint and a version upload carries
				// only files, so a title change cannot be applied at all; and
				// the `name` in SKILL.md is the skill's identity, which a new
				// version may not change.
				"display_name": {Immutable: true, Summary: true},
				"name":         {Immutable: true, Summary: true},
			}),
		},
		core.KindSpec{
			Kind:     KindEnvironment,
			IDPrefix: "env",
			Stamp:    updatedAt,
			Destroy:  byArchiving,
			Build:    build(prose{field: "description"}),
			Fields: fromAPI(KindEnvironment, with(named, archivable, metadata, core.Fields{
				"description": {Clearable: true},
			})),
		},
		core.KindSpec{
			// The store is the container; the memories inside it are the
			// agent's runtime data and are never read, diffed or removed here.
			Kind:     KindMemoryStore,
			IDPrefix: "memstore",
			Stamp:    updatedAt,
			Destroy:  byArchiving,
			Build:    build(prose{field: "description"}),
			Fields: fromAPI(KindMemoryStore, with(named, archivable, metadata, core.Fields{
				// Rendered into the system prompt of every session the store is
				// attached to, so it is worth showing in the plan.
				"description": {Clearable: true, Summary: true},
			})),
		},
		core.KindSpec{
			// The vault is the container; the credentials inside it are added by
			// whoever owns the secrets and are never read or diffed here. Archiving
			// is the exception to the note on byArchiving: it purges their secrets.
			Kind:     KindVault,
			IDPrefix: "vlt",
			Stamp:    updatedAt,
			Destroy:  byArchiving,
			Build:    buildVault,
			Fields: fromAPI(KindVault, with(archivable, metadata, core.Fields{
				// A vault has no `name`; this is what the API calls it.
				"display_name": {Summary: true},
			})),
		},
		core.KindSpec{
			Kind:                 KindAgent,
			IDPrefix:             "agent",
			VersionField:         "version",
			VersionIsInt:         true,
			UpdateNeedsVersion:   true,
			Stamp:                agentVersion,
			StampMovesOnlyOnEdit: true,
			StampName:            "version",
			Destroy:              byArchiving,
			Build:                build(prose{field: "system"}),
			Fields: fromAPI(KindAgent, with(named, archivable, metadata, core.Fields{
				// `model: x` is accepted, but read back as `{id: x, speed: …}`.
				"model":       {Shorthand: "id", Summary: true},
				"description": {Clearable: true},
				"system":      {Clearable: true},
				// A plan pairs list elements up before diffing them; these are
				// the members that say which element is which. Both lists also
				// take `- ./path.yml` entries spliced in from a data file, for
				// tool specs a framework generates rather than a person types.
				"tools":       {Clearable: true, MatchBy: []string{"name", "type"}, Include: true},
				"mcp_servers": {Clearable: true, MatchBy: []string{"name", "url"}, Include: true},
				"multiagent":  {Clearable: true},
				"skills":      {Clearable: true, MatchBy: []string{"skill_id"}},
			})),
		},
		core.KindSpec{
			Kind:     KindDeployment,
			IDPrefix: "depl",
			Stamp:    updatedAt,
			Destroy:  byArchiving,
			Build:    build(firstMessage),
			Fields: fromAPI(KindDeployment, with(named, archivable, metadata, core.Fields{
				"description": {Clearable: true},
				// A mixed list: files and repositories are written as the API
				// takes them and pass through; a memory store may be a path to
				// its file, optionally with `access` and `instructions` alongside.
				"resources":      {Clearable: true, MatchBy: []string{"memory_store_id", "url", "type"}},
				"vault_ids":      {Clearable: true},
				"schedule":       {Summary: true},
				"agent":          {Summary: true},
				"environment_id": {Summary: true},
			})),
		},
	)
}
