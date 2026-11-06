package sandbox

import (
	"context"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	kubefake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

// directClient builds a Client on the add-on path with a SandboxTemplate that
// has a fetch container, returning its nonce registry and the Kubernetes fake
// for inspection.
func directClient(t *testing.T, objects ...runtime.Object) (*Client, *dynamicfake.FakeDynamicClient, *fakeNonces, *kubefake.Clientset) {
	t.Helper()
	return directClientWithFetchEnv(t, []any{map[string]any{"name": "CMA_DISPATCHER_URL", "value": "https://d"}}, objects...)
}

// directClientWithFetchEnv is directClient with the template's fetch
// container carrying fetchEnv.
func directClientWithFetchEnv(t *testing.T, fetchEnv []any, objects ...runtime.Object) (*Client, *dynamicfake.FakeDynamicClient, *fakeNonces, *kubefake.Clientset) {
	t.Helper()
	shape, _ := ShapeFor("v1alpha1")
	tpl := object("extensions.agents.x-k8s.io/v1alpha1", "SandboxTemplate", "tpl", "tpl-uid", map[string]any{
		"spec": map[string]any{"podTemplate": map[string]any{"spec": map[string]any{"containers": []any{
			map[string]any{"name": "fetch", "env": fetchEnv},
			map[string]any{"name": "worker"},
		}}}},
	})
	listKinds := map[schema.GroupVersionResource]string{
		shape.Templates(): "SandboxTemplateList",
		shape.Sandboxes(): "SandboxList",
	}
	dyn := dynamicfake.NewSimpleDynamicClientWithCustomListKinds(runtime.NewScheme(), listKinds)
	seed := append([]runtime.Object{tpl}, objects...)
	for _, o := range seed {
		u := o.(*unstructured.Unstructured)
		gvr := shape.Templates()
		if u.GetKind() == "Sandbox" {
			gvr = shape.Sandboxes()
		}
		if _, err := dyn.Resource(gvr).Namespace(ns).Create(context.Background(), u, metav1.CreateOptions{}); err != nil {
			t.Fatal(err)
		}
	}
	nonces := newFakeNonces()
	kube := kubefake.NewSimpleClientset()
	c, err := New(dyn, kube, Options{Shape: shape, Namespace: ns, EnvironmentID: "env_1", TemplateName: "tpl", Nonces: nonces})
	if err != nil {
		t.Fatal(err)
	}
	return c, dyn, nonces, kube
}

// nonceEnvEntries returns the fetch container's EnvAssignmentNonce env
// entries from a stored Sandbox.
func nonceEnvEntries(t *testing.T, obj *unstructured.Unstructured) []map[string]any {
	t.Helper()
	containers, _, _ := unstructured.NestedSlice(obj.Object, "spec", "podTemplate", "spec", "containers")
	var entries []map[string]any
	for _, raw := range containers {
		container, ok := raw.(map[string]any)
		if !ok || container["name"] != FetchContainer {
			continue
		}
		env, _, _ := unstructured.NestedSlice(container, "env")
		for _, e := range env {
			if entry, ok := e.(map[string]any); ok && entry["name"] == EnvAssignmentNonce {
				entries = append(entries, entry)
			}
		}
	}
	return entries
}

// On the add-on path EnsureClaim creates a Sandbox directly from the
// template, one replica and Retain, plus a nonce Secret owned by the Sandbox
// whose value is registered under the sandbox name; the fetch container reads
// it through an env valueFrom so the nonce never appears in the pod spec.
func TestEnsureClaimCreatesTheSandboxDirectlyOnTheAddOn(t *testing.T) {
	c, dyn, nonces, kube := directClient(t)
	// The API server assigns a created Sandbox its UID. The fake does not, so
	// the reactor stands in.
	dyn.PrependReactor("create", "sandboxes", func(action k8stesting.Action) (bool, runtime.Object, error) {
		action.(k8stesting.CreateAction).GetObject().(*unstructured.Unstructured).SetUID("sb-uid")
		return false, nil, nil
	})
	claim, err := c.EnsureClaim(context.Background(), "session_01X")
	if err != nil {
		t.Fatal(err)
	}
	if claim.SandboxName != claim.Name || claim.SessionID != "session_01X" {
		t.Errorf("handle = %+v; the Sandbox is its own bound sandbox", claim)
	}
	stored, err := dyn.Resource(c.shape.Sandboxes()).Namespace(ns).Get(context.Background(), claim.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if labels := stored.GetLabels(); labels[LabelSessionID] != "session_01X" || labels[LabelManagedBy] != ManagedByValue || labels[LabelEnvironmentID] != "env_1" {
		t.Errorf("labels = %v", labels)
	}
	replicas, _, _ := unstructured.NestedInt64(stored.Object, "spec", "replicas")
	policy, _, _ := unstructured.NestedString(stored.Object, "spec", "shutdownPolicy")
	if replicas != 1 || policy != "Retain" {
		t.Errorf("spec.replicas = %d, spec.shutdownPolicy = %q; want 1 and Retain", replicas, policy)
	}
	entries := nonceEnvEntries(t, stored)
	if len(entries) != 1 {
		t.Fatalf("fetch carries %d %s entries, want exactly one", len(entries), EnvAssignmentNonce)
	}
	ref, _, _ := unstructured.NestedMap(entries[0], "valueFrom", "secretKeyRef")
	if ref["name"] != NonceSecretName(claim.Name) || ref["key"] != NonceSecretKey {
		t.Errorf("secretKeyRef = %v; want the sandbox's nonce Secret and key %q", ref, NonceSecretKey)
	}
	if value, found, _ := unstructured.NestedString(entries[0], "value"); found {
		t.Errorf("the env entry carries plaintext value %q; the nonce must live only in the Secret", value)
	}
	secret, err := kube.CoreV1().Secrets(ns).Get(context.Background(), NonceSecretName(claim.Name), metav1.GetOptions{})
	if err != nil {
		t.Fatalf("the nonce Secret was not created: %v", err)
	}
	if got := secret.StringData[NonceSecretKey]; got == "" || got != nonces.get(claim.Name) {
		t.Errorf("Secret nonce %q, registered %q; want them equal and non-empty", got, nonces.get(claim.Name))
	}
	// Garbage collection matches an owner by UID, not by name.
	if len(secret.OwnerReferences) != 1 || secret.OwnerReferences[0].Kind != "Sandbox" || secret.OwnerReferences[0].Name != claim.Name || secret.OwnerReferences[0].UID != "sb-uid" {
		t.Errorf("Secret ownerReferences = %v; want the Sandbox by name and UID, so the Secret is garbage-collected with it", secret.OwnerReferences)
	}
}

// A template whose fetch container already carries an EnvAssignmentNonce
// entry must not end up with two: the kubelet resolves duplicate env names
// last-wins, so a kept entry could diverge from what the registry accepts.
func TestEnsureClaimReplacesAPresetNonceEntry(t *testing.T) {
	preset := []any{
		map[string]any{"name": "CMA_DISPATCHER_URL", "value": "https://d"},
		map[string]any{"name": EnvAssignmentNonce, "value": "stale-template-nonce"},
	}
	c, dyn, _, _ := directClientWithFetchEnv(t, preset)
	claim, err := c.EnsureClaim(context.Background(), "session_01X")
	if err != nil {
		t.Fatal(err)
	}
	stored, err := dyn.Resource(c.shape.Sandboxes()).Namespace(ns).Get(context.Background(), claim.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	entries := nonceEnvEntries(t, stored)
	if len(entries) != 1 {
		t.Fatalf("fetch carries %d %s entries, want the template's replaced by exactly one", len(entries), EnvAssignmentNonce)
	}
	if _, found, _ := unstructured.NestedString(entries[0], "value"); found {
		t.Errorf("entry = %v; want the template's plaintext entry replaced by the secretKeyRef", entries[0])
	}
}

// While the registry still holds a Sandbox's nonce, EnsureClaim adopts the
// existing Sandbox without recreating it.
func TestEnsureClaimAdoptsASandboxWhoseNonceIsKnown(t *testing.T) {
	name := SandboxName("session_01X")
	existing := object("agents.x-k8s.io/v1alpha1", "Sandbox", name, "sb-uid", map[string]any{
		"spec": map[string]any{"podTemplate": map[string]any{"spec": map[string]any{"containers": []any{map[string]any{"name": "fetch"}}}}},
	})
	existing.SetLabels(map[string]string{LabelSessionID: "session_01X", LabelEnvironmentID: "env_1", LabelManagedBy: ManagedByValue})
	c, dyn, nonces, _ := directClient(t, existing)
	nonces.Set(name, "the-armed-nonce")
	claim, err := c.EnsureClaim(context.Background(), "session_01X")
	if err != nil {
		t.Fatal(err)
	}
	if claim.Name != name || nonces.get(name) != "the-armed-nonce" {
		t.Errorf("claim = %+v, nonce = %q; want the Sandbox adopted with its nonce kept", claim, nonces.get(name))
	}
	if _, err := dyn.Resource(c.shape.Sandboxes()).Namespace(ns).Get(context.Background(), name, metav1.GetOptions{}); err != nil {
		t.Errorf("the adopted Sandbox must survive: %v", err)
	}
}

// A restarted dispatcher holds no nonce for an existing Sandbox and may not
// read the Secret back. A podless Sandbox with no recorded handoff holds
// nothing of the session, so EnsureClaim deletes it and fails the handoff;
// the redelivered work item recreates it with a fresh nonce. The delete does
// not read the dispatcher's current claim setting: here one is configured.
func TestEnsureClaimDeletesASandboxWhoseNonceIsUnknown(t *testing.T) {
	name := SandboxName("session_01X")
	existing := object("agents.x-k8s.io/v1alpha1", "Sandbox", name, "sb-uid", map[string]any{
		"spec": map[string]any{"podTemplate": map[string]any{"spec": map[string]any{"containers": []any{map[string]any{"name": "fetch"}}}}},
	})
	existing.SetLabels(map[string]string{LabelSessionID: "session_01X", LabelEnvironmentID: "env_1", LabelManagedBy: ManagedByValue})
	c, dyn, nonces, _ := directClient(t, existing)
	c.workspaceClaim = &WorkspaceClaim{Size: resource.MustParse("10Gi")}
	if _, err := c.EnsureClaim(context.Background(), "session_01X"); err == nil {
		t.Fatal("a Sandbox whose nonce the registry does not hold must not be adopted")
	}
	if _, err := dyn.Resource(c.shape.Sandboxes()).Namespace(ns).Get(context.Background(), name, metav1.GetOptions{}); !apierrors.IsNotFound(err) {
		t.Errorf("get after refusal = %v; want the Sandbox deleted so the retry recreates it", err)
	}
	if nonces.Known(name) {
		t.Error("a refused Sandbox must leave no nonce registered")
	}
}

// A claim-backed Sandbox may hold a turn's files before its handoff stamp
// lands, so a restarted dispatcher never deletes one: podless and unstamped,
// it is suspended and re-armed like a retained one, and its volume survives.
// The claim is read from the stored Sandbox. Here the dispatcher restarts
// with no claim configured, as after an operator turned the claim off.
func TestEnsureClaimKeepsAnUnstampedClaimBackedSandbox(t *testing.T) {
	name := SandboxName("session_01X")
	existing := object("agents.x-k8s.io/v1alpha1", "Sandbox", name, "sb-uid", map[string]any{
		"spec": map[string]any{
			"replicas":             int64(1),
			"podTemplate":          map[string]any{"spec": map[string]any{"containers": []any{map[string]any{"name": "fetch"}}}},
			"volumeClaimTemplates": []any{(&WorkspaceClaim{Size: resource.MustParse("10Gi")}).template()},
		},
	})
	existing.SetLabels(map[string]string{LabelSessionID: "session_01X", LabelEnvironmentID: "env_1", LabelManagedBy: ManagedByValue})
	c, dyn, nonces, _ := directClient(t, existing)
	claim, err := c.EnsureClaim(context.Background(), "session_01X")
	if err != nil {
		t.Fatal(err)
	}
	if claim.Name != name || !nonces.Known(name) {
		t.Errorf("claim = %+v, nonce known = %v; want the unstamped claim-backed Sandbox adopted and re-armed", claim, nonces.Known(name))
	}
	stored, err := dyn.Resource(c.shape.Sandboxes()).Namespace(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("a claim-backed Sandbox must survive a restart, stamped or not: %v", err)
	}
	if !c.shape.isSuspended(stored) {
		t.Error("a podless Sandbox must be suspended so the handoff's resume builds the pod with the new nonce")
	}
}

// A served schema that prunes volumeClaimTemplates leaves the worker's
// workspace mount naming a volume nothing defines, so the create fails
// naming the api-version and leaves no Sandbox behind.
func TestEnsureClaimFailsWhenTheClusterDropsTheWorkspaceClaim(t *testing.T) {
	c, dyn, nonces, _ := directClient(t)
	c.workspaceClaim = &WorkspaceClaim{Size: resource.MustParse("10Gi")}
	dyn.PrependReactor("create", "sandboxes", func(action k8stesting.Action) (bool, runtime.Object, error) {
		stored := action.(k8stesting.CreateAction).GetObject().DeepCopyObject().(*unstructured.Unstructured)
		unstructured.RemoveNestedField(stored.Object, "spec", "volumeClaimTemplates")
		return true, stored, nil
	})
	name := SandboxName("session_01X")
	if _, err := c.EnsureClaim(context.Background(), "session_01X"); err == nil || !strings.Contains(err.Error(), "api-version") {
		t.Fatalf("err = %v; want the create refused naming the api-version", err)
	}
	if _, err := dyn.Resource(c.shape.Sandboxes()).Namespace(ns).Get(context.Background(), name, metav1.GetOptions{}); !apierrors.IsNotFound(err) {
		t.Errorf("get after refusal = %v; want no Sandbox left behind", err)
	}
	if nonces.Known(name) {
		t.Error("a refused create must leave no nonce registered")
	}
}

// retainedSandbox is a Sandbox past its first collection (the handoff is
// stamped on it), whose fetch container reads the nonce Secret from before a
// dispatcher restart. suspended is spec.replicas 0.
func retainedSandbox(name string, suspended bool) *unstructured.Unstructured {
	replicas := int64(1)
	if suspended {
		replicas = 0
	}
	sb := object("agents.x-k8s.io/v1alpha1", "Sandbox", name, "sb-uid", map[string]any{
		"spec": map[string]any{
			"replicas":       replicas,
			"shutdownPolicy": "Retain",
			"podTemplate": map[string]any{"spec": map[string]any{"containers": []any{
				map[string]any{"name": "fetch", "env": []any{map[string]any{
					"name":      EnvAssignmentNonce,
					"valueFrom": map[string]any{"secretKeyRef": map[string]any{"name": NonceSecretName(name), "key": NonceSecretKey}},
				}}},
				map[string]any{"name": "worker"},
			}}},
		},
	})
	sb.SetLabels(map[string]string{LabelSessionID: "session_01X", LabelEnvironmentID: "env_1", LabelManagedBy: ManagedByValue})
	sb.SetAnnotations(map[string]string{AnnotationHandoffPodUID: "pod-1", AnnotationHandoffRestarts: "0"})
	return sb
}

// A suspended Sandbox past its first collection holds the session's
// workspace, so a restarted dispatcher keeps it and re-arms it: a fresh
// nonce in a new Secret owned by the Sandbox, with the fetch container's
// env pointed at that Secret so the resumed pod starts with the new nonce.
func TestEnsureClaimRearmsASuspendedRetainedSandbox(t *testing.T) {
	name := SandboxName("session_01X")
	c, dyn, nonces, kube := directClient(t, retainedSandbox(name, true))
	claim, err := c.EnsureClaim(context.Background(), "session_01X")
	if err != nil {
		t.Fatal(err)
	}
	if claim.Name != name || !nonces.Known(name) {
		t.Errorf("claim = %+v, nonce known = %v; want the Sandbox adopted and re-armed", claim, nonces.Known(name))
	}
	stored, err := dyn.Resource(c.shape.Sandboxes()).Namespace(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("the re-armed Sandbox must survive: %v", err)
	}
	entries := nonceEnvEntries(t, stored)
	if len(entries) != 1 {
		t.Fatalf("fetch carries %d %s entries, want exactly one", len(entries), EnvAssignmentNonce)
	}
	ref, _, _ := unstructured.NestedMap(entries[0], "valueFrom", "secretKeyRef")
	secretName, _ := ref["name"].(string)
	if secretName == NonceSecretName(name) || !strings.HasPrefix(secretName, NonceSecretName(name)+"-") || ref["key"] != NonceSecretKey {
		t.Fatalf("secretKeyRef = %v; want a fresh Secret named under %s-, key %q", ref, NonceSecretName(name), NonceSecretKey)
	}
	secret, err := kube.CoreV1().Secrets(ns).Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("the re-armed nonce Secret was not created: %v", err)
	}
	if got := secret.StringData[NonceSecretKey]; got == "" || got != nonces.get(name) {
		t.Errorf("Secret nonce %q, registered %q; want them equal and non-empty", got, nonces.get(name))
	}
	// Garbage collection matches an owner by UID, not by name.
	if len(secret.OwnerReferences) != 1 || secret.OwnerReferences[0].Kind != "Sandbox" || secret.OwnerReferences[0].Name != name || secret.OwnerReferences[0].UID != "sb-uid" {
		t.Errorf("Secret ownerReferences = %v; want the Sandbox by name and UID, so the Secret is garbage-collected with it", secret.OwnerReferences)
	}
}

// A Sandbox with no recorded handoff whose pod exists may have collected
// work whose stamp never landed, so a restarted dispatcher must not delete
// it out from under a possibly running turn.
func TestEnsureClaimKeepsAnUnstampedSandboxWhosePodExists(t *testing.T) {
	name := SandboxName("session_01X")
	existing := object("agents.x-k8s.io/v1alpha1", "Sandbox", name, "sb-uid", map[string]any{
		"spec": map[string]any{"replicas": int64(1), "podTemplate": map[string]any{"spec": map[string]any{"containers": []any{map[string]any{"name": "fetch"}}}}},
	})
	existing.SetLabels(map[string]string{LabelSessionID: "session_01X", LabelEnvironmentID: "env_1", LabelManagedBy: ManagedByValue})
	c, dyn, nonces, kube := directClient(t, existing)
	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns, UID: "pod-1",
		OwnerReferences: []metav1.OwnerReference{{Kind: "Sandbox", Name: name, UID: "sb-uid", Controller: ptr(true)}}}}
	if _, err := kube.CoreV1().Pods(ns).Create(context.Background(), pod, metav1.CreateOptions{}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.EnsureClaim(context.Background(), "session_01X"); err == nil {
		t.Fatal("a Sandbox whose pod exists must not be adopted while its nonce is unknown")
	}
	if _, err := dyn.Resource(c.shape.Sandboxes()).Namespace(ns).Get(context.Background(), name, metav1.GetOptions{}); err != nil {
		t.Errorf("a Sandbox whose pod exists must survive: %v", err)
	}
	if nonces.Known(name) {
		t.Error("a Sandbox that was not adopted must leave no nonce registered")
	}
}

// A retained Sandbox that is not suspended and still has its pod may be
// running a turn, so a restarted dispatcher neither deletes nor re-arms it:
// the handoff fails and retries once the worker's exit suspends the Sandbox.
func TestEnsureClaimLeavesARunningRetainedSandboxUntilItSuspends(t *testing.T) {
	name := SandboxName("session_01X")
	c, dyn, nonces, kube := directClient(t, retainedSandbox(name, false))
	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns, UID: "pod-1",
		OwnerReferences: []metav1.OwnerReference{{Kind: "Sandbox", Name: name, UID: "sb-uid", Controller: ptr(true)}}}}
	if _, err := kube.CoreV1().Pods(ns).Create(context.Background(), pod, metav1.CreateOptions{}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.EnsureClaim(context.Background(), "session_01X"); err == nil {
		t.Fatal("a running Sandbox whose nonce the registry does not hold must not be adopted")
	}
	stored, err := dyn.Resource(c.shape.Sandboxes()).Namespace(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("the running Sandbox must survive: %v", err)
	}
	if c.shape.isSuspended(stored) {
		t.Error("a Sandbox whose pod exists must not be suspended out from under it")
	}
	if secrets, _ := kube.CoreV1().Secrets(ns).List(context.Background(), metav1.ListOptions{}); len(secrets.Items) != 0 {
		t.Errorf("secrets = %v; want none created before the Sandbox suspends", secrets.Items)
	}
	if nonces.Known(name) {
		t.Error("a Sandbox that was not re-armed must leave no nonce registered")
	}
}

