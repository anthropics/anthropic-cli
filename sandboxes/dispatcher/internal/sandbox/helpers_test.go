package sandbox

import (
	"context"
	"sync"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	kubefake "k8s.io/client-go/kubernetes/fake"
)

const ns = "sandboxes"

func newClient(t *testing.T, version string, objects ...runtime.Object) (*Client, *dynamicfake.FakeDynamicClient, *kubefake.Clientset) {
	t.Helper()
	shape, err := ShapeFor(version)
	if err != nil {
		t.Fatal(err)
	}
	var dynObjects, kubeObjects []runtime.Object
	for _, o := range objects {
		if _, ok := o.(*unstructured.Unstructured); ok {
			dynObjects = append(dynObjects, o)
		} else {
			kubeObjects = append(kubeObjects, o)
		}
	}
	listKinds := map[schema.GroupVersionResource]string{
		shape.Claims():    "SandboxClaimList",
		shape.Sandboxes(): "SandboxList",
	}
	// The fake seeds objects under a resource name guessed from the kind
	// ("sandboxs"), so they are created through the real resource instead.
	dyn := dynamicfake.NewSimpleDynamicClientWithCustomListKinds(runtime.NewScheme(), listKinds)
	for _, o := range dynObjects {
		u := o.(*unstructured.Unstructured)
		gvr := shape.Claims()
		if u.GetKind() == "Sandbox" {
			gvr = shape.Sandboxes()
		}
		if _, err := dyn.Resource(gvr).Namespace(ns).Create(context.Background(), u, metav1.CreateOptions{}); err != nil {
			t.Fatal(err)
		}
	}
	kube := kubefake.NewSimpleClientset(kubeObjects...)
	opts := Options{Shape: shape, Namespace: ns, EnvironmentID: "env_1", TemplateName: "tpl", WarmPoolName: "pool"}
	if shape.Direct() {
		opts.Nonces = newFakeNonces()
	}
	c, err := New(dyn, kube, opts)
	if err != nil {
		t.Fatal(err)
	}
	return c, dyn, kube
}

// fakeNonces records the nonce for each sandbox, standing in for the registry.
type fakeNonces struct {
	mu  sync.Mutex
	set map[string]string
}

func newFakeNonces() *fakeNonces { return &fakeNonces{set: map[string]string{}} }

func (f *fakeNonces) Mint(name string) (string, error) {
	nonce := "nonce-for-" + name
	f.Set(name, nonce)
	return nonce, nil
}

func (f *fakeNonces) Set(name, nonce string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.set[name] = nonce
}

func (f *fakeNonces) Clear(name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.set, name)
}

func (f *fakeNonces) Known(name string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.set[name]
	return ok
}

func (f *fakeNonces) get(name string) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.set[name]
}

func object(apiVersion, kind, name string, uid types.UID, extra map[string]any) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": apiVersion,
		"kind":       kind,
		"metadata":   map[string]any{"name": name, "namespace": ns, "uid": string(uid)},
	}}
	for k, v := range extra {
		obj.Object[k] = v
	}
	return obj
}

func ownedBy(kind, name string, uid types.UID) map[string]any {
	return map[string]any{"ownerReferences": []any{map[string]any{
		"apiVersion": "x/v1", "kind": kind, "name": name, "uid": string(uid), "controller": true,
	}}}
}

// ownedGraph is a bound claim, its Ready Sandbox and the pod the Sandbox
// names, with the pod Pending.
func ownedGraph(claimLabels map[string]string) []runtime.Object {
	claim := object("extensions.agents.x-k8s.io/v1beta1", "SandboxClaim", "claim-a", "claim-a-uid", map[string]any{
		"status": map[string]any{"sandbox": map[string]any{"name": "sb-a"}},
	})
	claim.SetLabels(claimLabels)
	sb := object("agents.x-k8s.io/v1beta1", "Sandbox", "sb-a", "sb-a-uid", map[string]any{
		"status": map[string]any{"conditions": []any{map[string]any{"type": "Ready", "status": "True"}}},
	})
	sb.Object["metadata"].(map[string]any)["ownerReferences"] = ownedBy("SandboxClaim", "claim-a", "claim-a-uid")["ownerReferences"]
	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "sb-a", Namespace: ns, UID: "pod-a-uid",
		OwnerReferences: []metav1.OwnerReference{{Kind: "Sandbox", Name: "sb-a", UID: "sb-a-uid", Controller: ptr(true)}}}}
	return []runtime.Object{claim, sb, pod}
}

func ptr[T any](v T) *T { return &v }
