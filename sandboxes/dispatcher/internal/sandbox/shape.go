package sandbox

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

const (
	sandboxGroup    = "agents.x-k8s.io"
	extensionsGroup = "extensions.agents.x-k8s.io"
)

const (
	v1alpha1 = "v1alpha1"
	v1beta1  = "v1beta1"
)

// ClaimMode is how a session gets its Sandbox.
type ClaimMode string

const (
	// ClaimModeSandboxClaim takes the Sandbox from a warm pool through a
	// SandboxClaim. Its pod proves itself with a projected token.
	ClaimModeSandboxClaim ClaimMode = "sandbox_claim"
	// ClaimModeDirect creates the Sandbox with no claim. Its pod proves itself
	// with a nonce. The GKE add-on needs it, because it bans projected tokens.
	ClaimModeDirect ClaimMode = "direct"
)

// Shape is how the dispatcher uses the Agent Sandbox CRDs. The version sets how
// a Sandbox is suspended. The mode sets which object stands for a session.
type Shape struct {
	version string
	mode    ClaimMode
}

// ShapeFor returns the Shape for an agents.x-k8s.io API version and a claim
// mode. An empty mode is direct on v1alpha1 and sandbox_claim on v1beta1. The
// dispatcher writes no v1alpha1 claim, so sandbox_claim is refused there.
func ShapeFor(version string, mode ClaimMode) (Shape, error) {
	if version != v1alpha1 && version != v1beta1 {
		return Shape{}, fmt.Errorf("unsupported agents.x-k8s.io API version %q (want v1alpha1 or v1beta1)", version)
	}
	switch {
	case mode == "" && version == v1alpha1:
		mode = ClaimModeDirect
	case mode == "":
		mode = ClaimModeSandboxClaim
	case mode != ClaimModeDirect && mode != ClaimModeSandboxClaim:
		return Shape{}, fmt.Errorf("unsupported claim mode %q (want %s or %s)", mode, ClaimModeSandboxClaim, ClaimModeDirect)
	case mode == ClaimModeSandboxClaim && version == v1alpha1:
		return Shape{}, fmt.Errorf("claim mode %s needs the v1beta1 API, and the version in use is v1alpha1; set -claim-mode %s", mode, ClaimModeDirect)
	}
	return Shape{version: version, mode: mode}, nil
}

// ServedVersion returns the API version at which the cluster serves both
// Sandboxes and SandboxTemplates. It prefers v1beta1, because a cluster that
// serves both versions is on its way to dropping v1alpha1.
func ServedVersion(d discovery.DiscoveryInterface) (string, error) {
	for _, version := range []string{v1beta1, v1alpha1} {
		s := Shape{version: version}
		served, err := discovery.IsResourceEnabled(d, s.Sandboxes())
		if served {
			served, err = discovery.IsResourceEnabled(d, s.Templates())
		}
		if err != nil {
			return "", fmt.Errorf("asking the cluster about Agent Sandbox %s: %w", version, err)
		}
		if served {
			return version, nil
		}
	}
	return "", errors.New("the cluster serves sandboxes and sandboxtemplates at neither v1beta1 nor v1alpha1; is Agent Sandbox installed?")
}

func (s Shape) Version() string { return s.version }

// ClaimMode returns the mode in use, which ShapeFor never leaves empty.
func (s Shape) ClaimMode() ClaimMode { return s.mode }

// Direct reports whether the claim mode is ClaimModeDirect.
func (s Shape) Direct() bool { return s.mode == ClaimModeDirect }

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
// would delete the session's volume. Only the claim path creates claims.
func (s Shape) claimSpec(warmPool string) map[string]any {
	return map[string]any{"warmPoolRef": map[string]any{"name": warmPool}}
}

// suspendSpec is the part of a Sandbox's spec that says whether it is suspended.
func (s Shape) suspendSpec(suspended bool) map[string]any {
	if s.version == v1alpha1 {
		replicas := int64(1)
		if suspended {
			replicas = 0
		}
		return map[string]any{"replicas": replicas}
	}
	mode := "Running"
	if suspended {
		mode = "Suspended"
	}
	return map[string]any{"operatingMode": mode}
}

// isReady reports whether the Sandbox's Ready condition is True. The GKE add-on
// is not trusted to report the condition, so on the direct path the pod phase
// alone decides.
func (s Shape) isReady(sandbox *unstructured.Unstructured) bool {
	if s.Direct() {
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

// suspendState reads the field suspendSpec holds: whether it says
// suspended, and whether the served object carries it at all.
func (s Shape) suspendState(sandbox *unstructured.Unstructured) (suspended, found bool) {
	if s.version == v1alpha1 {
		replicas, found, err := unstructured.NestedInt64(sandbox.Object, "spec", "replicas")
		return replicas == 0, found && err == nil
	}
	mode, found, _ := unstructured.NestedString(sandbox.Object, "spec", "operatingMode")
	return mode == "Suspended", found
}