// A retained Sandbox that is not suspended but whose pod is gone has no turn
// to disturb and nothing else that would ever suspend it, so the adoption
// suspends it and re-arms in the same pass; the resume then builds a pod
// carrying the new nonce.
func TestEnsureClaimSuspendsAndRearmsAPodlessRetainedSandbox(t *testing.T) {
	name := SandboxName("session_01X")
	c, dyn, nonces, _ := directClient(t, retainedSandbox(name, false))
	claim, err := c.EnsureClaim(context.Background(), "session_01X")
	if err != nil {
		t.Fatal(err)
	}
	if claim.Name != name || !nonces.Known(name) {
		t.Errorf("claim = %+v, nonce known = %v; want the podless Sandbox adopted and re-armed", claim, nonces.Known(name))
	}
	stored, err := dyn.Resource(c.shape.Sandboxes()).Namespace(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("the re-armed Sandbox must survive: %v", err)
	}
	if !c.shape.isSuspended(stored) {
		t.Error("a podless Sandbox must be suspended so the handoff's resume builds the pod with the new nonce")
	}
}

// A re-arm that could not finish must leave the nonce unregistered: a
// registered nonce the pod does not carry would let every later handoff run
// against a pod that can never authenticate.
func TestEnsureClaimRearmFailuresLeaveTheNonceUnregistered(t *testing.T) {
	name := SandboxName("session_01X")

	c, _, nonces, kube := directClient(t, retainedSandbox(name, true))
	kube.PrependReactor("create", "secrets", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, apierrors.NewForbidden(schema.GroupResource{Resource: "secrets"}, "any", nil)
	})
	if _, err := c.EnsureClaim(context.Background(), "session_01X"); err == nil {
		t.Fatal("a re-arm whose Secret create failed must fail the adoption")
	}
	if nonces.Known(name) {
		t.Error("a failed Secret create must leave no nonce registered")
	}

	stale := retainedSandbox(name, true)
	c, dyn, nonces, _ := directClient(t, stale)
	dyn.PrependReactor("patch", "sandboxes", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, stale.DeepCopy(), nil
	})
	if _, err := c.EnsureClaim(context.Background(), "session_01X"); err == nil {
		t.Fatal("a re-arm the cluster did not record must fail the adoption")
	}
	if nonces.Known(name) {
		t.Error("an unrecorded re-arm must leave no nonce registered")
	}
}

