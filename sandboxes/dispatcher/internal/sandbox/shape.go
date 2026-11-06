package sandbox

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	sandboxGroup    = "agents.x-k8s.io"
	extensionsGroup = "extensions.agents.x-k8s.io"
)

// Shape is the wire layout of the Agent Sandbox CRDs at one API version.
// The GKE managed add-on serves v1alpha1; the upstream controller serves
// v1beta1. The two differ in how a claim names its template and how a
// Sandbox is suspended.
type Shape struct {
	version string
}

// ShapeFor returns the Shape for an agents.x-k8s.io API version.
func ShapeFor(version string) (Shape, error) {
	switch version {
	case "v1alpha1", "v1beta1":
		return Shape{version: version}, nil
	}
	return Shape{}, fmt.Errorf("unsupported agents.x-k8s.io API version %q (want v1alpha1 or v1beta1)", version)
}

func (s Shape) Version() string { return s.version }

// Direct reports whether this version drives Sandbox objects directly rather
// than through SandboxClaims. The GKE add-on (v1alpha1) bans the projected
// token a claim's pod would prove itself with and never reconciles warm
// pools, so its sandboxes are created and authenticated the direct way.
func (s Shape) Direct() bool { return s.version == "v1alpha1" }

func (s Shape) Claims() schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: extensionsGroup, Version: s.version, Resource: "sandboxclaims"}
}

func (s Shape) Templates() schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: extensionsGroup, Version: s.version, Resource: "sandboxtemplates"}
}

func (s Shape) WarmPools() schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: extensionsGroup, Version: s.version, Resource: "sandboxwarmpools"}
}

func (s Shape) Sandboxes() schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: sandboxGroup, Version: s.version, Resource: "sandboxes"}
}

// claimSpec is a v1beta1 claim spec with no lifecycle: a shutdownTime or TTL
// would delete the session's volume. Only the upstream path creates claims;
// the add-on creates Sandboxes directly.
func (s Shape) claimSpec(warmPool string) map[string]any {
	return map[string]any{"warmPoolRef": map[string]any{"name": warmPool}}
}

func (s Shape) suspendPatch(suspended bool) map[string]any {
	if s.version == "v1alpha1" {
		replicas := int64(1)
		if suspended {
			replicas = 0
		}
		return map[string]any{"spec": map[string]any{"replicas": replicas}}
	}
	mode := "Running"
	if suspended {
		mode = "Suspended"
	}
	return map[string]any{"spec": map[string]any{"operatingMode": mode}}
}

// isReady reports whether the Sandbox's Ready condition is True. v1alpha1
// (the GKE add-on) is not trusted to report the condition, so there the pod
// phase alone decides.
func (s Shape) isReady(sandbox *unstructured.Unstructured) bool {
	if s.version == "v1alpha1" {
		return true
	}
	conditions, _, _ := unstructured.NestedSlice(sandbox.Object, "status", "conditions")
	for _, raw := range conditions {
		c, _ := raw.(map[string]any)
		if c["type"] == "Ready" {
			return c["status"] == "True"
		}
	}
	return false
}

func (s Shape) isSuspended(sandbox *unstructured.Unstructured) bool {
	suspended, found := s.suspendState(sandbox)
	return found && suspended
}

// suspendState reads the field suspendPatch writes: whether it says
// suspended, and whether the served object carries it at all.
func (s Shape) suspendState(sandbox *unstructured.Unstructured) (suspended, found bool) {
	if s.version == "v1alpha1" {
		replicas, found, err := unstructured.NestedInt64(sandbox.Object, "spec", "replicas")
		return replicas == 0, found && err == nil
	}
	mode, found, _ := unstructured.NestedString(sandbox.Object, "spec", "operatingMode")
	return mode == "Suspended", found
}
