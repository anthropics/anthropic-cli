package sandbox

import (
	"context"
	"reflect"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestEnsureClaimCreatesTheUpstreamClaim(t *testing.T) {
	c, dyn, _ := newClient(t, "v1beta1")
	claim, err := c.EnsureClaim(context.Background(), "session_01X")
	if err != nil {
		t.Fatal(err)
	}
	stored, err := dyn.Resource(c.shape.Claims()).Namespace(ns).Get(context.Background(), claim.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if labels := stored.GetLabels(); labels[LabelSessionID] != "session_01X" || labels[LabelManagedBy] != ManagedByValue || labels[LabelEnvironmentID] != "env_1" {
		t.Errorf("labels = %v", labels)
	}
	wantSpec := map[string]any{"warmPoolRef": map[string]any{"name": "pool"}}
	spec, _, _ := unstructured.NestedMap(stored.Object, "spec")
	if !reflect.DeepEqual(spec, wantSpec) {
		t.Errorf("spec = %v, want exactly %v (no lifecycle expiry)", spec, wantSpec)
	}
}

func TestEnsureClaimReturnsExisting(t *testing.T) {
	name := ClaimName("session_01X")
	existing := object("extensions.agents.x-k8s.io/v1beta1", "SandboxClaim", name, "claim-uid", map[string]any{
		"spec":   map[string]any{"warmPoolRef": map[string]any{"name": "pool"}},
		"status": map[string]any{"sandbox": map[string]any{"name": "sb-1"}},
	})
	existing.SetLabels(map[string]string{LabelSessionID: "session_01X", LabelEnvironmentID: "env_1", LabelManagedBy: ManagedByValue})
	c, _, _ := newClient(t, "v1beta1", existing)
	claim, err := c.EnsureClaim(context.Background(), "session_01X")
	if err != nil {
		t.Fatal(err)
	}
	if claim.SandboxName != "sb-1" || claim.SessionID != "session_01X" {
		t.Errorf("claim = %+v", claim)
	}
}

func TestEnsureClaimRefusesAClaimOfAnotherOwner(t *testing.T) {
	name := ClaimName("session_01X")
	for tc, labels := range map[string]map[string]string{
		"no labels":         nil,
		"another session":   {LabelSessionID: "session_02Y", LabelEnvironmentID: "env_1", LabelManagedBy: ManagedByValue},
		"another env":       {LabelSessionID: "session_01X", LabelEnvironmentID: "env_2", LabelManagedBy: ManagedByValue},
		"another manager":   {LabelSessionID: "session_01X", LabelEnvironmentID: "env_1", LabelManagedBy: "someone-else"},
		"session label off": {LabelEnvironmentID: "env_1", LabelManagedBy: ManagedByValue},
	} {
		existing := object("extensions.agents.x-k8s.io/v1beta1", "SandboxClaim", name, "claim-uid", map[string]any{
			"status": map[string]any{"sandbox": map[string]any{"name": "sb-attacker"}},
		})
		existing.SetLabels(labels)
		c, _, _ := newClient(t, "v1beta1", existing)
		if _, err := c.EnsureClaim(context.Background(), "session_01X"); err == nil {
			t.Errorf("%s: a pre-existing claim under the session's name that the dispatcher does not own must be refused", tc)
		}
	}
}

func TestListAndDeleteClaims(t *testing.T) {
	objects := ownedGraph(map[string]string{LabelManagedBy: ManagedByValue, LabelEnvironmentID: "env_1", LabelSessionID: "session_a"})
	other := object("extensions.agents.x-k8s.io/v1beta1", "SandboxClaim", "foreign", "foreign-uid", nil)
	otherEnvironment := object("extensions.agents.x-k8s.io/v1beta1", "SandboxClaim", "claim-b", "claim-b-uid", nil)
	otherEnvironment.SetLabels(map[string]string{LabelManagedBy: ManagedByValue, LabelEnvironmentID: "env_2", LabelSessionID: "session_b"})
	c, _, _ := newClient(t, "v1beta1", append(objects, other, otherEnvironment)...)
	claims, err := c.ListClaims(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(claims) != 1 || claims[0].Name != "claim-a" {
		t.Errorf("claims = %+v, want only this environment's claim-a", claims)
	}
	if err := c.StampLastActive(context.Background(), "claim-a", time.Date(2026, 9, 3, 12, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	claim, _ := c.GetClaim(context.Background(), "claim-a")
	if claim.LastActive.IsZero() {
		t.Error("last-active annotation not read back")
	}
	if claim.HandoffRestarts != nil {
		t.Error("a claim with no handoff must read back no restart count")
	}
	if err := c.StampHandoff(context.Background(), "claim-a", "pod-uid-1", 2); err != nil {
		t.Fatal(err)
	}
	claim, _ = c.GetClaim(context.Background(), "claim-a")
	if claim.HandoffRestarts == nil || *claim.HandoffRestarts != 2 || claim.HandoffPodUID != "pod-uid-1" {
		t.Errorf("HandoffRestarts = %v, HandoffPodUID = %q; want the stamped handoff read back", claim.HandoffRestarts, claim.HandoffPodUID)
	}
	if err := c.DeleteClaim(context.Background(), "claim-a"); err != nil {
		t.Fatal(err)
	}
	if err := c.DeleteClaim(context.Background(), "claim-a"); err != nil {
		t.Errorf("deleting an absent claim: %v", err)
	}
}