// EnsureClaim refuses a Sandbox under the session's name whose labels name
// another session, environment or manager. It leaves the Sandbox in place
// for its owner to remove. Sandbox names derive from the session id alone,
// so two dispatchers serving different environments in one namespace name a
// session's Sandbox alike. The seeded Sandbox is one a restarted dispatcher
// keeps and re-arms. A dropped label clause would hand the session to
// another environment's Sandbox.
func TestEnsureClaimRefusesADirectSandboxOfAnotherOwner(t *testing.T) {
	name := SandboxName("session_01X")
	for tc, labels := range map[string]map[string]string{
		"no labels":         nil,
		"another session":   {LabelSessionID: "session_02Y", LabelEnvironmentID: "env_1", LabelManagedBy: ManagedByValue},
		"another env":       {LabelSessionID: "session_01X", LabelEnvironmentID: "env_2", LabelManagedBy: ManagedByValue},
		"another manager":   {LabelSessionID: "session_01X", LabelEnvironmentID: "env_1", LabelManagedBy: "someone-else"},
		"session label off": {LabelEnvironmentID: "env_1", LabelManagedBy: ManagedByValue},
	} {
		existing := retainedSandbox(name, true)
		existing.SetLabels(labels)
		c, dyn, nonces, _ := directClient(t, existing)
		if _, err := c.EnsureClaim(context.Background(), "session_01X"); err == nil {
			t.Errorf("%s: a Sandbox under the session's name that the dispatcher does not own must be refused", tc)
		}
		if nonces.get(name) != "" {
			t.Errorf("%s: a refused Sandbox must leave no nonce registered", tc)
		}
		if _, err := dyn.Resource(c.shape.Sandboxes()).Namespace(ns).Get(context.Background(), name, metav1.GetOptions{}); err != nil {
			t.Errorf("%s: get after refusal = %v; want the Sandbox left in place", tc, err)
		}
	}
}
