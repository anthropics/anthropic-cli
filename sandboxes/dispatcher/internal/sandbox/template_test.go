package sandbox

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// With no claim configured, a built Sandbox carries no volumeClaimTemplates
// (the add-on's SandboxTemplate schema stores none, so none can be copied)
// and the workspace mount is backed by an explicit emptyDir.
func TestSandboxFromTemplateBacksTheWorkspaceWithAnEmptyDir(t *testing.T) {
	shape, _ := ShapeFor("v1alpha1")
	tpl := object("extensions.agents.x-k8s.io/v1alpha1", "SandboxTemplate", "tpl", "tpl-uid", map[string]any{
		"spec": map[string]any{
			"volumeClaimTemplates": []any{map[string]any{"metadata": map[string]any{"name": WorkspaceVolume}}},
			"podTemplate": map[string]any{"spec": map[string]any{
				"containers": []any{
					map[string]any{"name": "fetch"},
					map[string]any{"name": "worker", "volumeMounts": []any{map[string]any{"name": WorkspaceVolume, "mountPath": "/workspace"}}},
				},
				"volumes": []any{map[string]any{"name": "tmp", "emptyDir": map[string]any{}}},
			}},
		},
	})
	obj, err := SandboxFromTemplate(shape, tpl, ns, "sb", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if vct, found, _ := unstructured.NestedSlice(obj.Object, "spec", "volumeClaimTemplates"); found {
		t.Errorf("spec.volumeClaimTemplates = %v; want none copied from the template", vct)
	}
	volumes, _, _ := unstructured.NestedSlice(obj.Object, "spec", "podTemplate", "spec", "volumes")
	var workspace map[string]any
	for _, raw := range volumes {
		if v, ok := raw.(map[string]any); ok && v["name"] == WorkspaceVolume {
			workspace = v
		}
	}
	if workspace == nil {
		t.Fatalf("volumes = %v; want an explicit %s volume backing the worker's mount", volumes, WorkspaceVolume)
	}
	if _, found, _ := unstructured.NestedMap(workspace, "emptyDir"); !found {
		t.Errorf("workspace volume = %v; want an emptyDir", workspace)
	}
	// Compare with the literal, not workspaceSizeLimit: the cap must match the
	// recipe's default claim size, so a changed constant fails here.
	if limit, _, _ := unstructured.NestedString(workspace, "emptyDir", "sizeLimit"); limit != "10Gi" {
		t.Errorf("workspace emptyDir sizeLimit = %q; want 10Gi", limit)
	}
	if len(volumes) != 2 {
		t.Errorf("volumes = %v; want the template's tmp volume kept beside the workspace", volumes)
	}
}

// A configured claim becomes the Sandbox's own volumeClaimTemplates entry,
// which the controller binds to the worker's workspace mount by name, so
// no emptyDir stands in and the workspace outlives a suspend.
func TestSandboxFromTemplateGivesTheWorkspaceAClaim(t *testing.T) {
	shape, _ := ShapeFor("v1alpha1")
	tpl := object("extensions.agents.x-k8s.io/v1alpha1", "SandboxTemplate", "tpl", "tpl-uid", map[string]any{
		"spec": map[string]any{
			"podTemplate": map[string]any{"spec": map[string]any{
				"containers": []any{
					map[string]any{"name": "fetch"},
					map[string]any{"name": "worker", "volumeMounts": []any{map[string]any{"name": WorkspaceVolume, "mountPath": "/workspace"}}},
				},
				"volumes": []any{map[string]any{"name": "tmp", "emptyDir": map[string]any{}}},
			}},
		},
	})
	claim := &WorkspaceClaim{Size: resource.MustParse("10Gi"), StorageClass: "fast"}
	obj, err := SandboxFromTemplate(shape, tpl, ns, "sb", nil, claim)
	if err != nil {
		t.Fatal(err)
	}
	vct, found, _ := unstructured.NestedSlice(obj.Object, "spec", "volumeClaimTemplates")
	if !found || len(vct) != 1 {
		t.Fatalf("spec.volumeClaimTemplates = %v; want exactly the workspace claim", vct)
	}
	entry := vct[0].(map[string]any)
	if name, _, _ := unstructured.NestedString(entry, "metadata", "name"); name != WorkspaceVolume {
		t.Errorf("claim name = %q, want %q so the controller binds it to the worker's mount", name, WorkspaceVolume)
	}
	if size, _, _ := unstructured.NestedString(entry, "spec", "resources", "requests", "storage"); size != "10Gi" {
		t.Errorf("claim size = %q, want 10Gi", size)
	}
	if class, _, _ := unstructured.NestedString(entry, "spec", "storageClassName"); class != "fast" {
		t.Errorf("claim storageClassName = %q, want fast", class)
	}
	if modes, _, _ := unstructured.NestedStringSlice(entry, "spec", "accessModes"); len(modes) != 1 || modes[0] != "ReadWriteOnce" {
		t.Errorf("claim accessModes = %v, want [ReadWriteOnce]", modes)
	}
	volumes, _, _ := unstructured.NestedSlice(obj.Object, "spec", "podTemplate", "spec", "volumes")
	for _, raw := range volumes {
		if v, ok := raw.(map[string]any); ok && v["name"] == WorkspaceVolume {
			t.Errorf("volumes = %v; a claim-backed workspace must not also get an emptyDir", volumes)
		}
	}
	if len(volumes) != 1 {
		t.Errorf("volumes = %v; want only the template's tmp volume", volumes)
	}
}

// A claim with no storage class leaves storageClassName unset, so the
// cluster's default class applies.
func TestSandboxFromTemplateClaimTakesTheDefaultStorageClass(t *testing.T) {
	shape, _ := ShapeFor("v1alpha1")
	tpl := object("extensions.agents.x-k8s.io/v1alpha1", "SandboxTemplate", "tpl", "tpl-uid", map[string]any{
		"spec": map[string]any{"podTemplate": map[string]any{"spec": map[string]any{
			"containers": []any{map[string]any{"name": "fetch"}, map[string]any{"name": "worker"}},
		}}},
	})
	obj, err := SandboxFromTemplate(shape, tpl, ns, "sb", nil, &WorkspaceClaim{Size: resource.MustParse("1Gi")})
	if err != nil {
		t.Fatal(err)
	}
	vct, _, _ := unstructured.NestedSlice(obj.Object, "spec", "volumeClaimTemplates")
	if _, found, _ := unstructured.NestedString(vct[0].(map[string]any), "spec", "storageClassName"); found {
		t.Errorf("claim = %v; want no storageClassName so the default class applies", vct[0])
	}
}

// A template that defines its own workspace volume keeps it: the operator's
// backing wins over the dispatcher's claim or emptyDir stand-in.
func TestSandboxFromTemplateKeepsAnExplicitWorkspaceVolume(t *testing.T) {
	shape, _ := ShapeFor("v1alpha1")
	tpl := object("extensions.agents.x-k8s.io/v1alpha1", "SandboxTemplate", "tpl", "tpl-uid", map[string]any{
		"spec": map[string]any{
			"podTemplate": map[string]any{"spec": map[string]any{
				"containers": []any{map[string]any{"name": "fetch"}, map[string]any{"name": "worker"}},
				"volumes":    []any{map[string]any{"name": WorkspaceVolume, "persistentVolumeClaim": map[string]any{"claimName": "shared-ws"}}},
			}},
		},
	})
	obj, err := SandboxFromTemplate(shape, tpl, ns, "sb", nil, &WorkspaceClaim{Size: resource.MustParse("10Gi")})
	if err != nil {
		t.Fatal(err)
	}
	if vct, found, _ := unstructured.NestedSlice(obj.Object, "spec", "volumeClaimTemplates"); found {
		t.Errorf("spec.volumeClaimTemplates = %v; the template's own workspace volume must win over the configured claim", vct)
	}
	volumes, _, _ := unstructured.NestedSlice(obj.Object, "spec", "podTemplate", "spec", "volumes")
	if len(volumes) != 1 {
		t.Fatalf("volumes = %v; want only the template's own workspace volume", volumes)
	}
	claim, _, _ := unstructured.NestedString(volumes[0].(map[string]any), "persistentVolumeClaim", "claimName")
	if claim != "shared-ws" {
		t.Errorf("workspace volume = %v; want the template's persistentVolumeClaim kept", volumes[0])
	}
}
